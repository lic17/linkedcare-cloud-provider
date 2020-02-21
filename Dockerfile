FROM alpine:3.5

ADD cloud-controller-manager /usr/local/bin/cloud-controller-manager

CMD ["cloud-controller-manager"]
