NS=$1

kubectl create ns $NS

kubectl create -n $NS -f server.yaml

kubectl expose deployment -n $NS up-and-down-the-river --name uad-exposed --type LoadBalancer --port 5932 --target-port 5932
IP=$(kubectl get nodes -o json | jq -r '.items[0].status.addresses | map(select(.type == "ExternalIP"))[0].address')
PORT=$(kubectl get svc -n $NS uad-exposed -o json | jq '.spec.ports[0].nodePort')

echo http://$IP:$PORT/model