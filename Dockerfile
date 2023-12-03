FROM golang:latest as builder
WORKDIR /app
ADD main.go go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dev_null .

FROM scratch
ENV GIN_MODE=release
COPY --from=builder /app/dev_null /dev_null
ENTRYPOINT [ "/dev_null" ]
