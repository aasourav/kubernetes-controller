kind create cluster --name aascluster --config ./kind-metallb-manifests/kindCreateClusterManifest.yaml
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml
kubectl wait --namespace metallb-system \
                    --for=condition=ready pod \
                    --selector=app=metallb \
                    --timeout=390s
kubectl create -f ./kind-metallb-manifests/metallb.yaml
make install-ingress-controller