NS=$1

kubectl create ns $NS

kubectl create -n $NS -f server.yaml
