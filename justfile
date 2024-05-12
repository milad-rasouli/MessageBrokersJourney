build:
    go build -o ./bin/rabbit-consumer ./rabbitmq/hello-world/consumer/main.go
    go build -o ./bin/rabbit-publisher ./rabbitmq/hello-world/publisher/main.go
    go build -o ./bin/percy-publisher ./rabbitmq/percy/cmd/producer/main.go

rc-run: build
    ./bin/rabbit-consumer 

rp-run: build
    ./bin/rabbit-publisher

pp-run: build
    ./bin/percy-publisher