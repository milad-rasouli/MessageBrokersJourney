build:
    go build -o ./bin/rabbit-consumer ./rabbitmq/cmd/consumer/main.go
    go build -o ./bin/rabbit-publisher ./rabbitmq/cmd/publisher/main.go

rc-run: build
    ./bin/rabbit-consumer 

rp-run: build
    ./bin/rabbit-publisher