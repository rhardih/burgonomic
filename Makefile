deploy: build push
	docker stack deploy -c docker-compose.yml burgonomic

build:
	docker compose -f docker-compose.yml build

push:
	docker compose -f docker-compose.yml push
