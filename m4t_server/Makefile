protos:
	@echo "===更新生成 m4t.proto python 端代码==="
	@mkdir -p pb
	@python -m grpc_tools.protoc --proto_path=../api/m4t/ --python_out=./pb --pyi_out=./pb --grpc_python_out=./pb ../api/m4t/m4t.proto

clean:
	@echo "clean protos"
	@rm -rf pb/* 

# docker:


docker_cn:
	@docker buildx build --build-arg BUILD_COUNTRY="CN" -t m4t-server:0.0.1 ./