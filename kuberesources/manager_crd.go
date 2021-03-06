package kuberesources

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
  scope: Namespaced
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
                            logLevel:
                              type: string
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
                              properties:
                                bGPPeer:
                                  properties:
                                    number:
                                      type: string
                                    up:
                                      type: string
                                  type: object
                                connections:
                                  items:
                                    properties:
                                      name:
                                        type: string
                                      nodes:
                                        items:
                                          type: string
                                        type: array
                                      status:
                                        type: string
                                      type:
                                        type: string
                                    type: object
                                  type: array
                                numberOfRoutingInstances:
                                  type: string
                                numberOfXMPPPeers:
                                  type: string
                                state:
                                  type: string
                                staticRoutes:
                                  properties:
                                    down:
                                      type: string
                                    number:
                                      type: string
                                  type: object
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
    storage: true
`
