FROM golang:latest as builder
WORKDIR /app
ADD main.go go.mod go.sum ./
ARG f=1
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dev_null .

FROM scratch
COPY --from=builder /app/dev_null /dev_null

ENTRYPOINT [ "/dev_null" ]
