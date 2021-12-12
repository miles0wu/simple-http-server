# prometheus scrape and grafana dashboard setting

## install loki-stack
### download loki-stack
```
$ helm pull grafana/loki-stack --version=2.4.1
$ tar -xvf loki-stack-2.4.1.tgz
```

### replace all `rbac.authorization.k8s.io/v1beta1` with `rbac.authorization.k8s.io/v1` if your k8s version is too new
```
$ cd loki-stack
$ find . -type f -name "*.yaml" -print0 | xargs -0 sed -i s#rbac.authorization.k8s.io/v1beta1#rbac.authorization.k8s.io/v1#g
```
### install loki
```
$ helm upgrade --install loki ./loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

## configure simple-http-server scrape annonation
* add annonations to deployment.yaml
    ```yaml
    ...
    template:
      metadata:
        annonations:
          promethus.io/scrape: "true"
          promethus.io/port: 80
        lables:
          app: simple-http-server
    ...
    ```

* apply the config
    ```
    $ kubectl apply -f deployment.yaml
    ```

## make prometheus a quick test with NodePort
* make sure you can access prometheus with NodePort
    ```
    $ kubctl get svc -l app=prometheus
    NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
    loki-prometheus-alertmanager    ClusterIP   10.104.60.154    <none>        80/TCP         109m
    loki-prometheus-node-exporter   ClusterIP   None             <none>        9100/TCP       109m
    loki-prometheus-pushgateway     ClusterIP   10.110.173.5     <none>        9091/TCP       109m
    loki-prometheus-server          NodePort    10.104.229.240   <none>        80:32457/TCP   109m
    ```

* access prometheus and test with query `simplehttpserver_execution_latency_seconds_bucket`
![prometheus](https://imgur.com/A7utNo9.jpg)

## access grafana and import dashboard
* access grafana with NodePort
    ```
    $ kubectl get svc -lapp.kubernetes.io/name=grafana
    NAME           TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
    loki-grafana   NodePort   10.109.166.12   <none>        80:30035/TCP   119m
    ```

* get admin password from secret loki-grafana
    ```
    $ kubectl get secret loki-grafana -oyaml
    apiVersion: v1
    data:
    admin-password: {base64_encode_password}
    admin-user: {base64_encode_username}
    ldap-toml: ""
    kind: Secret
    ...

    $ echo $(base64_encode_password) | base64 -d
    {admin_password}
    ```

* login grafana and follow the step to import dashboard
    * Click the icon `+`(create)
    * Import
    * Import via panel json
    * Paste `simplehttpserver-dashboard.json` content
    * Load

* Now you have a dashbaord to monitor http server latency
![dashboard](https://imgur.com/DcmJauX.jpg)