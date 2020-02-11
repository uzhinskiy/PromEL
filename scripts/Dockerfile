# stage 1: build
FROM golang:latest AS builder
LABEL maintainer="Uzhinskiy Boris <boris.ujinsky@gmail.com>"

# Add source code
RUN mkdir -p /go/src/promel
ADD . /go/src/promel

# Build binary
RUN cd /go/src/promel/ && \ 
    make

# stage 2: lightweight "release"
FROM debian:latest
LABEL maintainer="Uzhinskiy Boris <boris.ujinsky@gmail.com>"

EXPOSE 9090/tcp 9091/tcp

COPY --from=builder /go/src/promel/build/ /usr/local/sbin/
COPY --from=builder /go/src/promel/scripts/promel-docker.yml /etc/promel/promel.yml

ENTRYPOINT [ "/usr/local/sbin/promel" ]
CMD [ "-config", "/etc/promel/promel.yml" ]