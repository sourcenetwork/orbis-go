build:
	KO_DOCKER_REPO=ko.local ko build -B --platform=linux/$$(go env GOARCH) ./cmd/orbisd

run:
	docker-compose -f demo/compose.yaml -f demo/compose-sourcehub.yaml down -v
	docker-compose -f demo/compose.yaml -f demo/compose-sourcehub.yaml up

## Don't spin up sourcehub, use the one that's already running on the host.
run-no-sourcehub:
	docker-compose -f demo/compose.yaml down -v
	docker-compose -f demo/compose.yaml up
