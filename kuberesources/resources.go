package kuberesources

import (
	"fmt"
	"time"

	contrailv1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/michaelhenkel/ckube/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextinstall "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	extscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ContrailResources struct {
	Namespace  corev1.Namespace
	Rbac       RbacResources
	CRDs       CRDResources
	Deployment appsv1.Deployment
	Manager    contrailv1.Manager
}

type RbacResources struct {
	Role               rbacv1.Role
	RoleBinding        rbacv1.RoleBinding
	ClusterRole        rbacv1.ClusterRole
	ClusterRoleBinding rbacv1.ClusterRoleBinding
	ServiceAccount     corev1.ServiceAccount
}

type CRDResources struct {
	Cassandra        apiextensionsv1.CustomResourceDefinition
	Zookeeper        apiextensionsv1.CustomResourceDefinition
	Rabbitmq         apiextensionsv1.CustomResourceDefinition
	Config           apiextensionsv1.CustomResourceDefinition
	Control          apiextensionsv1.CustomResourceDefinition
	Kubemanager      apiextensionsv1.CustomResourceDefinition
	Webui            apiextensionsv1.CustomResourceDefinition
	Vrouter          apiextensionsv1.CustomResourceDefinition
	Manager          apiextensionsv1.CustomResourceDefinition
	Provisionmanager apiextensionsv1.CustomResourceDefinition
}

type managerClient struct {
	restClient rest.Interface
	ns         string
}

type crdClient struct {
	restClient rest.Interface
	ns         string
}

func (c *managerClient) Create(name string, object *contrailv1.Manager) (*contrailv1.Manager, error) {
	result := contrailv1.Manager{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("managers").
		Name(name).
		Body(object).
		Do().
		Into(&result)

	return &result, err
}

func kubeClient(configPath string) (*kubernetes.Clientset, *rest.RESTClient, *apiextensionsclientset.Clientset, error) {
	var err error
	clientset := &kubernetes.Clientset{}
	restClient := &rest.RESTClient{}
	crdClient := &apiextensionsclientset.Clientset{}
	kubeConfig := &rest.Config{}

	kubeConfig, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return clientset, restClient, crdClient, err
	}

	// create the clientset
	contrailv1.SchemeBuilder.AddToScheme(scheme.Scheme)
	apiextinstall.Install(scheme.Scheme)
	if err := extscheme.AddToScheme(scheme.Scheme); err != nil {
		return clientset, restClient, crdClient, err
	}
	//apiextensionsv1.SchemeBuilder.AddToScheme(scheme.Scheme)

	crConfig := kubeConfig
	crConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: contrailv1.SchemeGroupVersion.Group, Version: contrailv1.SchemeGroupVersion.Version}
	crConfig.APIPath = "/apis"

	//crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	crConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err = rest.UnversionedRESTClientFor(crConfig)
	if err != nil {
		return clientset, restClient, crdClient, err
	}

	clientset, err = kubernetes.NewForConfig(crConfig)
	if err != nil {
		return clientset, restClient, crdClient, err
	}

	crdClient, err = apiextensionsclientset.NewForConfig(crConfig)
	if err != nil {
		return clientset, restClient, crdClient, err
	}

	return clientset, restClient, crdClient, nil
}

