build:
	go build -o shuttle cmd/server.go

build-plugins:
	go build -buildmode=plugin -o plugins/ss.plugin plugins/ss/shadowsocks.go
	go build -buildmode=plugin -o plugins/policy-path.plugin plugins/policy-path/policy-path.go

build-image: build build-plugins
