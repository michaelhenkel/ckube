package kuberesources

import (
	contrailv1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
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
	Cassandra        apiextensions.CustomResourceDefinition
	Zookeeper        apiextensions.CustomResourceDefinition
	Rabbitmq         apiextensions.CustomResourceDefinition
	Config           apiextensions.CustomResourceDefinition
	Control          apiextensions.CustomResourceDefinition
	Kubemanager      apiextensions.CustomResourceDefinition
	Webui            apiextensions.CustomResourceDefinition
	Vrouter          apiextensions.CustomResourceDefinition
	Manager          apiextensions.CustomResourceDefinition
	Provisionmanager apiextensions.CustomResourceDefinition
}

type managerClient struct {
	restClient rest.Interface
	ns         string
}

func (c *managerClient) Create(name string, object *contrailv1.Manager) (*contrailv1.Manager, error) {
	result := contrailv1.Manager{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource("managers").
		Name(name).
		Body(object).
		Do().
		Into(&result)

	return &result, err
}

func kubeClient(configPath string) (*kubernetes.Clientset, *rest.RESTClient, error) {
	var err error
	clientset := &kubernetes.Clientset{}
	restClient := &rest.RESTClient{}
	kubeConfig := &rest.Config{}

	kubeConfig, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return clientset, restClient, err
	}

	// create the clientset
	contrailv1.SchemeBuilder.AddToScheme(scheme.Scheme)

	crdConfig := kubeConfig
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: contrailv1.SchemeGroupVersion.Group, Version: contrailv1.SchemeGroupVersion.Version}
	crdConfig.APIPath = "/apis"

	//crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err = rest.UnversionedRESTClientFor(crdConfig)
	if err != nil {
		return clientset, restClient, err
	}
	clientset, err = kubernetes.NewForConfig(crdConfig)
	if err != nil {
		return clientset, restClient, err
	}
	return clientset, restClient, nil
}

func CreateContrailResources(configPath string) error {
	contrailResources := newContrailResources()
	clientSet, restClient, err := kubeClient(configPath)
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

	managerClient := &managerClient{
		ns:         contrailResources.Namespace.Name,
		restClient: restClient,
	}
	_, err = managerClient.Create(contrailResources.Manager.Name, &contrailResources.Manager)
	if err != nil {
		return err
	}

	/*
		podList, err := clientSet.CoreV1().Pods(config.Namespace).List(metav1.ListOptions{LabelSelector: "control=" + config.NodeName})
		if err != nil {
	*/

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

	contrailResources.CRDs.Cassandra = createExtension(crdCassandra)
	contrailResources.CRDs.Zookeeper = createExtension(crdZookeeper)
	contrailResources.CRDs.Rabbitmq = createExtension(crdRabbitmq)
	contrailResources.CRDs.Config = createExtension(crdConfig)
	contrailResources.CRDs.Control = createExtension(crdControl)
	contrailResources.CRDs.Kubemanager = createExtension(crdKubemanager)
	contrailResources.CRDs.Webui = createExtension(crdWebui)
	contrailResources.CRDs.Vrouter = createExtension(crdVrouter)
	contrailResources.CRDs.Manager = createExtension(crdManager)
	contrailResources.CRDs.Provisionmanager = createExtension(crdProvisionmanager)

	return contrailResources
}

