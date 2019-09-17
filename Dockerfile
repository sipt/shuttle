FROM golang:latest as builder
RUN wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz -O /tmp/GeoLite2-Country.tar.gz && \
    tar zxvf /tmp/GeoLite2-Country.tar.gz -C /tmp && \
    mv /tmp/GeoLite2-Country_*/GeoLite2-Country.mmdb /

WORKDIR /shuttle/src/
COPY . /shuttle/src/
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o /shuttle/bin/shuttle cmd/server.go
RUN GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o /shuttle/bin/plugins/ss.plugin plugins/ss/shadowsocks.go
RUN GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o /shuttle/bin/plugins/policy-path.plugin plugins/policy-path/policy-path.go

FROM alpine:latest
WORKDIR /shuttle/bin/
RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
COPY --from=builder /GeoLite2-Country.mmdb /shuttle/bin/
COPY --from=builder /shuttle/bin/plugins/ss.plugin /shuttle/bin/plugins/
COPY --from=builder /shuttle/bin/plugins/policy-path.plugin /shuttle/bin/plugins/
COPY --from=builder /shuttle/bin/shuttle /usr/local/bin/

ENTRYPOINT ["shuttle"]