

# ----------------------------
# Build and Run Commands
# ----------------------------

run:
	go run cmd/app/main.go -cfg configs/dev

build-models:
	docker run --rm --user $(shell id -u):$(shell id -g) -e GOPATH=$(go env GOPATH):/go -v ${HOME}:${HOME} -w $(shell pwd) quay.io/goswagger/swagger generate model --spec=../swagger.yaml
	go mod tidy