func createExtension(crd string) apiextensions.CustomResourceDefinition {
	extensionRes := apiextensions.CustomResourceDefinition{}
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

var crdCassandra = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: cassandras.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Cassandra
    listKind: CassandraList
    plural: cassandras
    singular: cassandra
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Cassandra is the Schema for the cassandras API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: CassandraSpec is the Spec for the cassandras API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: CassandraConfiguration is the Spec for the cassandras API.
              properties:
                clusterName:
                  type: string
                containers: {}
                cqlPort:
                  type: integer
                jmxLocalPort:
                  type: integer
                listenAddress:
                  type: string
                maxHeapSize:
                  type: string
                minHeapSize:
                  type: string
                port:
                  type: integer
                sslStoragePort:
                  type: integer
                startRPC:
                  type: boolean
                storagePath:
                  type: string
                storagePort:
                  type: integer
                storageSize:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          description: CassandraStatus defines the status of the cassandra object.
          properties:
            active:
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              description: CassandraStatusPorts defines the status of the ports of
                the cassandra object.
              properties:
                cqlPort:
                  type: string
                jmxPort:
                  type: string
                port:
                  type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdConfig = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: configs.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Config
    listKind: ConfigList
    plural: configs
    singular: config
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Config is the Schema for the configs API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ConfigSpec is the Spec for the cassandras API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: ConfigConfiguration is the Spec for the cassandras API.
              properties:
                analyticsPort:
                  type: integer
                apiPort:
                  type: integer
                cassandraInstance:
                  type: string
                collectorPort:
                  type: integer
                containers: {}
                nodeManager:
                  type: boolean
                rabbitmqPassword:
                  type: string
                rabbitmqUser:
                  type: string
                rabbitmqVhost:
                  type: string
                redisPort:
                  type: integer
                zookeeperInstance:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            configChanged:
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              properties:
                analyticsPort:
                  type: string
                apiPort:
                  type: string
                collectorPort:
                  type: string
                redisPort:
                  type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdControl = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: controls.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Control
    listKind: ControlList
    plural: controls
    singular: control
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Control is the Schema for the controls API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ControlSpec is the Spec for the controls API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: ControlConfiguration is the Spec for the controls API.
              properties:
                asnNumber:
                  type: integer
                bgpPort:
                  type: integer
                cassandraInstance:
                  type: string
                containers: {}
                dnsIntrospectPort:
                  type: integer
                dnsPort:
                  type: integer
                nodeManager:
                  type: boolean
                rabbitmqPassword:
                  type: string
                rabbitmqUser:
                  type: string
                rabbitmqVhost:
                  type: string
                xmppPort:
                  type: integer
                zookeeperInstance:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              properties:
                asnNumber:
                  type: string
                bgpPort:
                  type: string
                dnsIntrospectPort:
                  type: string
                dnsPort:
                  type: string
                xmppPort:
                  type: string
              type: object
            serviceStatus:
              additionalProperties:
                type: object
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdKubemanager = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kubemanagers.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Kubemanager
    listKind: KubemanagerList
    plural: kubemanagers
    singular: kubemanager
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubemanager is the Schema for the kubemanagers API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: KubemanagerSpec is the Spec for the kubemanagers API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: KubemanagerConfiguration is the Spec for the kubemanagers
                API.
              properties:
                cassandraInstance:
                  type: string
                cloudOrchestrator:
                  type: string
                clusterRole:
                  type: string
                clusterRoleBinding:
                  type: string
                containers: {}
                hostNetworkService:
                  type: boolean
                ipFabricForwarding:
                  type: boolean
                ipFabricSnat:
                  type: boolean
                ipFabricSubnets:
                  type: string
                kubernetesAPIPort:
                  type: integer
                kubernetesAPISSLPort:
                  type: integer
                kubernetesAPIServer:
                  type: string
                kubernetesClusterName:
                  type: string
                kubernetesTokenFile:
                  type: string
                podSubnets:
                  type: string
                rabbitmqPassword:
                  type: string
                rabbitmqUser:
                  type: string
                rabbitmqVhost:
                  type: string
                serviceAccount:
                  type: string
                serviceSubnets:
                  type: string
                useKubeadmConfig:
                  type: boolean
                zookeeperInstance:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            configChanged:
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdManager = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: managers.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Manager
    listKind: ManagerList
    plural: managers
    singular: manager
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Manager is the Schema for the managers API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagerSpec defines the desired state of Manager.
          properties:
            commonConfiguration:
              description: 'INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
                Important: Run "operator-sdk generate k8s" to regenerate code after
                modifying this file Add custom validation using kubebuilder tags:
                https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            services:
              description: Services defines the desired state of Services.
              properties:
                cassandras:
                  items:
                    description: Cassandra is the Schema for the cassandras API.
                    properties:
                      apiVersion:
                        description: 'APIVersion defines the versioned schema of this
                          representation of an object. Servers should convert recognized
                          schemas to the latest internal value, and may reject unrecognized
                          values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                        type: string
                      kind:
                        description: 'Kind is a string value representing the REST
                          resource this object represents. Servers may infer this
                          from the endpoint the client submits requests to. Cannot
                          be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        type: string
                      metadata:
                        type: object
                      spec:
                        description: CassandraSpec is the Spec for the cassandras
                          API.
                        properties:
                          commonConfiguration:
                            description: CommonConfiguration is the common services
                              struct.
                            properties:
                              activate:
                                description: Activate defines if the service will
                                  be activated by Manager.
                                type: boolean
                              create:
                                description: Create defines if the service will be
                                  created by Manager.
                                type: boolean
                              hostNetwork:
                                description: Host networking requested for this pod.
                                  Use the host's network namespace. If this option
                                  is set, the ports that will be used must be specified.
                                  Default to false.
                                type: boolean
                              imagePullSecrets:
                                description: ImagePullSecrets is an optional list
                                  of references to secrets in the same namespace to
                                  use for pulling any of the images used by this PodSpec.
                                items:
                                  type: string
                                type: array
                              nodeSelector:
                                additionalProperties:
                                  type: string
                                description: 'NodeSelector is a selector which must
                                  be true for the pod to fit on a node. Selector which
                                  must match a node''s labels for the pod to be scheduled
                                  on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                                type: object
                              replicas:
                                description: Number of desired pods. This is a pointer
                                  to distinguish between explicit zero and not specified.
                                  Defaults to 1.
                                format: int32
                                type: integer
                              tolerations:
                                description: If specified, the pod's tolerations.
                                items:
                                  description: The pod this Toleration is attached
                                    to tolerates any taint that matches the triple
                                    <key,value,effect> using the matching operator
                                    <operator>.
                                  properties:
                                    effect:
                                      description: Effect indicates the taint effect
                                        to match. Empty means match all taint effects.
                                        When specified, allowed values are NoSchedule,
                                        PreferNoSchedule and NoExecute.
                                      type: string
                                    key:
                                      description: Key is the taint key that the toleration
                                        applies to. Empty means match all taint keys.
                                        If the key is empty, operator must be Exists;
                                        this combination means to match all values
                                        and all keys.
                                      type: string
                                    operator:
                                      description: Operator represents a key's relationship
                                        to the value. Valid operators are Exists and
                                        Equal. Defaults to Equal. Exists is equivalent
                                        to wildcard for value, so that a pod can tolerate
                                        all taints of a particular category.
                                      type: string
                                    tolerationSeconds:
                                      description: TolerationSeconds represents the
                                        period of time the toleration (which must
                                        be of effect NoExecute, otherwise this field
                                        is ignored) tolerates the taint. By default,
                                        it is not set, which means tolerate the taint
                                        forever (do not evict). Zero and negative
                                        values will be treated as 0 (evict immediately)
                                        by the system.
                                      format: int64
                                      type: integer
                                    value:
                                      description: Value is the taint value the toleration
                                        matches to. If the operator is Exists, the
                                        value should be empty, otherwise just a regular
                                        string.
                                      type: string
                                  type: object
                                type: array
                            type: object
                          serviceConfiguration:
                            description: CassandraConfiguration is the Spec for the
                              cassandras API.
                            properties:
                              clusterName:
                                type: string
                              containers: {}
                              cqlPort:
                                type: integer
                              jmxLocalPort:
                                type: integer
                              listenAddress:
                                type: string
                              maxHeapSize:
                                type: string
                              minHeapSize:
                                type: string
                              port:
                                type: integer
                              sslStoragePort:
                                type: integer
                              startRPC:
                                type: boolean
                              storagePath:
                                type: string
                              storagePort:
                                type: integer
                              storageSize:
                                type: string
                            type: object
                        required:
                        - commonConfiguration
                        - serviceConfiguration
                        type: object
                      status:
                        description: CassandraStatus defines the status of the cassandra
                          object.
                        properties:
                          active:
                            type: boolean
                          nodes:
                            additionalProperties:
                              type: string
                            type: object
                          ports:
                            description: CassandraStatusPorts defines the status of
                              the ports of the cassandra object.
                            properties:
                              cqlPort:
                                type: string
                              jmxPort:
                                type: string
                              port:
                                type: string
                            type: object
                        type: object
                    type: object
                  type: array
                config:
                  description: Config is the Schema for the configs API.
                  properties:
                    apiVersion:
                      description: 'APIVersion defines the versioned schema of this
                        representation of an object. Servers should convert recognized
                        schemas to the latest internal value, and may reject unrecognized
                        values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                      type: string
                    kind:
                      description: 'Kind is a string value representing the REST resource
                        this object represents. Servers may infer this from the endpoint
                        the client submits requests to. Cannot be updated. In CamelCase.
                        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    metadata:
                      type: object
                    spec:
                      description: ConfigSpec is the Spec for the cassandras API.
                      properties:
                        commonConfiguration:
                          description: CommonConfiguration is the common services
                            struct.
                          properties:
                            activate:
                              description: Activate defines if the service will be
                                activated by Manager.
                              type: boolean
                            create:
                              description: Create defines if the service will be created
                                by Manager.
                              type: boolean
                            hostNetwork:
                              description: Host networking requested for this pod.
                                Use the host's network namespace. If this option is
                                set, the ports that will be used must be specified.
                                Default to false.
                              type: boolean
                            imagePullSecrets:
                              description: ImagePullSecrets is an optional list of
                                references to secrets in the same namespace to use
                                for pulling any of the images used by this PodSpec.
                              items:
                                type: string
                              type: array
                            nodeSelector:
                              additionalProperties:
                                type: string
                              description: 'NodeSelector is a selector which must
                                be true for the pod to fit on a node. Selector which
                                must match a node''s labels for the pod to be scheduled
                                on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                              type: object
                            replicas:
                              description: Number of desired pods. This is a pointer
                                to distinguish between explicit zero and not specified.
                                Defaults to 1.
                              format: int32
                              type: integer
                            tolerations:
                              description: If specified, the pod's tolerations.
                              items:
                                description: The pod this Toleration is attached to
                                  tolerates any taint that matches the triple <key,value,effect>
                                  using the matching operator <operator>.
                                properties:
                                  effect:
                                    description: Effect indicates the taint effect
                                      to match. Empty means match all taint effects.
                                      When specified, allowed values are NoSchedule,
                                      PreferNoSchedule and NoExecute.
                                    type: string
                                  key:
                                    description: Key is the taint key that the toleration
                                      applies to. Empty means match all taint keys.
                                      If the key is empty, operator must be Exists;
                                      this combination means to match all values and
                                      all keys.
                                    type: string
                                  operator:
                                    description: Operator represents a key's relationship
                                      to the value. Valid operators are Exists and
                                      Equal. Defaults to Equal. Exists is equivalent
                                      to wildcard for value, so that a pod can tolerate
                                      all taints of a particular category.
                                    type: string
                                  tolerationSeconds:
                                    description: TolerationSeconds represents the
                                      period of time the toleration (which must be
                                      of effect NoExecute, otherwise this field is
                                      ignored) tolerates the taint. By default, it
                                      is not set, which means tolerate the taint forever
                                      (do not evict). Zero and negative values will
                                      be treated as 0 (evict immediately) by the system.
                                    format: int64
                                    type: integer
                                  value:
                                    description: Value is the taint value the toleration
                                      matches to. If the operator is Exists, the value
                                      should be empty, otherwise just a regular string.
                                    type: string
                                type: object
                              type: array
                          type: object
                        serviceConfiguration:
                          description: ConfigConfiguration is the Spec for the cassandras
                            API.
                          properties:
                            analyticsPort:
                              type: integer
                            apiPort:
                              type: integer
                            cassandraInstance:
                              type: string
                            collectorPort:
                              type: integer
                            containers: {}
                            nodeManager:
                              type: boolean
                            rabbitmqPassword:
                              type: string
                            rabbitmqUser:
                              type: string
                            rabbitmqVhost:
                              type: string
                            redisPort:
                              type: integer
                            zookeeperInstance:
                              type: string
                          type: object
                      required:
                      - commonConfiguration
                      - serviceConfiguration
                      type: object
                    status:
                      properties:
                        active:
                          description: 'INSERT ADDITIONAL STATUS FIELD - define observed
                            state of cluster Important: Run "operator-sdk generate
                            k8s" to regenerate code after modifying this file Add
                            custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                          type: boolean
                        configChanged:
                          type: boolean
                        nodes:
                          additionalProperties:
                            type: string
                          type: object
                        ports:
                          properties:
                            analyticsPort:
                              type: string
                            apiPort:
                              type: string
                            collectorPort:
                              type: string
                            redisPort:
                              type: string
                          type: object
                      type: object
                  type: object
                controls:
                  items:
                    description: Control is the Schema for the controls API.
                    properties:
                      apiVersion:
                        description: 'APIVersion defines the versioned schema of this
                          representation of an object. Servers should convert recognized
                          schemas to the latest internal value, and may reject unrecognized
                          values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                        type: string
                      kind:
                        description: 'Kind is a string value representing the REST
                          resource this object represents. Servers may infer this
                          from the endpoint the client submits requests to. Cannot
                          be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        type: string
                      metadata:
                        type: object
                      spec:
                        description: ControlSpec is the Spec for the controls API.
                        properties:
                          commonConfiguration:
                            description: CommonConfiguration is the common services
                              struct.
                            properties:
                              activate:
                                description: Activate defines if the service will
                                  be activated by Manager.
                                type: boolean
                              create:
                                description: Create defines if the service will be
                                  created by Manager.
                                type: boolean
                              hostNetwork:
                                description: Host networking requested for this pod.
                                  Use the host's network namespace. If this option
                                  is set, the ports that will be used must be specified.
                                  Default to false.
                                type: boolean
                              imagePullSecrets:
                                description: ImagePullSecrets is an optional list
                                  of references to secrets in the same namespace to
                                  use for pulling any of the images used by this PodSpec.
                                items:
                                  type: string
                                type: array
                              nodeSelector:
                                additionalProperties:
                                  type: string
                                description: 'NodeSelector is a selector which must
                                  be true for the pod to fit on a node. Selector which
                                  must match a node''s labels for the pod to be scheduled
                                  on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                                type: object
                              replicas:
                                description: Number of desired pods. This is a pointer
                                  to distinguish between explicit zero and not specified.
                                  Defaults to 1.
                                format: int32
                                type: integer
                              tolerations:
                                description: If specified, the pod's tolerations.
                                items:
                                  description: The pod this Toleration is attached
                                    to tolerates any taint that matches the triple
                                    <key,value,effect> using the matching operator
                                    <operator>.
                                  properties:
                                    effect:
                                      description: Effect indicates the taint effect
                                        to match. Empty means match all taint effects.
                                        When specified, allowed values are NoSchedule,
                                        PreferNoSchedule and NoExecute.
                                      type: string
                                    key:
                                      description: Key is the taint key that the toleration
                                        applies to. Empty means match all taint keys.
                                        If the key is empty, operator must be Exists;
                                        this combination means to match all values
                                        and all keys.
                                      type: string
                                    operator:
                                      description: Operator represents a key's relationship
                                        to the value. Valid operators are Exists and
                                        Equal. Defaults to Equal. Exists is equivalent
                                        to wildcard for value, so that a pod can tolerate
                                        all taints of a particular category.
                                      type: string
                                    tolerationSeconds:
                                      description: TolerationSeconds represents the
                                        period of time the toleration (which must
                                        be of effect NoExecute, otherwise this field
                                        is ignored) tolerates the taint. By default,
                                        it is not set, which means tolerate the taint
                                        forever (do not evict). Zero and negative
                                        values will be treated as 0 (evict immediately)
                                        by the system.
                                      format: int64
                                      type: integer
                                    value:
                                      description: Value is the taint value the toleration
                                        matches to. If the operator is Exists, the
                                        value should be empty, otherwise just a regular
                                        string.
                                      type: string
                                  type: object
                                type: array
                            type: object
                          serviceConfiguration:
                            description: ControlConfiguration is the Spec for the
                              controls API.
                            properties:
                              asnNumber:
                                type: integer
                              bgpPort:
                                type: integer
                              cassandraInstance:
                                type: string
                              containers: {}
                              dnsIntrospectPort:
                                type: integer
                              dnsPort:
                                type: integer
                              nodeManager:
                                type: boolean
                              rabbitmqPassword:
                                type: string
                              rabbitmqUser:
                                type: string
                              rabbitmqVhost:
                                type: string
                              xmppPort:
                                type: integer
                              zookeeperInstance:
                                type: string
                            type: object
                        required:
                        - commonConfiguration
                        - serviceConfiguration
                        type: object
                      status:
                        properties:
                          active:
                            type: boolean
                          nodes:
                            additionalProperties:
                              type: string
                            type: object
                          ports:
                            properties:
                              asnNumber:
                                type: string
                              bgpPort:
                                type: string
                              dnsIntrospectPort:
                                type: string
                              dnsPort:
                                type: string
                              xmppPort:
                                type: string
                            type: object
                          serviceStatus:
                            additionalProperties:
                              type: object
                            type: object
                        type: object
                    type: object
                  type: array
                kubemanagers:
                  items:
                    description: Kubemanager is the Schema for the kubemanagers API.
                    properties:
                      apiVersion:
                        description: 'APIVersion defines the versioned schema of this
                          representation of an object. Servers should convert recognized
                          schemas to the latest internal value, and may reject unrecognized
                          values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                        type: string
                      kind:
                        description: 'Kind is a string value representing the REST
                          resource this object represents. Servers may infer this
                          from the endpoint the client submits requests to. Cannot
                          be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        type: string
                      metadata:
                        type: object
                      spec:
                        description: KubemanagerSpec is the Spec for the kubemanagers
                          API.
                        properties:
                          commonConfiguration:
                            description: CommonConfiguration is the common services
                              struct.
                            properties:
                              activate:
                                description: Activate defines if the service will
                                  be activated by Manager.
                                type: boolean
                              create:
                                description: Create defines if the service will be
                                  created by Manager.
                                type: boolean
                              hostNetwork:
                                description: Host networking requested for this pod.
                                  Use the host's network namespace. If this option
                                  is set, the ports that will be used must be specified.
                                  Default to false.
                                type: boolean
                              imagePullSecrets:
                                description: ImagePullSecrets is an optional list
                                  of references to secrets in the same namespace to
                                  use for pulling any of the images used by this PodSpec.
                                items:
                                  type: string
                                type: array
                              nodeSelector:
                                additionalProperties:
                                  type: string
                                description: 'NodeSelector is a selector which must
                                  be true for the pod to fit on a node. Selector which
                                  must match a node''s labels for the pod to be scheduled
                                  on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                                type: object
                              replicas:
                                description: Number of desired pods. This is a pointer
                                  to distinguish between explicit zero and not specified.
                                  Defaults to 1.
                                format: int32
                                type: integer
                              tolerations:
                                description: If specified, the pod's tolerations.
                                items:
                                  description: The pod this Toleration is attached
                                    to tolerates any taint that matches the triple
                                    <key,value,effect> using the matching operator
                                    <operator>.
                                  properties:
                                    effect:
                                      description: Effect indicates the taint effect
                                        to match. Empty means match all taint effects.
                                        When specified, allowed values are NoSchedule,
                                        PreferNoSchedule and NoExecute.
                                      type: string
                                    key:
                                      description: Key is the taint key that the toleration
                                        applies to. Empty means match all taint keys.
                                        If the key is empty, operator must be Exists;
                                        this combination means to match all values
                                        and all keys.
                                      type: string
                                    operator:
                                      description: Operator represents a key's relationship
                                        to the value. Valid operators are Exists and
                                        Equal. Defaults to Equal. Exists is equivalent
                                        to wildcard for value, so that a pod can tolerate
                                        all taints of a particular category.
                                      type: string
                                    tolerationSeconds:
                                      description: TolerationSeconds represents the
                                        period of time the toleration (which must
                                        be of effect NoExecute, otherwise this field
                                        is ignored) tolerates the taint. By default,
                                        it is not set, which means tolerate the taint
                                        forever (do not evict). Zero and negative
                                        values will be treated as 0 (evict immediately)
                                        by the system.
                                      format: int64
                                      type: integer
                                    value:
                                      description: Value is the taint value the toleration
                                        matches to. If the operator is Exists, the
                                        value should be empty, otherwise just a regular
                                        string.
                                      type: string
                                  type: object
                                type: array
                            type: object
                          serviceConfiguration:
                            description: KubemanagerConfiguration is the Spec for
                              the kubemanagers API.
                            properties:
                              cassandraInstance:
                                type: string
                              cloudOrchestrator:
                                type: string
                              clusterRole:
                                type: string
                              clusterRoleBinding:
                                type: string
                              containers: {}
                              hostNetworkService:
                                type: boolean
                              ipFabricForwarding:
                                type: boolean
                              ipFabricSnat:
                                type: boolean
                              ipFabricSubnets:
                                type: string
                              kubernetesAPIPort:
                                type: integer
                              kubernetesAPISSLPort:
                                type: integer
                              kubernetesAPIServer:
                                type: string
                              kubernetesClusterName:
                                type: string
                              kubernetesTokenFile:
                                type: string
                              podSubnets:
                                type: string
                              rabbitmqPassword:
                                type: string
                              rabbitmqUser:
                                type: string
                              rabbitmqVhost:
                                type: string
                              serviceAccount:
                                type: string
                              serviceSubnets:
                                type: string
                              useKubeadmConfig:
                                type: boolean
                              zookeeperInstance:
                                type: string
                            type: object
                        required:
                        - commonConfiguration
                        - serviceConfiguration
                        type: object
                      status:
                        properties:
                          active:
                            description: 'INSERT ADDITIONAL STATUS FIELD - define
                              observed state of cluster Important: Run "operator-sdk
                              generate k8s" to regenerate code after modifying this
                              file Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                            type: boolean
                          configChanged:
                            type: boolean
                          nodes:
                            additionalProperties:
                              type: string
                            type: object
                        type: object
                    type: object
                  type: array
                provisionManager:
                  description: ProvisionManager is the Schema for the provisionmanagers
                    API
                  properties:
                    apiVersion:
                      description: 'APIVersion defines the versioned schema of this
                        representation of an object. Servers should convert recognized
                        schemas to the latest internal value, and may reject unrecognized
                        values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                      type: string
                    kind:
                      description: 'Kind is a string value representing the REST resource
                        this object represents. Servers may infer this from the endpoint
                        the client submits requests to. Cannot be updated. In CamelCase.
                        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    metadata:
                      type: object
                    spec:
                      description: ProvisionManagerSpec defines the desired state
                        of ProvisionManager
                      properties:
                        commonConfiguration:
                          description: CommonConfiguration is the common services
                            struct.
                          properties:
                            activate:
                              description: Activate defines if the service will be
                                activated by Manager.
                              type: boolean
                            create:
                              description: Create defines if the service will be created
                                by Manager.
                              type: boolean
                            hostNetwork:
                              description: Host networking requested for this pod.
                                Use the host's network namespace. If this option is
                                set, the ports that will be used must be specified.
                                Default to false.
                              type: boolean
                            imagePullSecrets:
                              description: ImagePullSecrets is an optional list of
                                references to secrets in the same namespace to use
                                for pulling any of the images used by this PodSpec.
                              items:
                                type: string
                              type: array
                            nodeSelector:
                              additionalProperties:
                                type: string
                              description: 'NodeSelector is a selector which must
                                be true for the pod to fit on a node. Selector which
                                must match a node''s labels for the pod to be scheduled
                                on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                              type: object
                            replicas:
                              description: Number of desired pods. This is a pointer
                                to distinguish between explicit zero and not specified.
                                Defaults to 1.
                              format: int32
                              type: integer
                            tolerations:
                              description: If specified, the pod's tolerations.
                              items:
                                description: The pod this Toleration is attached to
                                  tolerates any taint that matches the triple <key,value,effect>
                                  using the matching operator <operator>.
                                properties:
                                  effect:
                                    description: Effect indicates the taint effect
                                      to match. Empty means match all taint effects.
                                      When specified, allowed values are NoSchedule,
                                      PreferNoSchedule and NoExecute.
                                    type: string
                                  key:
                                    description: Key is the taint key that the toleration
                                      applies to. Empty means match all taint keys.
                                      If the key is empty, operator must be Exists;
                                      this combination means to match all values and
                                      all keys.
                                    type: string
                                  operator:
                                    description: Operator represents a key's relationship
                                      to the value. Valid operators are Exists and
                                      Equal. Defaults to Equal. Exists is equivalent
                                      to wildcard for value, so that a pod can tolerate
                                      all taints of a particular category.
                                    type: string
                                  tolerationSeconds:
                                    description: TolerationSeconds represents the
                                      period of time the toleration (which must be
                                      of effect NoExecute, otherwise this field is
                                      ignored) tolerates the taint. By default, it
                                      is not set, which means tolerate the taint forever
                                      (do not evict). Zero and negative values will
                                      be treated as 0 (evict immediately) by the system.
                                    format: int64
                                    type: integer
                                  value:
                                    description: Value is the taint value the toleration
                                      matches to. If the operator is Exists, the value
                                      should be empty, otherwise just a regular string.
                                    type: string
                                type: object
                              type: array
                          type: object
                        serviceConfiguration:
                          description: ProvisionManagerConfiguration defines the provision
                            manager configuration
                          properties:
                            containers: {}
                          type: object
                      required:
                      - commonConfiguration
                      - serviceConfiguration
                      type: object
                    status:
                      description: ProvisionManagerStatus defines the observed state
                        of ProvisionManager
                      properties:
                        active:
                          description: 'INSERT ADDITIONAL STATUS FIELD - define observed
                            state of cluster Important: Run "operator-sdk generate
                            k8s" to regenerate code after modifying this file Add
                            custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html'
                          type: boolean
                        globalConfiguration:
                          additionalProperties:
                            type: string
                          type: object
                        nodes:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                  type: object
                rabbitmq:
                  description: Rabbitmq is the Schema for the rabbitmqs API.
                  properties:
                    apiVersion:
                      description: 'APIVersion defines the versioned schema of this
                        representation of an object. Servers should convert recognized
                        schemas to the latest internal value, and may reject unrecognized
                        values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                      type: string
                    kind:
                      description: 'Kind is a string value representing the REST resource
                        this object represents. Servers may infer this from the endpoint
                        the client submits requests to. Cannot be updated. In CamelCase.
                        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    metadata:
                      type: object
                    spec:
                      description: RabbitmqSpec is the Spec for the cassandras API.
                      properties:
                        commonConfiguration:
                          description: CommonConfiguration is the common services
                            struct.
                          properties:
                            activate:
                              description: Activate defines if the service will be
                                activated by Manager.
                              type: boolean
                            create:
                              description: Create defines if the service will be created
                                by Manager.
                              type: boolean
                            hostNetwork:
                              description: Host networking requested for this pod.
                                Use the host's network namespace. If this option is
                                set, the ports that will be used must be specified.
                                Default to false.
                              type: boolean
                            imagePullSecrets:
                              description: ImagePullSecrets is an optional list of
                                references to secrets in the same namespace to use
                                for pulling any of the images used by this PodSpec.
                              items:
                                type: string
                              type: array
                            nodeSelector:
                              additionalProperties:
                                type: string
                              description: 'NodeSelector is a selector which must
                                be true for the pod to fit on a node. Selector which
                                must match a node''s labels for the pod to be scheduled
                                on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                              type: object
                            replicas:
                              description: Number of desired pods. This is a pointer
                                to distinguish between explicit zero and not specified.
                                Defaults to 1.
                              format: int32
                              type: integer
                            tolerations:
                              description: If specified, the pod's tolerations.
                              items:
                                description: The pod this Toleration is attached to
                                  tolerates any taint that matches the triple <key,value,effect>
                                  using the matching operator <operator>.
                                properties:
                                  effect:
                                    description: Effect indicates the taint effect
                                      to match. Empty means match all taint effects.
                                      When specified, allowed values are NoSchedule,
                                      PreferNoSchedule and NoExecute.
                                    type: string
                                  key:
                                    description: Key is the taint key that the toleration
                                      applies to. Empty means match all taint keys.
                                      If the key is empty, operator must be Exists;
                                      this combination means to match all values and
                                      all keys.
                                    type: string
                                  operator:
                                    description: Operator represents a key's relationship
                                      to the value. Valid operators are Exists and
                                      Equal. Defaults to Equal. Exists is equivalent
                                      to wildcard for value, so that a pod can tolerate
                                      all taints of a particular category.
                                    type: string
                                  tolerationSeconds:
                                    description: TolerationSeconds represents the
                                      period of time the toleration (which must be
                                      of effect NoExecute, otherwise this field is
                                      ignored) tolerates the taint. By default, it
                                      is not set, which means tolerate the taint forever
                                      (do not evict). Zero and negative values will
                                      be treated as 0 (evict immediately) by the system.
                                    format: int64
                                    type: integer
                                  value:
                                    description: Value is the taint value the toleration
                                      matches to. If the operator is Exists, the value
                                      should be empty, otherwise just a regular string.
                                    type: string
                                type: object
                              type: array
                          type: object
                        serviceConfiguration:
                          description: RabbitmqConfiguration is the Spec for the cassandras
                            API.
                          properties:
                            containers: {}
                            erlangCookie:
                              type: string
                            password:
                              type: string
                            port:
                              type: integer
                            secret:
                              type: string
                            sslPort:
                              type: integer
                            user:
                              type: string
                            vhost:
                              type: string
                          type: object
                      required:
                      - commonConfiguration
                      - serviceConfiguration
                      type: object
                    status:
                      properties:
                        active:
                          description: 'INSERT ADDITIONAL STATUS FIELD - define observed
                            state of cluster Important: Run "operator-sdk generate
                            k8s" to regenerate code after modifying this file Add
                            custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                          type: boolean
                        nodes:
                          additionalProperties:
                            type: string
                          type: object
                        ports:
                          properties:
                            port:
                              type: string
                            sslPort:
                              type: string
                          type: object
                        secret:
                          type: string
                      type: object
                  type: object
                vrouters:
                  items:
                    description: Vrouter is the Schema for the vrouters API.
                    properties:
                      apiVersion:
                        description: 'APIVersion defines the versioned schema of this
                          representation of an object. Servers should convert recognized
                          schemas to the latest internal value, and may reject unrecognized
                          values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                        type: string
                      kind:
                        description: 'Kind is a string value representing the REST
                          resource this object represents. Servers may infer this
                          from the endpoint the client submits requests to. Cannot
                          be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        type: string
                      metadata:
                        type: object
                      spec:
                        description: VrouterSpec is the Spec for the cassandras API.
                        properties:
                          commonConfiguration:
                            description: CommonConfiguration is the common services
                              struct.
                            properties:
                              activate:
                                description: Activate defines if the service will
                                  be activated by Manager.
                                type: boolean
                              create:
                                description: Create defines if the service will be
                                  created by Manager.
                                type: boolean
                              hostNetwork:
                                description: Host networking requested for this pod.
                                  Use the host's network namespace. If this option
                                  is set, the ports that will be used must be specified.
                                  Default to false.
                                type: boolean
                              imagePullSecrets:
                                description: ImagePullSecrets is an optional list
                                  of references to secrets in the same namespace to
                                  use for pulling any of the images used by this PodSpec.
                                items:
                                  type: string
                                type: array
                              nodeSelector:
                                additionalProperties:
                                  type: string
                                description: 'NodeSelector is a selector which must
                                  be true for the pod to fit on a node. Selector which
                                  must match a node''s labels for the pod to be scheduled
                                  on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                                type: object
                              replicas:
                                description: Number of desired pods. This is a pointer
                                  to distinguish between explicit zero and not specified.
                                  Defaults to 1.
                                format: int32
                                type: integer
                              tolerations:
                                description: If specified, the pod's tolerations.
                                items:
                                  description: The pod this Toleration is attached
                                    to tolerates any taint that matches the triple
                                    <key,value,effect> using the matching operator
                                    <operator>.
                                  properties:
                                    effect:
                                      description: Effect indicates the taint effect
                                        to match. Empty means match all taint effects.
                                        When specified, allowed values are NoSchedule,
                                        PreferNoSchedule and NoExecute.
                                      type: string
                                    key:
                                      description: Key is the taint key that the toleration
                                        applies to. Empty means match all taint keys.
                                        If the key is empty, operator must be Exists;
                                        this combination means to match all values
                                        and all keys.
                                      type: string
                                    operator:
                                      description: Operator represents a key's relationship
                                        to the value. Valid operators are Exists and
                                        Equal. Defaults to Equal. Exists is equivalent
                                        to wildcard for value, so that a pod can tolerate
                                        all taints of a particular category.
                                      type: string
                                    tolerationSeconds:
                                      description: TolerationSeconds represents the
                                        period of time the toleration (which must
                                        be of effect NoExecute, otherwise this field
                                        is ignored) tolerates the taint. By default,
                                        it is not set, which means tolerate the taint
                                        forever (do not evict). Zero and negative
                                        values will be treated as 0 (evict immediately)
                                        by the system.
                                      format: int64
                                      type: integer
                                    value:
                                      description: Value is the taint value the toleration
                                        matches to. If the operator is Exists, the
                                        value should be empty, otherwise just a regular
                                        string.
                                      type: string
                                  type: object
                                type: array
                            type: object
                          serviceConfiguration:
                            description: VrouterConfiguration is the Spec for the
                              cassandras API.
                            properties:
                              cassandraInstance:
                                type: string
                              clusterRole:
                                type: string
                              clusterRoleBinding:
                                type: string
                              containers: {}
                              controlInstance:
                                type: string
                              distribution:
                                type: string
                              gateway:
                                type: string
                              metaDataSecret:
                                type: string
                              nodeManager:
                                type: boolean
                              physicalInterface:
                                type: string
                              serviceAccount:
                                type: string
                            type: object
                        required:
                        - commonConfiguration
                        - serviceConfiguration
                        type: object
                      status:
                        properties:
                          active:
                            type: boolean
                          nodes:
                            additionalProperties:
                              type: string
                            type: object
                          ports:
                            description: 'INSERT ADDITIONAL STATUS FIELD - define
                              observed state of cluster Important: Run "operator-sdk
                              generate k8s" to regenerate code after modifying this
                              file Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                            properties:
                              analyticsPort:
                                type: string
                              apiPort:
                                type: string
                              collectorPort:
                                type: string
                              redisPort:
                                type: string
                            type: object
                        type: object
                    type: object
                  type: array
                webui:
                  description: Webui is the Schema for the webuis API.
                  properties:
                    apiVersion:
                      description: 'APIVersion defines the versioned schema of this
                        representation of an object. Servers should convert recognized
                        schemas to the latest internal value, and may reject unrecognized
                        values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                      type: string
                    kind:
                      description: 'Kind is a string value representing the REST resource
                        this object represents. Servers may infer this from the endpoint
                        the client submits requests to. Cannot be updated. In CamelCase.
                        More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    metadata:
                      type: object
                    spec:
                      description: WebuiSpec is the Spec for the cassandras API.
                      properties:
                        commonConfiguration:
                          description: CommonConfiguration is the common services
                            struct.
                          properties:
                            activate:
                              description: Activate defines if the service will be
                                activated by Manager.
                              type: boolean
                            create:
                              description: Create defines if the service will be created
                                by Manager.
                              type: boolean
                            hostNetwork:
                              description: Host networking requested for this pod.
                                Use the host's network namespace. If this option is
                                set, the ports that will be used must be specified.
                                Default to false.
                              type: boolean
                            imagePullSecrets:
                              description: ImagePullSecrets is an optional list of
                                references to secrets in the same namespace to use
                                for pulling any of the images used by this PodSpec.
                              items:
                                type: string
                              type: array
                            nodeSelector:
                              additionalProperties:
                                type: string
                              description: 'NodeSelector is a selector which must
                                be true for the pod to fit on a node. Selector which
                                must match a node''s labels for the pod to be scheduled
                                on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                              type: object
                            replicas:
                              description: Number of desired pods. This is a pointer
                                to distinguish between explicit zero and not specified.
                                Defaults to 1.
                              format: int32
                              type: integer
                            tolerations:
                              description: If specified, the pod's tolerations.
                              items:
                                description: The pod this Toleration is attached to
                                  tolerates any taint that matches the triple <key,value,effect>
                                  using the matching operator <operator>.
                                properties:
                                  effect:
                                    description: Effect indicates the taint effect
                                      to match. Empty means match all taint effects.
                                      When specified, allowed values are NoSchedule,
                                      PreferNoSchedule and NoExecute.
                                    type: string
                                  key:
                                    description: Key is the taint key that the toleration
                                      applies to. Empty means match all taint keys.
                                      If the key is empty, operator must be Exists;
                                      this combination means to match all values and
                                      all keys.
                                    type: string
                                  operator:
                                    description: Operator represents a key's relationship
                                      to the value. Valid operators are Exists and
                                      Equal. Defaults to Equal. Exists is equivalent
                                      to wildcard for value, so that a pod can tolerate
                                      all taints of a particular category.
                                    type: string
                                  tolerationSeconds:
                                    description: TolerationSeconds represents the
                                      period of time the toleration (which must be
                                      of effect NoExecute, otherwise this field is
                                      ignored) tolerates the taint. By default, it
                                      is not set, which means tolerate the taint forever
                                      (do not evict). Zero and negative values will
                                      be treated as 0 (evict immediately) by the system.
                                    format: int64
                                    type: integer
                                  value:
                                    description: Value is the taint value the toleration
                                      matches to. If the operator is Exists, the value
                                      should be empty, otherwise just a regular string.
                                    type: string
                                type: object
                              type: array
                          type: object
                        serviceConfiguration:
                          description: WebuiConfiguration is the Spec for the cassandras
                            API.
                          properties:
                            cassandraInstance:
                              type: string
                            clusterRole:
                              type: string
                            clusterRoleBinding:
                              type: string
                            containers: {}
                            serviceAccount:
                              type: string
                          type: object
                      required:
                      - commonConfiguration
                      - serviceConfiguration
                      type: object
                    status:
                      properties:
                        active:
                          description: 'INSERT ADDITIONAL STATUS FIELD - define observed
                            state of cluster Important: Run "operator-sdk generate
                            k8s" to regenerate code after modifying this file Add
                            custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                          type: boolean
                        nodes:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                  type: object
                zookeepers:
                  items:
                    description: Zookeeper is the Schema for the zookeepers API.
                    properties:
                      apiVersion:
                        description: 'APIVersion defines the versioned schema of this
                          representation of an object. Servers should convert recognized
                          schemas to the latest internal value, and may reject unrecognized
                          values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                        type: string
                      kind:
                        description: 'Kind is a string value representing the REST
                          resource this object represents. Servers may infer this
                          from the endpoint the client submits requests to. Cannot
                          be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        type: string
                      metadata:
                        type: object
                      spec:
                        description: ZookeeperSpec is the Spec for the zookeepers
                          API.
                        properties:
                          commonConfiguration:
                            description: CommonConfiguration is the common services
                              struct.
                            properties:
                              activate:
                                description: Activate defines if the service will
                                  be activated by Manager.
                                type: boolean
                              create:
                                description: Create defines if the service will be
                                  created by Manager.
                                type: boolean
                              hostNetwork:
                                description: Host networking requested for this pod.
                                  Use the host's network namespace. If this option
                                  is set, the ports that will be used must be specified.
                                  Default to false.
                                type: boolean
                              imagePullSecrets:
                                description: ImagePullSecrets is an optional list
                                  of references to secrets in the same namespace to
                                  use for pulling any of the images used by this PodSpec.
                                items:
                                  type: string
                                type: array
                              nodeSelector:
                                additionalProperties:
                                  type: string
                                description: 'NodeSelector is a selector which must
                                  be true for the pod to fit on a node. Selector which
                                  must match a node''s labels for the pod to be scheduled
                                  on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                                type: object
                              replicas:
                                description: Number of desired pods. This is a pointer
                                  to distinguish between explicit zero and not specified.
                                  Defaults to 1.
                                format: int32
                                type: integer
                              tolerations:
                                description: If specified, the pod's tolerations.
                                items:
                                  description: The pod this Toleration is attached
                                    to tolerates any taint that matches the triple
                                    <key,value,effect> using the matching operator
                                    <operator>.
                                  properties:
                                    effect:
                                      description: Effect indicates the taint effect
                                        to match. Empty means match all taint effects.
                                        When specified, allowed values are NoSchedule,
                                        PreferNoSchedule and NoExecute.
                                      type: string
                                    key:
                                      description: Key is the taint key that the toleration
                                        applies to. Empty means match all taint keys.
                                        If the key is empty, operator must be Exists;
                                        this combination means to match all values
                                        and all keys.
                                      type: string
                                    operator:
                                      description: Operator represents a key's relationship
                                        to the value. Valid operators are Exists and
                                        Equal. Defaults to Equal. Exists is equivalent
                                        to wildcard for value, so that a pod can tolerate
                                        all taints of a particular category.
                                      type: string
                                    tolerationSeconds:
                                      description: TolerationSeconds represents the
                                        period of time the toleration (which must
                                        be of effect NoExecute, otherwise this field
                                        is ignored) tolerates the taint. By default,
                                        it is not set, which means tolerate the taint
                                        forever (do not evict). Zero and negative
                                        values will be treated as 0 (evict immediately)
                                        by the system.
                                      format: int64
                                      type: integer
                                    value:
                                      description: Value is the taint value the toleration
                                        matches to. If the operator is Exists, the
                                        value should be empty, otherwise just a regular
                                        string.
                                      type: string
                                  type: object
                                type: array
                            type: object
                          serviceConfiguration:
                            description: ZookeeperConfiguration is the Spec for the
                              zookeepers API.
                            properties:
                              clientPort:
                                type: integer
                              containers: {}
                              electionPort:
                                type: integer
                              serverPort:
                                type: integer
                            type: object
                        required:
                        - commonConfiguration
                        - serviceConfiguration
                        type: object
                      status:
                        description: ZookeeperStatus defines the status of the zookeeper
                          object.
                        properties:
                          active:
                            description: 'INSERT ADDITIONAL STATUS FIELD - define
                              observed state of cluster Important: Run "operator-sdk
                              generate k8s" to regenerate code after modifying this
                              file Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
                            type: boolean
                          nodes:
                            additionalProperties:
                              type: string
                            type: object
                          ports:
                            description: ZookeeperStatusPorts defines the status of
                              the ports of the zookeeper object.
                            properties:
                              clientPort:
                                type: string
                            type: object
                        type: object
                    type: object
                  type: array
              type: object
          type: object
        status:
          description: ManagerStatus defines the observed state of Manager.
          properties:
            cassandras:
              items:
                description: ServiceStatus provides information on the current status
                  of the service.
                properties:
                  active:
                    type: boolean
                  controllerRunning:
                    type: boolean
                  created:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
            config:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              properties:
                active:
                  type: boolean
                controllerRunning:
                  type: boolean
                created:
                  type: boolean
                name:
                  type: string
              type: object
            controls:
              items:
                description: ServiceStatus provides information on the current status
                  of the service.
                properties:
                  active:
                    type: boolean
                  controllerRunning:
                    type: boolean
                  created:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
            crdStatus:
              items:
                description: CrdStatus tracks status of CRD.
                properties:
                  active:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
            kubemanagers:
              items:
                description: ServiceStatus provides information on the current status
                  of the service.
                properties:
                  active:
                    type: boolean
                  controllerRunning:
                    type: boolean
                  created:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
            provisionManager:
              description: ServiceStatus provides information on the current status
                of the service.
              properties:
                active:
                  type: boolean
                controllerRunning:
                  type: boolean
                created:
                  type: boolean
                name:
                  type: string
              type: object
            rabbitmq:
              description: ServiceStatus provides information on the current status
                of the service.
              properties:
                active:
                  type: boolean
                controllerRunning:
                  type: boolean
                created:
                  type: boolean
                name:
                  type: string
              type: object
            vrouters:
              items:
                description: ServiceStatus provides information on the current status
                  of the service.
                properties:
                  active:
                    type: boolean
                  controllerRunning:
                    type: boolean
                  created:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
            webui:
              description: ServiceStatus provides information on the current status
                of the service.
              properties:
                active:
                  type: boolean
                controllerRunning:
                  type: boolean
                created:
                  type: boolean
                name:
                  type: string
              type: object
            zookeepers:
              items:
                description: ServiceStatus provides information on the current status
                  of the service.
                properties:
                  active:
                    type: boolean
                  controllerRunning:
                    type: boolean
                  created:
                    type: boolean
                  name:
                    type: string
                type: object
              type: array
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdProvisionmanager = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: provisionmanagers.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: ProvisionManager
    listKind: ProvisionManagerList
    plural: provisionmanagers
    singular: provisionmanager
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: ProvisionManager is the Schema for the provisionmanagers API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ProvisionManagerSpec defines the desired state of ProvisionManager
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: ProvisionManagerConfiguration defines the provision manager
                configuration
              properties:
                containers: {}
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          description: ProvisionManagerStatus defines the observed state of ProvisionManager
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            globalConfiguration:
              additionalProperties:
                type: string
              type: object
            nodes:
              additionalProperties:
                type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdRabbitmq = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: rabbitmqs.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Rabbitmq
    listKind: RabbitmqList
    plural: rabbitmqs
    singular: rabbitmq
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Rabbitmq is the Schema for the rabbitmqs API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: RabbitmqSpec is the Spec for the cassandras API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: RabbitmqConfiguration is the Spec for the cassandras API.
              properties:
                containers: {}
                erlangCookie:
                  type: string
                password:
                  type: string
                port:
                  type: integer
                secret:
                  type: string
                sslPort:
                  type: integer
                user:
                  type: string
                vhost:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              properties:
                port:
                  type: string
                sslPort:
                  type: string
              type: object
            secret:
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdVrouter = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: vrouters.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Vrouter
    listKind: VrouterList
    plural: vrouters
    singular: vrouter
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Vrouter is the Schema for the vrouters API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: VrouterSpec is the Spec for the cassandras API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: VrouterConfiguration is the Spec for the cassandras API.
              properties:
                cassandraInstance:
                  type: string
                clusterRole:
                  type: string
                clusterRoleBinding:
                  type: string
                containers: {}
                controlInstance:
                  type: string
                distribution:
                  type: string
                gateway:
                  type: string
                metaDataSecret:
                  type: string
                nodeManager:
                  type: boolean
                physicalInterface:
                  type: string
                serviceAccount:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              properties:
                analyticsPort:
                  type: string
                apiPort:
                  type: string
                collectorPort:
                  type: string
                redisPort:
                  type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

var crdWebui = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: webuis.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Webui
    listKind: WebuiList
    plural: webuis
    singular: webui
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Webui is the Schema for the webuis API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: WebuiSpec is the Spec for the cassandras API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: WebuiConfiguration is the Spec for the cassandras API.
              properties:
                cassandraInstance:
                  type: string
                clusterRole:
                  type: string
                clusterRoleBinding:
                  type: string
                containers: {}
                serviceAccount:
                  type: string
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
	storage: true`

var crdZookeeper = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: zookeepers.contrail.juniper.net
spec:
  group: contrail.juniper.net
  names:
    kind: Zookeeper
    listKind: ZookeeperList
    plural: zookeepers
    singular: zookeeper
  scope: ""
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Zookeeper is the Schema for the zookeepers API.
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ZookeeperSpec is the Spec for the zookeepers API.
          properties:
            commonConfiguration:
              description: CommonConfiguration is the common services struct.
              properties:
                activate:
                  description: Activate defines if the service will be activated by
                    Manager.
                  type: boolean
                create:
                  description: Create defines if the service will be created by Manager.
                  type: boolean
                hostNetwork:
                  description: Host networking requested for this pod. Use the host's
                    network namespace. If this option is set, the ports that will
                    be used must be specified. Default to false.
                  type: boolean
                imagePullSecrets:
                  description: ImagePullSecrets is an optional list of references
                    to secrets in the same namespace to use for pulling any of the
                    images used by this PodSpec.
                  items:
                    type: string
                  type: array
                nodeSelector:
                  additionalProperties:
                    type: string
                  description: 'NodeSelector is a selector which must be true for
                    the pod to fit on a node. Selector which must match a node''s
                    labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.'
                  type: object
                replicas:
                  description: Number of desired pods. This is a pointer to distinguish
                    between explicit zero and not specified. Defaults to 1.
                  format: int32
                  type: integer
                tolerations:
                  description: If specified, the pod's tolerations.
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
              type: object
            serviceConfiguration:
              description: ZookeeperConfiguration is the Spec for the zookeepers API.
              properties:
                clientPort:
                  type: integer
                containers: {}
                electionPort:
                  type: integer
                serverPort:
                  type: integer
              type: object
          required:
          - commonConfiguration
          - serviceConfiguration
          type: object
        status:
          description: ZookeeperStatus defines the status of the zookeeper object.
          properties:
            active:
              description: 'INSERT ADDITIONAL STATUS FIELD - define observed state
                of cluster Important: Run "operator-sdk generate k8s" to regenerate
                code after modifying this file Add custom validation using kubebuilder
                tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html'
              type: boolean
            nodes:
              additionalProperties:
                type: string
              type: object
            ports:
              description: ZookeeperStatusPorts defines the status of the ports of
                the zookeeper object.
              properties:
                clientPort:
                  type: string
              type: object
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true`

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
            cassandra:
              image: cassandra:3.11.4
            init:
              image: python:alpine
            init2:
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
            analyticsapi:
              image: opencontrailnightly/contrail-analytics-api:1910-latest
            api:
              image: opencontrailnightly/contrail-controller-config-api:1910-latest
            collector:
              image: opencontrailnightly/contrail-analytics-collector:1910-latest
            devicemanager:
              image: opencontrailnightly/contrail-controller-config-devicemgr:1910-latest
            init:
              image: python:alpine
            init2:
              image: busybox
            nodeinit:
              image: opencontrailnightly/contrail-node-init:1910-latest
            redis:
              image: redis:4.0.2
            schematransformer:
              image: opencontrailnightly/contrail-controller-config-schema:1910-latest
            servicemonitor:
              image: opencontrailnightly/contrail-controller-config-svcmonitor:1910-latest
          logLevel: SYS_DEBUG
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
            control:
              image: opencontrailnightly/contrail-controller-control-control:1910-latest
            dns:
              image: opencontrailnightly/contrail-controller-control-dns:1910-latest
            init:
              image: python:alpine
            named:
              image: opencontrailnightly/contrail-controller-control-named:1910-latest
            nodeinit:
              image: opencontrailnightly/contrail-node-init:1910-latest
            statusmonitor:
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
            init:
              image: python:alpine
            kubemanager:
              image: michaelhenkel/contrail-kubernetes-kube-manager:1910-latest
            nodeinit:
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
            init:
              image: python:alpine
            provisioner:
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
            init:
              image: python:alpine
            rabbitmq:
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
            init:
              image: python:alpine
            nodeinit:
              image: opencontrailnightly/contrail-node-init:1910-latest
            vrouteragent:
              command: ["bash","-c","/usr/bin/contrail-vrouter-agent --config_file /etc/mycontrail/vrouter.${POD_IP}"]
              image: opencontrailnightly/contrail-vrouter-agent:1910-latest
            vroutercni:
              image: michaelhenkel/contrailcni:v0.0.1
            vrouterkernelbuildinit:
              image: opencontrailnightly/contrail-vrouter-kernel-build-init:1910-latest
            vrouterkernelinit:
              image: opencontrailnightly/contrail-vrouter-kernel-init:1910-latest
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
            init:
              image: python:alpine
            nodeinit:
              image: opencontrailnightly/contrail-node-init:1910-latest
            vrouteragent:
              image: opencontrailnightly/contrail-vrouter-agent:1910-latest
            vroutercni:
              image: michaelhenkel/contrailcni:v0.0.1
            vrouterkernelbuildinit:
              image: opencontrailnightly/contrail-vrouter-kernel-build-init:1910-latest
            vrouterkernelinit:
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
            init:
              image: python:alpine
            nodeinit:
              image: opencontrailnightly/contrail-node-init:1910-latest
            redis:
              image: redis:4.0.2
            webuijob:
              image: opencontrailnightly/contrail-controller-webui-job:1910-latest
            webuiweb:
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
            init:
              image: python:alpine
            zookeeper:
              image: docker.io/zookeeper:3.5.5`
