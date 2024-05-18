build:
    go build -o ./bin/rabbit-consumer ./rabbitmq/hello-world/consumer/main.go
    go build -o ./bin/rabbit-publisher ./rabbitmq/hello-world/publisher/main.go
    go build -o ./bin/percy-publisher ./rabbitmq/percy/cmd/producer/main.go
    go build -o ./bin/percy-consumer ./rabbitmq/percy/cmd/consumer/main.go
    
build-nats:
    go build -o ./bin/producer-nats ./nats/producer/main.go
    go build -o ./bin/consumer-nats ./nats/consumer/main.go

np-run: build-nats
    ./bin/producer-nats
nc-run: build-nats
    ./bin/consumer-nats

rc-run: build
    ./bin/rabbit-consumer 

rp-run: build
    ./bin/rabbit-publisher

pp-run: build
    ./bin/percy-publisher

pc-run: build
    ./bin/percy-consumer