func CreateContrailResources(configPath string) error {
	contrailResources := newContrailResources()
	clientSet, restClient, crdClientSet, err := kubeClient(configPath)
	if err != nil {
		return err
	}

	_, err = clientSet.CoreV1().Namespaces().Create(&contrailResources.Namespace)
	if err != nil {
		return err
	}
	_, err = clientSet.RbacV1().Roles(contrailResources.Namespace.Name).Create(&contrailResources.Rbac.Role)
	if err != nil {
		return err
	}
	_, err = clientSet.RbacV1().RoleBindings(contrailResources.Namespace.Name).Create(&contrailResources.Rbac.RoleBinding)
	if err != nil {
		return err
	}
	_, err = clientSet.RbacV1().ClusterRoles().Create(&contrailResources.Rbac.ClusterRole)
	if err != nil {
		return err
	}
	_, err = clientSet.RbacV1().ClusterRoleBindings().Create(&contrailResources.Rbac.ClusterRoleBinding)
	if err != nil {
		return err
	}
	_, err = clientSet.CoreV1().ServiceAccounts(contrailResources.Namespace.Name).Create(&contrailResources.Rbac.ServiceAccount)
	if err != nil {
		return err
	}

	fmt.Println("Creating CRDs")
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Manager)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Cassandra)
	if err != nil {
		return err
	}

	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Zookeeper)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Rabbitmq)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Config)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Control)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Kubemanager)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Webui)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Vrouter)
	if err != nil {
		return err
	}
	_, err = crdClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&contrailResources.CRDs.Provisionmanager)
	if err != nil {
		return err
	}

	fmt.Println("CRDs created")

	err = utils.Retry(10, 2*time.Second, func() (err error) {
		fmt.Println("trying to get Manager CRD")
		err = getResource(crdClientSet, contrailResources.CRDs.Manager.Name)
		return
	})

	if err != nil {
		return err
	}

	managerClient := &managerClient{
		ns:         contrailResources.Namespace.Name,
		restClient: restClient,
	}

	err = utils.Retry(10, 2*time.Second, func() (err error) {
		fmt.Println("Creating manager CR")
		_, err = managerClient.Create(contrailResources.Manager.Name, &contrailResources.Manager)
		return
	})
	if err != nil {
		return err
	}

	fmt.Println("Creating deployment")
	_, err = clientSet.AppsV1().Deployments(contrailResources.Namespace.Name).Create(&contrailResources.Deployment)
	if err != nil {
		return err
	}

	return nil
}

func getResource(clientSet *apiextensionsclientset.Clientset, name string) error {
	getOptions := metav1.GetOptions{}
	fmt.Println("Getting manager crd ", name)
	_, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Get(name, getOptions)
	if err != nil {
		return err
	}
	return nil
}

func newContrailResources() ContrailResources {
	contrailResources := ContrailResources{}

	namespaceRes := corev1.Namespace{}
	err := yaml.Unmarshal([]byte(namespace), &namespaceRes)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(namespace))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &namespaceRes)
	if err != nil {
		panic(err)
	}
	contrailResources.Namespace = namespaceRes

	rbacResources := RbacResources{}

	roleRes := rbacv1.Role{}
	err = yaml.Unmarshal([]byte(role), &roleRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(role))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &roleRes)
	if err != nil {
		panic(err)
	}
	rbacResources.Role = roleRes

	roleBindingRes := rbacv1.RoleBinding{}
	err = yaml.Unmarshal([]byte(roleBinding), &roleBindingRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(roleBinding))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &roleBindingRes)
	if err != nil {
		panic(err)
	}
	rbacResources.RoleBinding = roleBindingRes

	clusterRoleRes := rbacv1.ClusterRole{}
	err = yaml.Unmarshal([]byte(clusterRole), &clusterRoleRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(clusterRole))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &clusterRoleRes)
	if err != nil {
		panic(err)
	}
	rbacResources.ClusterRole = clusterRoleRes

	clusterRoleBindingRes := rbacv1.ClusterRoleBinding{}
	err = yaml.Unmarshal([]byte(clusterRoleBinding), &clusterRoleBindingRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(clusterRoleBinding))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &clusterRoleBindingRes)
	if err != nil {
		panic(err)
	}
	rbacResources.ClusterRoleBinding = clusterRoleBindingRes

	serviceAccountRes := corev1.ServiceAccount{}
	err = yaml.Unmarshal([]byte(serviceAccount), &serviceAccountRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(serviceAccount))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &serviceAccountRes)
	if err != nil {
		panic(err)
	}
	rbacResources.ServiceAccount = serviceAccountRes

	contrailResources.Rbac = rbacResources

	deploymentRes := appsv1.Deployment{}
	err = yaml.Unmarshal([]byte(deployment), &deploymentRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(deployment))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &deploymentRes)
	if err != nil {
		panic(err)
	}
	contrailResources.Deployment = deploymentRes

	managerCrRes := contrailv1.Manager{}
	err = yaml.Unmarshal([]byte(crManager), &managerCrRes)
	if err != nil {
		panic(err)
	}
	jsonData, err = yaml.YAMLToJSON([]byte(crManager))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &managerCrRes)
	if err != nil {
		panic(err)
	}
	contrailResources.Manager = managerCrRes

	contrailResources.CRDs.Manager = createExtension(crdManager)
	contrailResources.CRDs.Cassandra = createExtension(crdCassandra)
	contrailResources.CRDs.Zookeeper = createExtension(crdZookeeper)
	contrailResources.CRDs.Rabbitmq = createExtension(crdRabbitmq)
	contrailResources.CRDs.Config = createExtension(crdConfig)
	contrailResources.CRDs.Control = createExtension(crdControl)
	contrailResources.CRDs.Kubemanager = createExtension(crdKubemanager)
	contrailResources.CRDs.Webui = createExtension(crdWebui)
	contrailResources.CRDs.Vrouter = createExtension(crdVrouter)
	contrailResources.CRDs.Provisionmanager = createExtension(crdProvisionmanager)

	return contrailResources
}

