include .env

default: docker docker-up

gen-protoc:
	@echo "***** Gen protoc... *****"
	@protoc -I api/ api/protobuf.proto --go_out=plugins=grpc:api
	@echo "***** Success *****"

docker:
	docker-compose build

docker-up:
	docker-compose up
	docker-compose ps

docker-clean:
	docker-compose down --rmi=all