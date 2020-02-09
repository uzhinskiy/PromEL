# PromEL #

PromEL - A prometheus remote storage adapter for ElasticSearch

## INTRODUCTION ##


## INSTALL ##


## USAGE ##

### Standalone ###

### Docker ###

#### Tuning Linux kernel ####


    $ nano /etc/sysctl.conf
     net.ipv4.tcp_keepalive_time = 300
     net.ipv4.tcp_keepalive_intvl = 20
     net.ipv4.tcp_keepalive_probes = 3
     fs.file-max = 1199136

    $ sysctl -p


![Peak load ~ 12000 doc/sec](https://gitlab.insitu.co.il/boris.uzhinskiy/promel/raw/dev/docs/images/kibana_discovery_state.png)
![Normal load ~ 1000 doc/sec](https://gitlab.insitu.co.il/boris.uzhinskiy/promel/raw/dev/docs/images/kibana_discovery_state_2.png)
