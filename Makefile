build:
	go build -o _output/shuttle cmd/server.go

build-plugins:
	go build -buildmode=plugin -ldflags '-w -s' -o _output/plugins/ss.plugin plugins/ss/shadowsocks.go
	go build -buildmode=plugin -ldflags '-w -s' -o _output/plugins/policy-path.plugin plugins/policy-path/policy-path.go

upgrade-geo:
	mkdir -p tmp
	mkdir -p _output
	wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz -O tmp/GeoLite2-Country.tar.gz && \
	tar zxvf tmp/GeoLite2-Country.tar.gz -C tmp/ && \
	mv tmp/GeoLite2-Country_*/GeoLite2-Country.mmdb _output/ &&\
	rm -rf tmp

build-all: build build-plugins upgrade-geo

build-image: build-all
	docker build -t ${TARGET} -f Dockerfile .