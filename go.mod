module github.com/michaelhenkel/ckube

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/Azure/go-autorest v8.0.0+incompatible // indirect
	github.com/Juniper/contrail-operator v0.0.0-20200122001444-1b19a3665b33
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5
	github.com/Microsoft/hcsshim v0.8.7 // indirect
	github.com/ScaleFT/sshkeys v0.0.0-20181112160850-82451a803681 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/containerd v1.3.0
	github.com/creack/goselect v0.0.0-20180501195510-58854f77ee8d // indirect
	github.com/dchest/bcrypt_pbkdf v0.0.0-20150205184540-83f37f9c154a // indirect
	github.com/docker/cli v0.0.0-20200227165822-2298e6a3fe24
	github.com/docker/compose-on-kubernetes v0.4.24 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20181221150755-2cb26cfe9cbf
	github.com/docker/go v1.5.1-1.0.20160303222718-d30aec9fd63c // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-p9p v0.0.0-20170223181108-87ae8514a3a2 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-ini/ini v1.27.3-0.20170519023713-afbc45e87f3b // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gax-go v0.0.0-20170305230405-8c5154c0fe5b // indirect
	github.com/gophercloud/utils v0.0.0-20181029231510-34f5991525d1 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/go-ps v0.0.0-20170309133038-4fdf99ab2936 // indirect
	github.com/moby/datakit v0.0.0-20170703142523-97b3d2305353 // indirect
	github.com/moby/hyperkit v0.0.0-20180416161519-d65b09c1c28a
	github.com/moby/vpnkit v0.1.2-0.20171107134956-0e4293bb1058
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/moul/anonuuid v1.1.0 // indirect
	github.com/moul/gotty-client v1.7.1-0.20180526075433-e5589f6df359 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/runc v1.0.0-rc5.0.20180615140650-ad0f5255060d // indirect
	github.com/opencontainers/runtime-spec v1.0.1
	github.com/packethost/packngo v0.1.1-0.20171201154433-f1be085ecd6f // indirect
	github.com/prometheus/client_golang v1.4.1 // indirect
	github.com/radu-matei/azure-sdk-for-go v5.0.0-beta.0.20161118192335-3b1282355199+incompatible // indirect
	github.com/radu-matei/azure-vhd-utils v0.0.0-20170531165126-e52754d5569d // indirect
	github.com/renstrom/fuzzysearch v1.0.1-0.20180302113537-7a8f9a1c4bed // indirect
	github.com/rn/iso9660wrap v0.0.0-20171120145750-baf8d62ad315
	github.com/scaleway/scaleway-sdk-go v0.0.0-20190617160902-20b731586975 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/surma/gocpio v1.0.2-0.20160926205914-fcb68777e7dc
	github.com/theupdateframework/notary v0.6.0
	github.com/xeipuuv/gojsonschema v1.1.0
	github.com/zchee/go-vmnet v0.0.0-20161021174912-97ebf9174097
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/api v0.6.1-0.20190607001116-5213b8090861
	gopkg.in/dancannon/gorethink.v3 v3.0.5 // indirect
	gopkg.in/fatih/pool.v2 v2.0.0 // indirect
	gopkg.in/gorethink/gorethink.v3 v3.0.5 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools/v3 v3.0.2 // indirect
	k8s.io/api v0.0.0
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v12.0.0+incompatible
	//k8s.io/api v0.17.3 // indirect
	//k8s.io/client-go v11.0.0+incompatible // indirect
	//k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab // indirect
	vbom.ml/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm
