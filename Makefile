start:
	docker compose up --build

stop:
	docker compose down

test-api:
	./test-polygon-api.sh