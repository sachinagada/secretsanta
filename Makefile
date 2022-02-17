image:
	docker build -t secretsanta:v0.1.0 .

run:
	docker run -p 8080:8080 -p 3000:3000 secretsanta

minikube:
	minikube start --driver=docker

deploy:
	@eval $$(minikube docker-env) ;\
	docker build -t secretsanta:v0.1.0 . ;\
	kubectl apply -f ./k8s/deployment.yml

# apply the service
service:
	kubectl apply -f ./k8s/service.yml

connect-service:
	minikube service secretservice -n secret-santa

stop:
	minikube stop
