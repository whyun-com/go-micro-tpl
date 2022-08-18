PATH  := node_modules/.bin:$(PATH)
branch := $(GIT_BRANCH_FOR_MAKE)
now := $(shell date '+%Y%m%d%H%M%S')
EMPTY :=

protoPath := src/proto

ifeq ($(branch),$(EMPTY))
	branch := test
endif

all: tag pull config run

#ubuntu下安装protoc： sudo apt install -y protobuf-compiler
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
gen-micro-proto:
	protoc --go_out=./src   \
	--go-grpc_out=./src  \
	--go_opt=paths=import \
	--go-grpc_opt=paths=import \
	--go_opt=Msrc/proto/micro.proto=./micro \
	--go-grpc_opt=Msrc/proto/micro.proto=./micro \
	src/proto/micro.proto; \
	protoc --go_out=./src   \
	--go-grpc_out=./src  \
	--go_opt=paths=import \
	--go-grpc_opt=paths=import \
	--go_opt=Msrc/health/health.proto=./grpc_health_v1 \
	--go-grpc_opt=Msrc/health/health.proto=./grpc_health_v1 \
	src/health/health.proto

check:
	


switch2tag:
	python scripts/switch_tag.py

build:
	cd src && go build -o ../bin/micro
# 非root用户需要使用sudo运行，因为要运行docker
build-via-docker:
	docker build -f docker/complie.dockerfile --output type=local,dest=bin/micro .
# 需要安装依赖 sudo apt install graphviz -y
# go get github.com/kisielk/godepgraph
dep:
	cd src; godepgraph -s github.com/whyun-com/go-micro-tpl | dot -Tpng -o ../coverage/godepgraph.png

install:check
	go mod tidy

pull:check
	git checkout $(branch) && git pull origin $(branch)

test:check
	mkdir -p coverage && cd src/ && \
	go test ./... -v -timeout 20s -convey-story -cover -coverprofile=../coverage/coverage.out

coverage:test
	cd src && \
	go tool cover -func ../coverage/coverage.out && \
	go tool cover -html ../coverage/coverage.out -o ../coverage/index.html

run:check install
	

copy-ca:
	cp -rf scripts/cacert.pem /etc/ssl/certs
	

run-docker: check
	

clean:
	rm -rf bin/*



# 依赖环境变量 GIT_BRANCH_FOR_MAKE AREA_TYPE CI_COMMIT_TAG CI_PROJECT_NAME
build-docker:
	cd scripts && ./docker_build.sh

.PHONY: check install pull test coverage run clean gen-micro-proto build build-docker switch2tag dep build-via-docker copy-ca
