
.PHONY: run teardown

run:
	export GO_CONFIG_ENV=prod
	export TENCENT_SECRET_KEY=66weGkLm2AEBtPo0ymSzRo5B8wrCOQFw 
	go env -w GO111MODULE='auto'
	go env -w GOPROXY=https://goproxy.cn,direct
	go run cmd/server/main.go


teardown:
	-PID=$$(lsof -t -i :9197); \
	if [ -n "$$PID" ]; then \
		kill -9 $$PID; \
		echo "Force killed PID $$PID"; \
		sleep 2; \
	else \
		echo "No process listening on port 9197"; \
	fi

# curl -X POST 127.0.0.1:9197/knano/v1/upload/record -H "Content-Type: application/json" -d '[{"upload_id": "test.png", "status": 1,"file_path": "test.png"},{"upload_id": "test.png", "status": 0,"file_path": "test.png"}]'


# curl -X GET "http://127.0.0.1:9197/knano/v1/upload/token?application_code=imagine_hub"