func createExtension(crd string) apiextensionsv1.CustomResourceDefinition {
	extensionRes := apiextensionsv1.CustomResourceDefinition{}
	err := yaml.Unmarshal([]byte(crd), &extensionRes)
	if err != nil {
		panic(err)
	}
	jsonData, err := yaml.YAMLToJSON([]byte(crd))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(jsonData), &extensionRes)
	if err != nil {
		panic(err)
	}
	return extensionRes
}

var namespace = `apiVersion: v1
kind: Namespace
metadata:
  name: contrail`

var role = `apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: contrail-operator
  namespace: contrail
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  - endpoints
  - persistentvolumeclaims
  - events
  - configmaps
  - secrets
  verbs:
  - '*'
- apiGroups:
  - apps
  resources:
  - deployments
  - daemonsets
  - replicasets
  - statefulsets
  verbs:
  - '*'
- apiGroups:
  - monitoring.coreos.com
  resources:
  - servicemonitors
  verbs:
  - get
  - create
- apiGroups:
  - apps
  resourceNames:
  - contrail-operator
  resources:
  - deployments/finalizers
  verbs:
  - update
- apiGroups:
  - contrail.juniper.net
  resources:
  - '*'
  - managers
  - cassandras
  - zookeepers
  - rabbitmqs
  - controls
  - kubemanagers
  - webuis
  - vrouters
  - provisionmanagers
  verbs:
  - '*'
- apiGroups:
  - storage
  resources:
  - storageclasses
  verbs:
  - '*'`

var clusterRole = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: contrail-operator
  namespace: contrail
rules:
  - apiGroups:
    - "*"
    resources:
    - "*"
    verbs:
    - "*"`

var serviceAccount = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: contrail-operator
  namespace: contrail`

var roleBinding = `kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: contrail-operator
  namespace: contrail
subjects:
- kind: ServiceAccount
  name: contrail-operator
  namespace: contrail
roleRef:
  kind: Role
  name: contrail-operator
  apiGroup: rbac.authorization.k8s.io`

var clusterRoleBinding = `kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: contrail-operator
  namespace: contrail
subjects:
- kind: ServiceAccount
  name: contrail-operator
  namespace: contrail
roleRef:
  kind: ClusterRole
  name: contrail-operator
  apiGroup: rbac.authorization.k8s.io`

var deployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: contrail-operator
  namespace: contrail
spec:
  replicas: 1
  selector:
    matchLabels:
      name: contrail-operator
  template:
    metadata:
      labels:
        name: contrail-operator
    spec:
      serviceAccountName: contrail-operator
      hostNetwork: true
      tolerations:
      - operator: Exists
        effect: NoSchedule
      - operator: Exists
        effect: NoExecute
      containers:
       - name: contrail-operator
          # Replace this with the built image name
         image: michaelhenkel/contrail-operator:latest
         command:
         - contrail-operator
         imagePullPolicy: Always
         env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "contrail-operator"`

var crManager = `apiVersion: contrail.juniper.net/v1alpha1
kind: Manager
metadata:
  name: cluster1
  namespace: contrail
