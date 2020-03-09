FROM alpine:3.9
RUN wget https://github.com/rancher/k3s/releases/download/v1.17.2%2Bk3s1/k3s && \
      chmod +x k3s && \
      mkdir -p /var/lib/rancher/k3s/agent/etc/cni/net.d && \
      ln -s /etc/cni/net.d /var/lib/rancher/k3s/agent/etc/cni/net.d
VOLUME /var/lib/kubelet
VOLUME /var/lib/rancher/k3s
VOLUME /var/lib/cni
VOLUME /var/log
ENTRYPOINT ["/k3s", "server --cluster-cidr=10.32.0.0/12 --service-cidr=10.96.0.0/12 --no-flannel --no-deploy=traefik"]
