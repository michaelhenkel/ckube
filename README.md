# Prerequisite
- macOS 10.11+
- HyperKit
- DockerForMac (for building operating system images)

## Hyperkit installation
```
brew install hyperkit
```

## DockerForMac installation
```
https://docs.docker.com/docker-for-mac/install/
```

# Installation
## Get cKube binary
```
curl -LO https://github.com/michaelhenkel/ckube/releases/download/v0.1/ckube
chmod +x ckube
mv ckube /usr/local/bin
```
## Build OS images
```
ckube build
```
## Create inital cluster
```
ckube create cluster cluster1
```
## Add Master node to cluster
Nodes can be created using either vpnkit or vmnet networking.    
- vpnkit doesn't expose the VM to the outside
  ssh is needed on the MAC for a reverse tunnel    
- vmnet exposes the VM but requires root permissions

### VPNKIT
```
ckube cluster cluster1 add master -v vpnkit master1
```

### VMNET
```
sudo chown root:wheel /usr/local/bin/ckube && sudo chmod u+s /usr/local/bin/ckube
ckube cluster cluster1 add master master1
```

# Using
## Source kubeconfig
```
export KUBECONFIG=~/.ckube/clusters/cluster1/master/master1/k3s.yaml
```
## Check Contrail services
Wait for Contrail services are fully up and running
```
kubectl -n contrail get pods
NAME                                          READY   STATUS    RESTARTS   AGE
contrail-operator-8487b46d7f-28rhs            1/1     Running   0          14m
zookeeper1-zookeeper-statefulset-0            1/1     Running   0          13m
rabbitmq1-rabbitmq-statefulset-0              1/1     Running   0          13m
cassandra1-cassandra-statefulset-0            1/1     Running   0          13m
config1-config-statefulset-0                  7/7     Running   0          11m
kubemanager1-kubemanager-statefulset-0        1/1     Running   0          7m3s
provmanager1-provisionmanager-statefulset-0   1/1     Running   0          7m3s
webui1-webui-statefulset-0                    3/3     Running   0          7m3s
control1-control-statefulset-0                4/4     Running   0          7m3s
vroutermaster-vrouter-daemonset-gtxcx         2/2     Running   0          5m12s
```
## Create a POD
```
cat << EOF > /tmp/pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: busy1
  annotations:
  labels:
    app: busy1
spec:
  containers:
  - name: busy1
    image: busybox
    command: ["/bin/sh","-c", "while true; do echo hello; sleep 10;done"]
EOF
kubectl apply -f /tmp/pod.yaml
```
## Check the POD
```
kubectl get pods -owide
NAME    READY   STATUS    RESTARTS   AGE   IP              NODE                    NOMINATED NODE   READINESS GATES
busy1   1/1     Running   0          15s   10.47.255.249   linuxkit-025000000001   <none>           <none>
```
## Accessing Contrail Webui
Open browser at https://127.0.0.1:8143
