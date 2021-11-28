## network flow
![flow](https://imgur.com/eny5O7r.jpg)

## How to deploy on kubernetes
### create simple-http-server deployment and service
```
$ kubectl apply -f deployment.yaml
```

### create nginx-ingress controller
```
$ kubectl apply -f nginx-ingress-deployment.yaml
```

### create secret
```
# generate key & cert
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=*/O=foo"

# replace key & cert part in file and create
$ kubectl apply -f secret.yaml
```

### create ingress
```
$ kubectl create -f ingress.yaml
```

### test the result
```
$ k get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.96.180.217   <none>        80:30185/TCP,443:31231/TCP   29m
ingress-nginx-controller-admission   ClusterIP   10.102.79.180   <none>        443/TCP                      29m
```

```
curl -i -k -H "Host:demo.com" https://192.168.10.2:31231
HTTP/2 200
date: Sun, 28 Nov 2021 13:17:20 GMT
content-length: 0
accept: */*
user-agent: curl/7.64.1
version: 1.0.0
x-forwarded-for: 10.0.2.15
x-forwarded-host: demo.com
x-forwarded-port: 443
x-forwarded-proto: https
x-forwarded-scheme: https
x-real-ip: 10.0.2.15
x-request-id: 5dac2d0176ed99760bffc6aa669f4be9
x-scheme: https
strict-transport-security: max-age=15724800; includeSubDomains
```
