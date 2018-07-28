FROM alpine:3.7
LABEL Author=gambol99@gmail.com \
      Help=https://github.com/gambol99/kube-admission/issues \
      Name=kube-admission \
      Release=https://github.com/gambol99/kube-admission \
      Url=https://github.com/gambol99/kube-admission

RUN apk add --no-cache ca-certificates && \
    adduser -D controller

ADD bin/kube-admission /kube-admission

USER 1000

ENTRYPOINT [ "kube-admission" ]
