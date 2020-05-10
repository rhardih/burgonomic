build:
	docker-compose -f docker-compose.yml build

push:
	docker-compose -f docker-compose.yml push

deploy: build push
	docker-compose -f docker-compose.yml -f docker-compose.production.yml up -d
