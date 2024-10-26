.PHONY: nuke publisher consumer back run

nuke:
	docker container rm -f mq
	lsof -i:8080 -t | xargs kill -9
	lsof -i:8081 -t | xargs kill -9
	lsof -i:8082 -t | xargs kill -9
	lsof -i:3333 -t | xargs kill -9
	lsof -i:5173 -t | xargs kill -9
	
publisher:
	@export $$(cat .env | xargs); \
	PUBLISHER_PORT=8080 go run cmd/publish/main.go & \
	PUBLISHER_PORT=8081 go run cmd/publish/main.go & \
	PUBLISHER_PORT=8082 go run cmd/publish/main.go & \

consumer:
	CONSUMER_PORT=3333 go run cmd/consume/main.go & \


back: 
	make publisher; \
	make consumer; \

front:
	cd ./queue; \
	npm run dev &


wait_mq:
	@echo "Waiting for RabbitMQ to start..."
	@while ! docker logs mq --tail 1 | grep -q "Time to start RabbitMQ"; do \
		sleep 1; \
	done; \
	echo "RabbitMQ has started."

mq_up:
	docker-compose up -d; \
	make wait_mq; \
	./scripts/init_mq;
	


run: nuke mq_up
	@export $$(cat .env | xargs); \
	make back; \
	make front; \
