.PHONY: run
run:
	@echo "启动 Go 服务..."
	@go run cmd/main.go

.PHONY: build
build:
	@echo "构建二进制文件..."
	@go build -o bin/llm-center cmd/main.go

.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf bin/