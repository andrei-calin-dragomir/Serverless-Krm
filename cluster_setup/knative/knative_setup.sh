#Install Go
curl -OL https://golang.org/dl/go1.23.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.23.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin >> ~/.profile
source ~/.profile

#Install tuf-client

go install github.com/theupdateframework/go-tuf/cmd/tuf-client@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

#Install Cosign

go install github.com/sigstore/cosign/v2/cmd/cosign@latest
curl -o sigstore-root.json https://sigstore.github.io/root-signing/10.root.json
tuf-client init https://tuf-repo-cdn.sigstore.dev sigstore-root.json
tuf-client get https://tuf-repo-cdn.sigstore.dev artifact.pub > artifact.pub
curl -o cosign-release.sig -L https://github.com/sigstore/cosign/releases/download/v2.0.1/cosign-linux-amd64.sig
base64 -d cosign-release.sig > cosign-release.sig.decoded
curl -o cosign -L https://github.com/sigstore/cosign/releases/download/v2.0.1/cosign-linux-amd64
openssl dgst -sha256 -verify artifact.pub -signature cosign-release.sig.decoded cosign

#Install Jq

sudo apt install jq

#Get Knative CRDs

curl -sSL https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-core.yaml \
  | grep 'gcr.io/' | awk '{print $2}' | sort | uniq \
  | xargs -n 1 \
    cosign verify -o text \
      --certificate-identity=signer@knative-releases.iam.gserviceaccount.com \
      --certificate-oidc-issuer=https://accounts.google.com

#Deploy Knative Serving Components

kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-crds.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-core.yaml

#Deploy MetalLB for LoadBalancing

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/main/config/manifests/metallb-native.yaml
kubectl apply -f metallb-config.yaml
kubectl apply -f ipaddresspool.yaml
kubectl apply -f l2advertisement.yaml

#Deploy Kourier Networking Layer (and DNS config)

kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.16.0/kourier.yaml
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.16.0/serving-default-domain.yaml

#Deploy HPA Autoscaler

kubectl apply -f https://github.com/knative/serving/releases/download/knativ