build: 
	docker build -t task-etcd .

run: 
	docker-compose up task-etcd  
