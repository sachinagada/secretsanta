image:
	docker build -t secretsanta .

run:
	docker run -p 8080:8080 -p 3000:3000 secretsanta
