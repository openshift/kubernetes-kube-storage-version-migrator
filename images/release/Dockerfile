FROM registry.svc.ci.openshift.org/ocp/builder:rhel-8-golang-openshift-4.6 AS builder
WORKDIR /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator
COPY . .
RUN make build

FROM registry.svc.ci.openshift.org/ocp/4.6:base
COPY --from=builder /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator/migrator /usr/bin/
