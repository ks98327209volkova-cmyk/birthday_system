.PHONY: up 
up:
	docker-compose up --build -d

.PHONY: down  
down:
	docker-compose down

.PHONY: test 
test:
	docker-compose run --rm gift-manager go test ./...


.PHONY: clean
clean:
	docker-compose down -v
	docker system prune -f