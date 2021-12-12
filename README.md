# simple-http-server
Simple http server handle request with following rule:
* Response with header from all request headers
* Response with header VERSION get environment variable from OS
* Logging request ip and response code to stdout

## Getting Started
### install from source
```
# default $ARCH=amd64
$ make build ARCH=arm64
```

### build image and run with container
```
$ make docker
```

### run with docker
```
$ docker run -d -p 8080:8080 --name simple-http-server simple-http-server:latest
$ docker logs -f simple-http-server
```

### pull from dockerhub and run
```
$ docker run -d -P miles0wu/simple-http-server
```

## Image repository
* [miles0wu/simple-http-server](https://hub.docker.com/r/miles0wu/simple-http-server)

## Usage
```
NAME:
   simple-http-server - simple-http-server implement simple response from request

USAGE:
   simple-http-server [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value  HTTP port (default: 8080) [$PORT]
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)
```

### APIs
#### HealthCheck
```
GET /healthz
200 OK
```

### Example
```
$ curl -i -H "test:123" 127.0.0.1:8080/healthz
HTTP/1.1 200 OK
Accept: */*
Test: 123
User-Agent: curl/7.64.1
Version: 1.0.0
Date: Sat, 09 Oct 2021 17:09:57 GMT
Content-Length: 0

$ docker logs -f simple-http-server
2021/10/09 17:09:57 {"URI":"/healthz","IP":"172.17.0.1","Code":200}
```

### Deploy on kubernetes
* [deploy.md](k8s-yaml/deploy.md)