FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS builder
WORKDIR /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator
COPY . .
RUN make build

FROM registry.ci.openshift.org/ocp/4.19:base-rhel9
COPY --from=builder /go/src/github.com/kubernetes-sigs/kube-storage-version-migrator/migrator /usr/bin/
