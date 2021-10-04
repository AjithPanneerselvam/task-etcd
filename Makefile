build: 
	docker build -t code_task-etcd .

run: 
	docker-compose -p code up task-etcd 
