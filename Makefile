image:
	docker build -t secretsanta:v0.1.0 .

run:
	docker run -p 8080:8080 -p 3000:3000 secretsanta

minikube:
	minikube start --driver=docker

eval-minikube:
	@eval $$(minikube docker-env)


# to deploy to minikube, run:
# 1. minikube start --driver=docker
# 2. eval $(minikube docker-env)
# 3. make image (from the same terminal as step 2)
# 4. make deploy
# 5. make service to create a service once the deployment is up and running
deploy:
	kubectl apply -f ./k8s/deployment.yml

# apply the service
service:
	kubectl apply -f ./k8s/service.yml

stop:
	minikube stop
