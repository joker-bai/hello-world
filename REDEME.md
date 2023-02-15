# Hello World

This guide will walk you through deploying a MySQL database and a Hello World application with Kubernetes. 

## Get Started

### MySQL Deploy
1. Create the mysql.yaml file with the following contents:
```yaml
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: mysql
 spec:
   ports:
     - port: 3306
   selector:
     app: mysql
   clusterIP: None
 ---
 apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: mysql
 spec:
   selector:
     matchLabels:
       app: mysql
   strategy:
     type: Recreate
   template:
     metadata:
       labels:
         app: mysql
     spec:
       containers:
         - image: mysql:5.7
           name: mysql
           env:
             - name: MYSQL_ROOT_PASSWORD
               value: root
             - name: MYSQL_DATABASE
               value: hello-world
           ports:
             - containerPort: 3306
               name: mysql
           volumeMounts:
             - name: mysql-persistent-storage
               mountPath: /var/lib/mysql
       volumes:
         - name: mysql-persistent-storage
           emptyDir: {}
```

2. Apply the MySQL deployment and service YAML file.
```bash
kubectl apply -f mysql-deployment.yaml 
```

### Application Deploy
Next, we will create the YAML file for the Hello World application. This file will contain the deployment, service, and ingress for the application.

1. Create the application YAML file:
```yaml
 apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: helloworld
 spec:
   selector:
     matchLabels:
       app: helloworld
   replicas: 3
   template:
     metadata:
       labels:
         app: helloworld
     spec:
       containers:
         - name: helloworld
           image: baidjay/hello-world:v2.0
           imagePullPolicy: Always
           env:
             - name: MYSQL_USER
               value: root
             - name: MYSQL_PASSWORD
               value: root
             - name: MYSQL_ADDRESS
               value: mysql
             - name: MYSQL_DBNAME
               value: hello-world
             - name: CONTENT
               value: Hello World
           ports:
             - containerPort: 8080
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: helloworld
 spec:
   selector:
     app: helloworld
   ports:
     - protocol: TCP
       port: 80
       targetPort: 8080
 ---
 apiVersion: networking.k8s.io/v1
 kind: Ingress
 metadata:
   name: helloworld-ingress
 spec:
   rules:
     - host: helloworld.local
       http:
         paths:
           - path: /
             pathType: Prefix
             backend:
               service:
                 name: helloworld-service
                 port:
                   name: http
```

2. Apply the YAML file:
```shell
kubectl apply -f helloworld.yaml
```

Access the application by opening your web browser and navigating to http://helloworld.local.