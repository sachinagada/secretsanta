![Go](https://github.com/sachinagada/secretsanta/workflows/Go/badge.svg)

## Purpose
This application makes it easier for friends and family to continue the tradition of Secret Santa near the holiday season despite living far apart or quarantining during the pandemic. It provides a form to enter information about the people who want to participate in the tradition, assigns secret santas for all the participants, and emails them their assigned match.

## Installation:

Secret Santa requires go 1.17 or later

`go get github.com/sachinagada/secretsanta`

## Running the application

The application can be run locally with

```
go run ./cmd
```
followed by any flags that need to be passed in. All supported flags can be found with:

```
go run ./cmd --help
```

`mail-username` and `mail-password` need to be passed in either via flags or
environment variables in order to send emails to the chosen secret santas.

## Running with Kubernetes

The application can be deployed to a kubernetes environment (either running locally on minikube or on a cloud platform). Using [kubernetes secrets](https://kubernetes.io/docs/concepts/configuration/secret/) to hold the username and password for the email configuration is highly encouraged.

To run locally on minikube, run the following commands:
```
// run minikube with the docker driver
make minikube

// the following configures the local environment to re-use the Docker daemon inside the Minikube
// instances and builds the docker image within minkube so it can use the local image. It then deploys
// the application on minikube in the `secret-santa` name space
make deploy
```

A service can be added on top of the kubernetes deployment once the deployment is up and running:
```
make service
```

To access the service through minikube, run the following command:
```
make connect-service
```

This will open the default browser to an endpoint (127.0.0.1 with a randomly chosen port) that tunnels back to the service.

## Usage

After running the application locally, go to http://localhost:3000/static (replace with the endpoint from minikube service if running on minikube) to bring up the form. Insert the names and the corresponding email addresses and hit the submit button when done. The application will randomly assign the Secret Santa to each name and email everyone their paired match.

Metrics can be found at http://localhost:8080/metrics.
Go profiles can be found at http://localhost:8080/debug/pprof
Traces can be found at http://localhost:8080/debug/tracez

Note that port 8080 isn't exposed by the service so in order to hit that endpoint, port-forwarding from the pod will be necessary. That can be done with:
```
kubectl -n secret-santa port-forward PODNAME 8080:8080
```
