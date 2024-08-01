tidy:
	go mod tidy

MAIN_FILE_PATH=cmd/main.go
CONFIG_FILE_PATH=examples/config.yml

run:
	CONFIG_FILE_PATH=${CONFIG_FILE_PATH} go run ${MAIN_FILE_PATH}
