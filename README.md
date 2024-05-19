# Message Broker Journey
I've decided to explore more in the software development world to do so I made this repo to work around **message brokers** such as **NATS**, **RabbitMQ**, **Kafka**, etc. I will put all the snippet codes I've done here also, I'll try to show you all different aspects of each message broker. 
Message brokers are widely used in *EDA(Event-Driven Architecture)*.

### To get the project
```bash
git clone https://github.com/Milad75Rasouli/MessageBrokersJourney

``` 
## RabbitMQ
Here is the way that you can run the examples:
1. Go to the project directory.
```bash
cd MessageBrokersJourney
```

2. Perform this command to run RabbitMQ in your local machine:
```bash
sudo docker-compose up rabbitmq
```

> [!NOTE] 
> If you could perform the command successfully, you might want to take a look at the *RabbitMQ Management* at *http://localhost:15672* and username and password are up to you. You can find them in *docker-compose.yml* file. 

3. Perform this command to run a **consumer**(you can perform this command as many as you like in different terminals/CMD to get more consumer)
```bash
just pc-run
```
> [!IMPORTANT]
> If you have a problem with **just**https://github.com/casey/just you might need to [install it on you machine first](https://github.com/casey/just) and put the installation path in your Linux/Windows **PATH**

4. Perform this command to run a **producer**

```bash
just pp-run
```
5. Nothing! but if you'd like to see more examples you can checkout to other commits. they are listed here:
    - [RPC example with TLS and rabbitmq_definitions.json](https://github.com/Milad75Rasouli/MessageBrokersJourney/releases/tag/rabbitmq-rpc-tls-conf)
    - [RPC example](https://github.com/Milad75Rasouli/MessageBrokersJourney/releases/tag/rabbitmq-rpc)
    - [Fanout example](https://github.com/Milad75Rasouli/MessageBrokersJourney/releases/tag/rabbitmq-fanout)
    - [Work Queues example](https://github.com/Milad75Rasouhttps://github.com/Milad75Rasouli/MessageBrokersJourney/releases/tag/rabbitmq-fanoutli/MessageBrokersJourney/releases/tag/rabbitmq-worker)

6. Perform this command whenever you are done and want to shut the RabbitMQ service:
```bash
sudo docker-compose down rabbitmq -v
```
> [!TIP]
> The command must be performed in the main project directory.

[Take a look at here for more examples](https://www.rabbitmq.com/tutorials)

## NATS
Here is how you can run the examples:

1. Go to the project directory.
```bash
cd MessageBrokersJourney
```

2. Perform this command to run NATS in your local machine:
```bash
sudo docker-compose up nats
```

> [!NOTE] 
> If you could perform the command successfully, you might want to take a look at the *NATS monitoring* at *http://localhost:8222* you don't need any passwords for that. 



[Take a look at here for more examples](https://natsbyexample.com/)
