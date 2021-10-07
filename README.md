# simple-http-server
Simple http server handle request with following rule:
* Response with header from all request headers
* Response with header VERSION get environment variable from OS
* Logging request ip and response code to stdout

## Getting Started
### install from source
```
$ make build
```

### build image and run with container
```
$ make docker
```

### run with docker
```
$ docker run -d -p 8080:8080 simple-http-server:latest --name simple-http-server
$ docker logs -f simple-http-server
```

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