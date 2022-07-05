# Kubernetes-Hasher

The implementation of a hasher in golang, which calculates the checksum of files using different algorithms in kubernetes. ( **MD5, SHA256, SHA1, SHA224, SHA384, SHA512**).

+ we used standard libraries "crypto/sha1", "crypto/sha256","crypto/sha512"
+ you can see https://pkg.go.dev/crypto

A Kubernetes sidecar to watch for file changes and restart deployments and pods.
## Getting Started

See examples in manifests/hasher directory for how to add the hasher-sidecar to any pod, and the service account needed.
### Running locally
The code only works running inside a pod in Kubernetes
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using `minikube`.
## Configuration
To work properly, you first need to set the configuration files:
+ environmental variables in the `.env` file
+ config in file `manifests/hasher/configMap.yaml`
+ secret for database `manifests/database/postgres-secret.yaml`
## :hammer: Installing / Deploying

### Installation DATABASE
Apply all annotations in directory "manifests/db/..":
```
kubectl apply -f manifests/db/postgres-db-pv.yaml
kubectl apply -f manifests/db/postgres-db-pvc.yaml
kubectl apply -f manifests/db/postgres-secret.yaml
kubectl apply -f manifests/db/postgres-db-deployment.yaml
kubectl apply -f manifests/db/postgres-db-service.yaml
```

### Installation WEBHOOK
Generate ca in /tmp :
```
cfssl gencert -initca ./webhook/tls/ca-csr.json | cfssljson -bare /tmp/ca
```

Generate private key and certificate for SSL connection:
```
cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./webhook/tls/ca-config.json \
-hostname="k8s-webhook-injector,k8s-webhook-injector.default.svc.cluster.local,k8s-webhook-injector.default.svc,localhost,127.0.0.1" \
-profile=default \
./webhook/tls/ca-csr.json | cfssljson -bare /tmp/k8s-webhook-injector
```

Move your SSL key and certificate to the ssl directory:
```
mv /tmp/k8s-webhook-injector.pem ./webhook/ssl/k8s-webhook-injector.pem
mv /tmp/k8s-webhook-injector-key.pem ./webhook/ssl/k8s-webhook-injector.key
```

Update configuration data in the manifests/webhook/webhook-configMap.yaml file with your key in the appropriate field `data:server.key` and certificate in the appropriate field `data:server.crt:`:
```
cat ./webhook/ssl/k8s-webhook-injector.key | base64 | tr -d '\n'
cat ./webhook/ssl/k8s-webhook-injector.pem | base64 | tr -d '\n'
```

Update field `caBundle` value in the manifests/webhook/webhook-configuration.yaml file with your base64 encoded CA certificate:
```
cat /tmp/ca.pem | base64 | tr -d '\n'
```
Build docker images webhook and hasher:
```
eval $(minikube docker-env)
docker build -t webhook -f webhook/Dockerfile .
docker build -t hasher .
```
Apply webhook annotation:
```
kubectl apply -f manifests/webhook/webhook-configMap.yaml
kubectl apply -f manifests/webhook/webhook-deployment.yaml
kubectl apply -f manifests/webhook/webhook-service.yaml
kubectl apply -f manifests/webhook/webhook-configuration.yaml
```
Apply hasher annotation:
```
kubectl apply -f manifests/hasher/service-account-hasher.yaml
kubectl apply -f manifests/hasher/configMap.yaml
```

For example there is manifests/hasher/test-nginx-deploy.yaml DEPLOYMENT files:
```
kubectl apply -f manifests/hasher/test-nginx-deploy.yaml
```
##Pay attention!
If you want to use a hasher-webhook-injector-sidecar, then you need to specify the following data in your deployment:
+ `spec:template:metadata:labels:hasher-webhook-injector-sidecar: "true"`
+ `hasher-webhook-process-name: "your main process name"`

___________________________
### :notebook_with_decorative_cover: Godoc extracts and generates documentation for Go programs
#### Presents the documentation as a web page.
```go
godoc -http=:6060/sha256sum
go doc packge.function_name
```
for example
```go
go doc pkg/api.Result
```

### :mag: Running tests

You need to go to the folder where the file is located *_test.go and run the following command:
```go
go test -v
```

for example
```go
cd ../pkg/api
go test -v
```

### :minidisc: How to create a `bin` for other platforms:
`bin` will be located in internal/core/services/bin
```
bash build.sh
```

### :mag: Running linter "golangci-lint"
```
golangci-lint run
```