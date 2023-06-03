# stage 1 building the code
FROM golang:1.20 as builder

COPY / /golangci
WORKDIR /golangci
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -X main.version=master -X main.commit=master -X main.date=custom" -o golangci-lint ./cmd/golangci/main.go
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -buildmode=plugin -o wrap-err-checker.so ./cmd/plugin/plugin.go

# stage 2
FROM golang:1.20
# related to https://github.com/golangci/golangci-lint/issues/3107
ENV GOROOT /usr/local/go
# don't place it into $GOPATH/bin because Drone mounts $GOPATH as volume
COPY --from=builder /golangci/golangci-lint /usr/bin/
COPY --from=builder /golangci/wrap-err-checker.so /usr/bin/
CMD ["golangci-lint"]