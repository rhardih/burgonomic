docker.sock:
	tunnel-docker.sh

local-build:
	docker-compose -f docker-compose.yml build

local-push:
	docker-compose -f docker-compose.yml push

remote-pull: export DOCKER_HOST=unix://${PWD}/docker.sock
remote-pull: docker.sock
	docker-compose -f docker-compose.yml pull

deploy: export DOCKER_HOST=unix://${PWD}/docker.sock
deploy: docker.sock local-build local-push remote-pull
		docker-compose -f docker-compose.yml -f docker-compose.production.yml up -d;
		tunnel-docker.sh close
