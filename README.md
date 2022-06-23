# Kubernetes-Hasher

Calculates the hash sum of files in algorithms ( **MD5, SHA256, SHA1, SHA224, SHA384, SHA512**).

+ we used standard libraries "crypto/sha1", "crypto/sha256","crypto/sha512"
+ you can see https://pkg.go.dev/crypto

## :hammer: Installation

## Installation DATABASE
Apply all annotations in directory "manifests/db/..":
```
kubectl apply -f manifests/db/postgres-db-pv.yaml
kubectl apply -f manifests/db/postgres-db-pvc.yaml
kubectl apply -f manifests/db/postgres-secret.yaml
kubectl apply -f manifests/db/postgres-db-deployment.yaml
kubectl apply -f manifests/db/postgres-db-service.yaml
```

## Installation WEBHOOK

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
-hostname="tcpdump-webhook,tcpdump-webhook.default.svc.cluster.local,tcpdump-webhook.default.svc,localhost,127.0.0.1" \
-profile=default \
./webhook/tls/ca-csr.json | cfssljson -bare /tmp/tcpdump-webhook
```

Move your SSL key and certificate to the ssl directory:
```
mv /tmp/tcpdump-webhook.pem ./webhook/ssl/tcpdump.pem
mv /tmp/tcpdump-webhook-key.pem ./webhook/ssl/tcpdump.key
```

Update ConfigMap data in the manifests/webhook/webhook-deployment.yaml file with your key and certificate:
```
cat ./webhook/ssl/tcpdump.key | base64 | tr -d '\n'
cat ./webhook/ssl/tcpdump.pem | base64 | tr -d '\n'
```

Update caBundle value in the manifests/webhook/webhook-configuration.yaml file with your base64 encoded CA certificate:
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
kubectl apply -f manifests/webhook/webhook-deployment.yaml
kubectl apply -f manifests/webhook/webhook-configuration.yaml
```
For example there is DEPLOYMENT file:
```
kubectl apply -f manifests/hasher/test-deploy.yaml
```

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