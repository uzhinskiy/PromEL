# PromEL - Elasticsearch Prometheus Adapter #

## INTRODUCTION ##
This is an adapter that accepts Prometheus remote read/write requests, and sends them on to Elasticsearch. This allows using Elasticsearch as long term storage for Prometheus.

Requires Elasticsearch v7.0 or greater.


## INSTALL ##

    $ git clone https://github.com/uzhinskiy/PromEL.git promel
    $ cd promel
    $ make

## USAGE ##

### Standalone ###

	$ sudo mkdir /var/log/promel/
	$ sudo mkdir /etc/promel/
	$ sudo cp ./build/promel /usr/local/sbin/promel
	$ sudo cp ./scripts/promel.service /etc/systemd/system/
	$ sudo cp ./scripts/promel-standalone.yml /etc/promel/promel.yml
	$ edit /etc/promel/promel.yml (replace elk0X-ip with your actual IPs of elastic-nodes)
	$ sudo systemctl daemon-reload && systemctl start promel
	$ sudo systemctl enable promel

### Docker ###

    $ git clone https://github.com/uzhinskiy/PromEL.git promel
    $ cd promel
    $ edit /scripts/promel-docker.yml (replace elk0X-ip with your actual IPs of elastic-nodes)
    $ docker build -t promel -f scripts/Dockerfile .
    $ docker run -t -d -p 0.0.0.0:9090:9090 -p 0.0.0.0:9091:9091 --name promel-app promel

#### Tuning Linux kernel ####

    $ sudo nano /etc/sysctl.conf
     net.ipv4.tcp_keepalive_time = 300
     net.ipv4.tcp_keepalive_intvl = 20
     net.ipv4.tcp_keepalive_probes = 3
     fs.file-max = 1199136

    $ sudo sysctl -p


![Peak load ~ 12000 doc/sec](https://raw.githubusercontent.com/uzhinskiy/PromEL/master/docs/images/kibana_discovery_state.png)
![Normal load ~ 1000 doc/sec](https://raw.githubusercontent.com/uzhinskiy/PromEL/master/docs/images/kibana_discovery_state_2.png)
