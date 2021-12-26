### install istio
```
$ curl -L https://istio.io/downloadIstio | sh -
$ cd istio-1.12.1/
$ export PATH=$PWD/bin:$PATH
$ istioctl install --set profile=demo
```

### create ns tracing and services
You can confiure simple-http-server env `TRANSFER_ADDRESS` any address you want in yaml like [service0.yaml](service0.yaml)
```
$ kubectl create ns tracing
$ kubectl label ns tracing istio-injection=enabled
$ kubectl -n tracing apply -f service0.yaml
$ kubectl -n tracing apply -f service1.yaml
$ kubectl -n tracing apply -f service2.yaml
```

### install jaeger and configure sampling rate as 100%
```
$ kubectl apply -f jaeger.yaml
$ kubectl edit configmap istio -n istio-system
# set tracing.sampling=100
```

### create tls for https gateway
```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=demo /CN=*' -keyout tls.key -out tls.crt

kubectl create -n istio-system secret tls demo-credential --key=tls.key --cert=tls.crt
```

### create gateway and virtualService
```
kubectl apply -f istio-specs.yaml -n tracing
```

### check ingress gateway ip
```
$k get svc -nistio-system
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)
istio-egressgateway    ClusterIP      10.101.26.62    <none>        80/TCP,443TCP
istio-ingressgateway   LoadBalancer   10.108.109.37   <pending>     15021:31955/TCP,80:31325/TCP,443:32706/TCP,31400:30753/TCP,15443:30837/TCP
istiod                 ClusterIP      10.110.121.46   <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP
jaeger-collector       ClusterIP      10.107.77.128   <none>        14268/TCP,14250/TCP,9411/TCP
tracing                ClusterIP      10.110.90.202   <none>        80/TCP,16685/TCP
zipkin                 ClusterIP      10.99.60.96     <none>        9411/TCP
```

### access service0
```
curl --resolve httpsserver.demo.com:443:10.108.109.37 https://httpsserver.demo.com/service0 -v -k

...
===========Details of the http request header:===========
HTTP/1.1 200 OK
Transfer-Encoding: chunked
Content-Type: text/plain; charset=utf-8
Date: Sun, 26 Dec 2021 17:56:18 GMT
Server: envoy
X-Envoy-Upstream-Service-Time: 1034

372
===========Details of the http request header:===========
HTTP/1.1 200 OK
Content-Length: 655
Content-Type: text/plain; charset=utf-8
Date: Sun, 26 Dec 2021 17:56:18 GMT
Server: envoy
X-Envoy-Upstream-Service-Time: 299

===========Details of the http request header:===========
X-Envoy-Internal=[true]
X-Forwarded-Client-Cert=[By=spiffe://cluster.local/ns/tracing/sa/default;Hash=d529c3c6d0734359a327fdc4c7532f9ed4ed69955f674137437851c9c9367832;Subject="";URI=spiffe://cluster.local/ns/tracing/sa/default]
X-B3-Spanid=[23e878509328c8b3]
X-B3-Parentspanid=[ae2723ae2dda8d11]
X-Envoy-Attempt-Count=[1]
X-Forwarded-For=[10.0.2.15]
X-Forwarded-Proto=[https]
X-B3-Traceid=[18c440c3c6cb90eeabdeaaaa9850cb19]
Accept-Encoding=[gzip,gzip]
User-Agent=[Go-http-client/1.1,Go-http-client/1.1,curl/7.68.0]
X-Request-Id=[7fde97c8-c0a8-9c7e-ac91-41d391a3d8a1]
X-B3-Sampled=[1]
Accept=[*/*]

0

* Connection #0 to host httpsserver.demo.com left intact
```

### check jaeger dashboard
```
istioctl dashboard jaeger
http://localhost:16686
```

### example
![jaeger](https://imgur.com/woT7jwJ.jpg)