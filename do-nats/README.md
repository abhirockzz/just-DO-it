# NATS cluster on DigitalOcean

[NATS](https://nats.io/) is a cloud native messaging system with a focus on simplicity and high performance. Let's setup and play around with a **load balanced** NATS **cluster** on [DigitalOcean](https://digitalocean.com)

We will make use of

- [Droplet](https://www.digitalocean.com/products/droplets/) - for NATS
- [Load Balancer](https://www.digitalocean.com/products/load-balancer/) - front end NATS cluster 
- [REST API](https://developers.digitalocean.com/documentation/v2/)  and [Go client](https://github.com/digitalocean/godo) for [programmatic access](https://github.com/abhirockzz/just-DO-it/blob/master/do-nats/cleanup-do-resources.go) to above mentioned resources

## Create Droplets

Start by setting up a personal access token - refer to [this tutorial](https://www.digitalocean.com/docs/api/create-personal-access-token/)

> DigitalOcean has an excellent REST API - we'll leverage it with [cURL](https://curl.haxx.se/), although there are [multiple clients](https://developers.digitalocean.com/libraries/) out there too !

**Create first droplet**

`curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer DO_PERSONAL_ACCESS_TOKEN" -d '{"name":"nats-droplet-1","region":"blr1","size":"s-1vcpu-1gb","image":"ubuntu-18-04-x64","ssh_keys":null,"backups":false,"ipv6":true,"user_data":null,"private_networking":null,"volumes": null,"tags":["nats"]}' "https://api.digitalocean.com/v2/droplets"`

**... and the second one**

`curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer DO_PERSONAL_ACCESS_TOKEN" -d '{"name":"nats-droplet-2","region":"blr1","size":"s-1vcpu-1gb","image":"ubuntu-18-04-x64","ssh_keys":null,"backups":false,"ipv6":true,"user_data":null,"private_networking":null,"volumes": null,"tags":["nats"]}' "https://api.digitalocean.com/v2/droplets"`

You should now have

- Standard Droplets,
- with smallest shape (1 GB, 1vCPU), 
- in the Bangalore (India) data center
- with the `nats` tag

> you can obviously change these and the other parameters in the JSON payload

![](images/droplets.JPG)

## Setup NATS

[SSH into your Droplet](https://www.digitalocean.com/docs/droplets/how-to/connect-with-ssh/) and to begin NATS setup

For first Droplet

- Get NATS - `wget https://github.com/nats-io/gnatsd/releases/download/v1.2.0/gnatsd-v1.2.0-linux-amd64.zip`
- install **unzip** - `apt install unzip`
- `unzip gnatsd-v1.2.0-linux-amd64.zip`
- `cd gnatsd-v1.2.0-linux-amd64`
- start the NATS server - `./gnatsd -p 4222 -m 8222 -cluster nats://DROPLET_IP:5222 &`
	- this is in `clustered` mode
	- the first server (this one) will act as the `seed` node
	- `4222` is the TCP port for NATS clients
	- `8222` is the HTTP monitoring port
	- `5222` is the port for intra-cluster communication


![](images/nats-started.JPG)


On the second Droplet

[SSH into your Droplet](https://www.digitalocean.com/docs/droplets/how-to/connect-with-ssh/)

- Get NATS `wget https://github.com/nats-io/gnatsd/releases/download/v1.2.0/gnatsd-v1.2.0-linux-amd64.zip`
- install unzip `apt install unzip`
- `unzip gnatsd-v1.2.0-linux-amd64.zip`
- `cd gnatsd-v1.2.0-linux-amd64`
- start the second node in the cluster - `./gnatsd -p 4222 -m 8222 -cluster nats://DROPLET_IP:5222 -routes nats://SEED_NODE_IP:5222 &`

> `SEED_NODE_IP` is nothing but the IP of the first droplet i.e. the `seed` NATS node

![](images/nats-2-started.JPG)

## Provision a Load Balancer

We will place a load balancer front ending our NATS cluster. Let's provision a DigitalOcean Load Balancer for the two NATS Droplets and configure it

`curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer DO_PERSONAL_ACCESS_TOKEN" -d '{"name": "nats-lb", "region": "blr1", "forwarding_rules":[{"entry_protocol":"http","entry_port":8222,"target_protocol":"http","target_port":8222,"certificate_id":"","tls_passthrough":false}, {"entry_protocol":"tcp","entry_port":4222,"target_protocol":"tcp","target_port":4222,"certificate_id":"","tls_passthrough":false}], "health_check":{"protocol":"http","port":8222,"path":"/varz","check_interval_seconds":5,"response_timeout_seconds":5,"healthy_threshold":5,"unhealthy_threshold":3}, "sticky_sessions":{"type":"none"}, "tag": "nats"}' "https://api.digitalocean.com/v2/load_balancers"`

- provisioned in the same data center as the droplets (Bangalore in this case)
- uses droplets tag (`"tag": "nats"`) for automatic assignment (no manual steps involved)
- Round Robin load balancing
- Port forwarding rules
	- `4222` over TCP for clients
	- `8222` over HTTP for monitoring
- ping `http://DROPLET_IP:8222/varz` for health check 

![](images/lb.JPG)

**Load Balancer configuration**

![](images/lb-settings.JPG)

## Test the setup

- Access NATS monitoring console - In a browser, open `http://LOAD_BALANCER_IP:8222`

![](images/lb-nats.JPG)

- Access the `/varz` endpoint `http://LOAD_BALANCER_IP:8222/varz` to check general statistics. 

> Notice the `server_id` attribute in the returned JSON payload - it will be different for subsequent requests. This is because we had chose round robin load balancing style and the request will go to a different node each time


		{
		  "server_id": "eXS0FsfA1h7N83OKJCTrmT",
		  "version": "1.2.0",
		  "proto": 1,
		  "go": "go1.10.3",
		  "host": "0.0.0.0",
		  "connect_urls": [
		    "139.59.60.61:4222",
		    "10.47.0.5:4222",
		    "[2400:6180:100:d0::859:7001]:4222",
		    "139.59.35.7:4222",
		    "10.47.0.6:4222",
		    "[2400:6180:100:d0::85d:2001]:4222"
		  ],
		  "addr": "0.0.0.0",
		  "max_connections": 65536,
		  "ping_interval": 120000000000,
		  "ping_max": 2,
		  "http_host": "0.0.0.0",
		  "http_port": 8222,
		  "https_port": 0,
		  "auth_timeout": 1,
		  "max_control_line": 1024,
		  "cluster": {
		    "addr": "139.59.60.61",
		    "cluster_port": 5222,
		    "auth_timeout": 1
		  },
		  "tls_timeout": 0.5,
		  "port": 4222,
		  "max_payload": 1048576,
		  "start": "2018-07-15T10:31:15.8625413Z",
		  "now": "2018-07-15T10:53:05.722699974Z",
		  "uptime": "21m49s",
		  "mem": 9089024,
		  "cores": 1,
		  "cpu": 0,
		  "connections": 0,
		  "total_connections": 4,
		  "routes": 1,
		  "remotes": 1,
		  "in_msgs": 0,
		  "out_msgs": 0,
		  "in_bytes": 0,
		  "out_bytes": 0,
		  "slow_consumers": 0,
		  "max_pending": 268435456,
		  "write_deadline": 2000000000,
		  "subscriptions": 0,
		  "http_req_stats": {
		    "/": 2,
		    "/connz": 0,
		    "/routez": 0,
		    "/subsz": 0,
		    "/varz": 353
		  },
		  "config_load_time": "2018-07-15T10:31:15.8625413Z"
		}

- Access the `/routez` endpoint `http://LOAD_BALANCER_IP:8222/routez` to confirm that the cluster routing is has been configured properly

		{
		  "server_id": "eXS0FsfA1h7N83OKJCTrmT",
		  "now": "2018-07-15T10:54:12.041718718Z",
		  "num_routes": 1,
		  "routes": [
		    {
		      "rid": 5,
		      "remote_id": "Bho7R0xIUop1VG5XslvEJh",
		      "did_solicit": false,
		      "is_configured": false,
		      "ip": "139.59.35.7",
		      "port": 50158,
		      "pending_size": 0,
		      "in_msgs": 0,
		      "out_msgs": 0,
		      "in_bytes": 0,
		      "out_bytes": 0,
		      "subscriptions": 0
		    }
		  ]
		}

## Test the application

This is a simple pub-sub application which used the NATS Java and Go clients 

Set up your local machine

- You need `Docker`, `Docker Compose` and `Maven`
- get the project - `git clone https://github.com/abhirockzz/just-DO-it.git`
- `cd do-nats/nats-pub-sub`
- build the Java publisher service - `mvn clean install -f java_publisher/pom.xml`

To run

- populate the `NATS_SERVER` environment variable in `docker-compose.yml` - specify the load balancer IP e.g. `nats://174.138.123.164:4222`

		version: '3'
		services:
		    nats_server:
		        image: nats
		        container_name: natsio-server
		        ports:
		            - '8222:8222'
		    publisher:
		        build: java_publisher
		        depends_on:
		            - nats_server
		        environment:
		            - NATS_SERVER=nats://174.138.123.164:4222
		    subscriber:
		        build: go_subscriber
		        depends_on:
		            - nats_server
		        environment:
		            - NATS_SERVER=nats://174.138.123.164:4222

- build containers - `docker-compose build --no-cache`

![](images/docker-compose-build.JPG)

- start containers - `docker-compose up`
	- the `pub sub` process will start
	- the client **discovers** both the cluster nodes
	- the publisher randomly connects to one of the NATS cluster node 
	- the subscriber randomly connects to one of the NATS cluster node which might or might not be the same as the one which the publisher is connected to

![](images/docker-compose-up.JPG)

- check cluster monitoring info 
	- `http://LOAD_BALANCER_IP:8222/connz` for connection related stats
	- `http://LOAD_BALANCER_IP:8222/subsz` for info on specific subscription

**Connection stats**

		{
		  "server_id": "eXS0FsfA1h7N83OKJCTrmT",
		  "now": "2018-07-15T11:07:02.901448486Z",
		  "num_connections": 1,
		  "total": 1,
		  "offset": 0,
		  "limit": 1024,
		  "connections": [
		    {
		      "cid": 6,
		      "ip": "159.65.158.166",
		      "port": 30512,
		      "start": "2018-07-15T11:04:30.061962583Z",
		      "last_activity": "2018-07-15T11:07:00.629545897Z",
		      "uptime": "2m32s",
		      "idle": "2s",
		      "pending_bytes": 0,
		      "in_msgs": 31,
		      "out_msgs": 0,
		      "in_bytes": 1240,
		      "out_bytes": 0,
		      "subscriptions": 0,
		      "lang": "java",
		      "version": "1.0"
		    }
		  ]
		}

**subscription info**

		{
		  "num_subscriptions": 1,
		  "num_cache": 1,
		  "num_inserts": 1,
		  "num_removes": 0,
		  "num_matches": 1,
		  "cache_hit_rate": 0,
		  "max_fanout": 1,
		  "avg_fanout": 1
		}

- scale up subscribers - `docker-compose up --scale subscriber=2`
	- this subscriber process will connect to the node which the first subscriber instance is not connected to
	- you should see differences reflecting in the `/connz` and `/subsz` endpoint payloads

![](images/subscriber-2.JPG)

- after you've had enough of this pub sub, clean up using `docker-compose down -v`

![](images/docker-compose-down.JPG)

## Clean up DigitalOcean resources

We will use the DigitalOcean Go client to remove the Droplets and Load Balancer

- You need [Go](https://golang.org) [installed](https://golang.org/doc/install)
- Fetch libraries
	- `go get github.com/digitalocean/godo`
	- `go get golang.org/x/oauth2`
- `go run cleanup-do-resources.go --DO_PERSONAL_ACCESS_TOKEN=<your_DO_token> --DROPLET_TAG=<droplet tag>` e.g. `go run cleanup-do-resources.go --DO_PERSONAL_ACCESS_TOKEN=abcdefgijklmnopqrstuvwxyz12345678910 --DROPLET_TAG=nats`

![](images/delete-do-resources.JPG)

Check your DigitalOcean [dashboard](https://cloud.digitalocean.com/dashboard) to confirm resource(s) deletion