# Kubernetes-Hasher

Calculates the hash sum of files in algorithms ( **MD5, SHA256, SHA1, SHA224, SHA384, SHA512**).

+ we used standard libraries "crypto/sha1", "crypto/sha256","crypto/sha512"
+ you can see https://pkg.go.dev/crypto

## :hammer: Installation



## Installation WEBHOOK
```
cfssl gencert -initca ./webhook/tls/ca-csr.json | cfssljson -bare /tmp/ca

cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./webhook/tls/ca-config.json \
-hostname="tcpdump-webhook,tcpdump-webhook.default.svc.cluster.local,tcpdump-webhook.default.svc,localhost,127.0.0.1" \
-profile=default \
./webhook/tls/ca-csr.json | cfssljson -bare /tmp/tcpdump-webhook

mv /tmp/tcpdump-webhook.pem ./webhook/ssl/tcpdump.pem
mv /tmp/tcpdump-webhook-key.pem ./webhook/ssl/tcpdump.key

Update ConfigMap data in the manifest/webhook-deployment.yaml file with your key and certificate.
cat ./webhook/ssl/tcpdump.key | base64 | tr -d '\n'
cat ./webhook/ssl/tcpdump.pem | base64 | tr -d '\n'

Update caBundle value in the manifest/webhook-configuration.yaml file with your base64 encoded CA certificate.
cat /tmp/ca.pem | base64 | tr -d '\n'

eval $(minikube docker-env)
docker build -t webhook -f webhook/Dockerfile .
docker build -t hasher .

kubectl apply -f manifests/webhook/webhook-deployment.yaml
kubectl apply -f manifests/webhook/webhook-configuration.yaml

kubectl apply -f manifests/hasher/test-deploy.yaml
```