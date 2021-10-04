build: 
	docker build -t code_todo .

run: 
	docker-compose -p code up todo 