spec:
  commonConfiguration:
    hostNetwork: true
    replicas: 1
  services:
    cassandras:
    - metadata:
        labels:
          contrail_cluster: cluster1
        name: cassandra1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          maxHeapSize: 1024M
          minHeapSize: 100M
          containers:
          - name: cassandra
            image: cassandra:3.11.4
          - name: init
            image: python:alpine
          - name: init2
            image: cassandra:3.11.4
    config:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: config1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: analyticsapi
            image: opencontrailnightly/contrail-analytics-api:1910-latest
          - name: api
            image: opencontrailnightly/contrail-controller-config-api:1910-latest
          - name: collector
            image: opencontrailnightly/contrail-analytics-collector:1910-latest
          - name: devicemanager
            image: opencontrailnightly/contrail-controller-config-devicemgr:1910-latest
          - name: init
            image: python:alpine
          - name: init2
            image: busybox
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          - name: redis
            image: redis:4.0.2
          - name: schematransformer
            image: opencontrailnightly/contrail-controller-config-schema:1910-latest
          - name: servicemonitor
            image: opencontrailnightly/contrail-controller-config-svcmonitor:1910-latest
          zookeeperInstance: zookeeper1
    controls:
    - metadata:
        labels:
          contrail_cluster: cluster1
          control_role: master
        name: control1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: control
            image: opencontrailnightly/contrail-controller-control-control:1910-latest
          - name: dns
            image: opencontrailnightly/contrail-controller-control-dns:1910-latest
          - name: init
            image: python:alpine
          - name: named
            image: opencontrailnightly/contrail-controller-control-named:1910-latest
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          - name: statusmonitor
            image: michaelhenkel/contrail-statusmonitor:debug
          zookeeperInstance: zookeeper1
    kubemanagers:
    - metadata:
        labels:
          contrail_cluster: cluster1
        name: kubemanager1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: init
            image: python:alpine
          - name: kubemanager
            image: michaelhenkel/contrail-kubernetes-kube-manager:1910-latest
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          ipFabricForwarding: false
          ipFabricSnat: true
          kubernetesTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
          useKubeadmConfig: false
          zookeeperInstance: zookeeper1
    provisionManager:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: provmanager1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
          replicas: 1
        serviceConfiguration:
          containers:
          - name: init
            image: python:alpine
          - name: provisioner
            image: michaelhenkel/contrail-provisioner:debug
    rabbitmq:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: rabbitmq1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          containers:
          - name: init
            image: python:alpine
          - name: rabbitmq
            image: rabbitmq:3.7
    vrouters:
    - metadata:
        labels:
          contrail_cluster: cluster1
        name: vroutermaster
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: init
            image: python:alpine
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          - name: vrouteragent
            image: opencontrailnightly/contrail-vrouter-agent:1910-latest
            command: ["bash","-c","/usr/bin/contrail-vrouter-agent --config_file /etc/mycontrail/vrouter.${POD_IP}"]
          - name: vroutercni
            image: michaelhenkel/contrailcni:v0.0.1
          - name: vrouterkernelbuildinit
            image: opencontrailnightly/contrail-vrouter-kernel-build-init:1910-latest
          - name: vrouterkernelinit
            image: opencontrailnightly/contrail-vrouter-kernel-init:1910-latest
          controlInstance: control1
          controlInstance: control1
    - metadata:
        labels:
          contrail_cluster: cluster1
        name: vrouternodes
      spec:
        commonConfiguration:
          create: false
          nodeSelector:
            node-role.opencontrail.org: vrouter
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: init
            image: python:alpine
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          - name: vrouteragent
            image: opencontrailnightly/contrail-vrouter-agent:1910-latest
            command: ["bash","-c","/usr/bin/contrail-vrouter-agent --config_file /etc/mycontrail/vrouter.${POD_IP}"]
          - name: vroutercni
            image: michaelhenkel/contrailcni:v0.0.1
          - name: vrouterkernelbuildinit
            image: opencontrailnightly/contrail-vrouter-kernel-build-init:1910-latest
          - name: vrouterkernelinit
            image: opencontrailnightly/contrail-vrouter-kernel-init:1910-latest
          controlInstance: control1
    webui:
      metadata:
        labels:
          contrail_cluster: cluster1
        name: webui1
      spec:
        commonConfiguration:
          replicas: 1
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          cassandraInstance: cassandra1
          containers:
          - name: init
            image: python:alpine
          - name: nodeinit
            image: opencontrailnightly/contrail-node-init:1910-latest
          - name: redis
            image: redis:4.0.2
          - name: webuijob
            image: opencontrailnightly/contrail-controller-webui-job:1910-latest
          - name: webuiweb
            image: opencontrailnightly/contrail-controller-webui-web:1910-latest
    zookeepers:
    - metadata:
        labels:
          contrail_cluster: cluster1
        name: zookeeper1
      spec:
        commonConfiguration:
          create: true
          nodeSelector:
            node-role.kubernetes.io/master: "true"
        serviceConfiguration:
          containers:
          - name: init
            image: python:alpine
          - name: zookeeper
            image: docker.io/zookeeper:3.5.5`
