TAG := latest
PROJECT_NAME := hack
SVC_NAME := task-etcd

build:
	@TAG=$(TAG) docker-compose -p $(PROJECT_NAME)  build 

run: 
	docker-compose -p $(PROJECT_NAME) up $(SVC_NAME) 

stop:
	@docker-compose stop
