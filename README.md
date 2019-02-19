# Auth Web

The auth web is a dashboard for the platform auth-srv. 

## Dependence on Service
- [auth-srv](https://github.com/microhq/auth-srv)

## Getting started

1. Install Consul

	Consul is the default registry/auth for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -dev -advertise=127.0.0.1
	```

3. Download and start the service

	```shell
	go get github.com/micro/auth-web
	auth-web
	```

	OR as a docker container

	```shell
	docker run microhq/auth-web --registry_address=YOUR_REGISTRY_ADDRESS
	```

