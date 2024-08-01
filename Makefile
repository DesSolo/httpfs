tidy:
	go mod tidy

MAIN_FILE_PATH=cmd/main.go
CONFIG_FILE_PATH=examples/config.yml

run:
	CONFIG_FILE_PATH=${CONFIG_FILE_PATH} go run ${MAIN_FILE_PATH}

CONTAINER_RUNNER=docker

generate-mocks:
	${CONTAINER_RUNNER} run --rm -v $$PWD:/src:z -w /src vektra/mockery:v2.42 --all --output internal/mocks --case snake --with-expecter
