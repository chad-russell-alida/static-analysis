# Based on https://github.com/golangci/golangci-lint/blob/master/build/Dockerfile

# stage 1 building the code
FROM golang:1.21 as builder

COPY / /golangci
WORKDIR /golangci

# Build golangci-lint binary with CGO_ENABLED=1, so it's dynamically linked.
#
# NOTE: We can't use golangci/golangci-lint:latest image as a base, since it
# has statically linked binary, which can't load *.so plugins.
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -X main.version=master -X main.commit=master -X main.date=custom" -o golangci-lint ./cmd/golangci/main.go
RUN CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -buildmode=plugin -o wrap-err-checker.so ./cmd/plugin/plugin.go

# stage 2
FROM golang:1.21
# related to https://github.com/golangci/golangci-lint/issues/3107
ENV GOROOT /usr/local/go
# Set all directories as safe
RUN git config --global --add safe.directory '*'

COPY --from=builder /golangci/golangci-lint /usr/bin/
COPY --from=builder /golangci/wrap-err-checker.so /usr/bin/
CMD ["golangci-lint"]
