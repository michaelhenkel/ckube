package kuberesources

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
  scope: Namespaced
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
                containers:
                  items:
                    description: Container defines name, image and command.
                    properties:
                      command:
                        items:
                          type: string
                        type: array
                      image:
                        type: string
                      name:
                        type: string
                    type: object
                  type: array
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
    storage: true
`
