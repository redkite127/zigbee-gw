# BUILDER
FROM golang:1.10 AS builder
ARG SERVICE=zigbee-gw
WORKDIR /go/src/$SERVICE/src
COPY src .
RUN go get -insecure -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../$SERVICE -v

# RUNNER
FROM alpine:3.7
ARG SERVICE=zigbee-gw
RUN apk --no-cache --update upgrade add \
    ca-certificates \
    curl
ENV bin_dir /opt/zigbee-gw/bin
ENV etc_dir /opt/zigbee-gw/etc
ENV var_dir /opt/zigbee-gw/var

RUN mkdir -p ${bin_dir} && mkdir -p ${etc_dir} && mkdir -p ${var_dir}

WORKDIR ${bin_dir}

#COPY var/ ${var_dir}/
COPY --from=builder /go/src/$SERVICE/$SERVICE .
RUN chmod +x $SERVICE

CMD ["./zigbee-gw"]

#OLD DOCKERFILE
# FROM alpine:3.6

# ENV bin_dir /opt/zigbee-gw/bin
# ENV etc_dir /opt/zigbee-gw/etc
# ENV var_dir /opt/zigbee-gw/var

# RUN mkdir -p ${bin_dir} && mkdir -p ${etc_dir} && mkdir -p ${var_dir}

# COPY zigbee-gw ${bin_dir}/zigbee-gw

# RUN chmod +x ${bin_dir}/zigbee-gw

# WORKDIR ${bin_dir}

# # it does accept the variable ${etc_dir} in the parameters
# #CMD ["./zigbee-gw", "-config-dir", "/opt/tadaweb/etc"]
# CMD ["./zigbee-gw"]
