
all:
	docker-compose -f demo/compose.yaml -f demo/compose-sourcehub.yaml down -v
	docker-compose -f demo/compose.yaml -f demo/compose-sourcehub.yaml up

## Don't spin up sourcehub, use the one that's already running on the host.
no-sourcehub:
	docker-compose -f demo/compose.yaml down -v
	docker-compose -f demo/compose.yaml up
