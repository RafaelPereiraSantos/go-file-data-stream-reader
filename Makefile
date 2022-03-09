build:
	docker build -t go-file-stream-reder-example .

run:
	docker run --name go-file-stream-reder-example go-file-stream-reder-example

start:
	docker start go-file-stream-reder-example -a

remove:
	docker rm go-file-stream-reder-example
	docker image rm go-file-stream-reder-example
