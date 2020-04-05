# Up And Down The River (Merde!)

A simple, silly, and fun card game!



## How to run locally

Get set up with [golang on your platform](https://golang.org/dl/)

Clone the repo:
```
git clone git@github.com:mattfenwick/upanddowntheriver.git
```

Run the server:

```
cd upanddowntheriver/cmd/server

go run server.go conf.json
```

[Visit the UI](http://localhost:5932/main.html)



## How to deploy to kubernetes

Prereqs: `kubectl` set up and talking to a kubernetes cluster

Run the deploy script:

```
cd upanddowntheriver/deploy

./deploy.sh my-fave-namespace
```

Figure out how to expose the `up-and-down-the-river` service!



## Components

Server: [golang](./cmd/server/server.go)

UI: [javascript](./cmd/server/ui)