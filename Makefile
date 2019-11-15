SRCS := $(shell find ./src -type f -name '*.go')
OS ?=

.PHONY: all
all: linux mac win

.PHONY: linux mac win
linux: output/app_linux
mac: output/app_mac
win: output/app_win.exe

output/app_linux: $(SRCS)
	@OS="linux" $(MAKE) make_app
	@mv src/out output/app_linux

output/app_mac: $(SRCS)
	@OS="darwin" $(MAKE) make_app
	@mv src/out output/app_mac

output/app_win.exe: $(SRCS)
	@OS="windows" $(MAKE) make_app
	@mv src/out output/app_win.exe

# OS指定でビルド処理共通化
.PHONY: make_app
make_app:
	docker run -i --rm \
	-v `pwd`/src:/work \
	-v `pwd`/src/mod:/go/pkg/mod \
	-w /work \
	-e GOOS=${OS} \
	-e GOARCH=amd64 \
	-e CGO_ENABLED=0 \
	golang:1.13 \
	go build -o out

.PHONY: clean
clean:
	rm -f output/app_*
	docker run -it --rm \
	-v `pwd`/src/mod:/go/pkg/mod \
	golang:1.13 go clean -modcache || :

.PHONY: test
test: $(SRCS)
	docker run -i --rm \
	-v `pwd`/src:/work \
	-v `pwd`/src/mod:/go/pkg/mod \
	-w /work \
	-e GOOS=linux \
	-e GOARCH=amd64 \
	-e CGO_ENABLED=0 \
	golang:1.13 \
	go test -cover -v . 