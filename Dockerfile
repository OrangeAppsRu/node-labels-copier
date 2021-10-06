FROM docker.io/golang:1.17.1 as builder
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download
COPY ./*.go ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
RUN chmod +x ./main

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/OrangeAppsRu/node-labels-copier
COPY --from=builder ./app/main /main
ENTRYPOINT ["/main"]