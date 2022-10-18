BINARY_NAME=./bin/ovh-dynhost-ip-updater

all: build test

build:
	go build -o ${BINARY_NAME}

test:
	go test

run:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
