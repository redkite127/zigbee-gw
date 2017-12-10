FROM alpine:3.6

ENV bin_dir /opt/zigbee-gw/bin
ENV etc_dir /opt/zigbee-gw/etc
ENV var_dir /opt/zigbee-gw/var

RUN mkdir -p ${bin_dir} && mkdir -p ${etc_dir} && mkdir -p ${var_dir}

COPY zigbee-gw ${bin_dir}/zigbee-gw

RUN chmod +x ${bin_dir}/zigbee-gw

WORKDIR ${bin_dir}

# it does accept the variable ${etc_dir} in the parameters
#CMD ["./zigbee-gw", "-config-dir", "/opt/tadaweb/etc"]
CMD ["./zigbee-gw"]
