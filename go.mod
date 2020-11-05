module github.com/kubernetes-sigs/kube-storage-version-migrator

go 1.13

require (
	github.com/golang/groupcache v0.0.0-20191002201903-404acd9df4cc // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/openshift/build-machinery-go v0.0.0-20200512074546-3744767c4131
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/client_model v0.2.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.18.6
)

replace (
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/klog => k8s.io/klog v0.1.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.6
)
