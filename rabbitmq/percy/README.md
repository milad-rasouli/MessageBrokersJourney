
## Intro
Virtual Host (vhost) is like a namespace that contains channels and exchanges. 
it's commonly used for limiting resources. and it happens logically

## Delete User
```bash
$ sudo docker exec 63191ca rabbitmqctl delete_user guest
Deleting user "guest" ...
```

## Add User
```bash
$ sudo docker exec 63191ca rabbitmqctl add_user ninja 1234qwer
Adding user "ninja" ...
Done. Don't forget to grant the user permissions to some virtual hosts! See 'rabbitmqctl help set_permissions' to learn more.

$ sudo docker exec 63191ca rabbitmqctl set_user_tags ninja administrator
Setting tags for user "ninja" to [administrator] 
```

## Add Virtual Host
```bash
$ sudo docker exec 63191ca rabbitmqctl add_vhost customer
Adding vhost "customer" ...

```

## Giving Virtual Host Permissions  
There are 3 type of it:
1- configuration  "^customer.*"  for anything ".*"
2- read     "queue.*"
3- write    "^customers.*"

```bash
$ sudo docker exec 63191ca rabbitmqctl set_permissions -p  customer ninja  ".*" ".*" ".*"
Setting permissions for user "ninja" in vhost "customer" ...
```

## How RabbitMQ Works
Producer#N -> Exchange(Decides where message goes) -> Queue(A queue is message buffer default it's FIFO) -> Consumer
Exchange: it's like a broker
Bind: it's basically a rule of set of rules.

## Exchanges
### Direct Exchange
Exact Routing key match!
P(with routing key AKA Topic e.g. customer_created) -> E --> Queue(customer_created) -> C 

### FanOut Exchange
Ignore Routing key!
P(customer_created) -> E --> Queue(customer_created) -> C 
                        \-> Queue(customer_emailed)

### Topic Exchange
Rules on routing key delimited by "."

"#" means match zero or more (like customer.created.#)
"*" means match everything (like customer.*.february)

P(customer.created.february) -> E(customer.created.#) --> Queue(customer_created) -> C 
                            \-> E(customer.created.match) --> Queue(customer_emailed)

### Topic Exchange
Rules based on the extra header it's more like key-value

P(browser=Linux) -> E(browser=Linux) --> Queue(customer_linux) -> C 
                    \-> E(browser=Windows) --> Queue(customer_windows)

## Create Queue And Exchange
``` bash
$ sudo docker exec 63191cae1166 rabbitmqadmin declare exchange --vhost=customer name=customer_test2 type=topic -u ninja -p 1234qwer durable=true
exchange declared
```

Since, the user does not have the permission to send data to the  topic we have to give it the permission.

```bash
$ sudo docker exec 63191 rabbitmqctl set_topic_permissions -p customer ninja customer_test2  "^customer.*" "^customer.*" 
Setting topic permissions on "customer_test2" for user "ninja" in vhost "customer" ...

```
After, we need a function that binds the exchange to the queue.

## Use fanout exchange
first delete the previous exchange 
```bash
$ sudo docker exec 63191 rabbitmqadmin delete exchange name=customer_test2 --vhost=customer -u ninja -p 1234qwer
exchange deleted
```
then redeclare the queue:
ps: there is no way to modify the type of the exchange except deleting it and recreate it.
```bash
sudo docker exec 63191 rabbitmqadmin declare exchange name=customer_test2 --vhost=customer type=fanout durable=true -u ninja -p 1234qwer
exchange declared
``` 
then, we need to update the permission:
```bash
$ sudo docker exec 63191 rabbitmqctl set_topic_permissions -p customer ninja customer_test2 ".*" ".*" 
Setting topic permissions on "customer_test2" for user "ninja" in vhost "customer" ...
```
now, we are ready to go!

** caution:**
there is a drawback in this way. if no consumer creates queue the producer will sends to nowhere. meaning the data will be lost.

## Use RPC Exchange
we need 2 exchanges for this example:
1. for replying
2- fo the data

### For data
make sure you have this exchange:
```bash
sudo docker exec 63191 rabbitmqadmin declare exchange name=customer_test2 --vhost=customer type=fanout durable=true -u ninja -p 1234qwer
exchange declared
``` 
you might need to delete it and then create the exchange here you go!
```bash
$ sudo docker exec 63191 rabbitmqadmin delete exchange name=customer_test2 --vhost=customer -u ninja -p 1234qwer                          
exchange deleted 
```
### For replying
we need to have a direct exchange fo this part.
create the exchange:
```bash 
$ sudo docker exec 63191 rabbitmqadmin declare exchange name=customer_callback type=direct --vhost=customer durable=true -u ninja -p 1234qwer 
exchange declared
```

permission:
```bash
sudo docker exec 63191 rabbitmqctl set_topic_permissions -p customer ninja customer_callback ".*" ".*"
Setting topic permissions on "customer_callback" for user "ninja" in vhost "customer" ...

```

a good rule of thumb:
never use a connection for publishing and consuming.

## Encrypt the data
to to so, we are going to use TLS.
just follow the instruction:
```bash
$ git clone https://github.com/rabbitmq/tls-gen tls-gen

$ cd tls-gen/basic

$ make PASSWORD=

$ make verify

$ sudo chmod 644 result/*
```
then, we need to delete the rabbitmq instance and rerun it.
go to the root project path.
```bash
$ sudo docker compose down -v

$ sudo rm -rf rabbitmq-data

$ sudo docker compose up
```

then edit the rabbitmq.conf
just make sure the name of the files are correct as follow:
```bash
$ sudo docker exec -it a4f502f44e56  sh

$ ls certs
ca_certificate.pem            client_milad_certificate.pem  server_milad_certificate.pem
ca_key.pem                    client_milad_key.pem          server_milad_key.pem
client_milad.p12              server_milad.p12

```

## Generate password for TLS

```bash
$ cd script

$ sudo chmod +x encodepassword.sh

$ ./encodepassword.sh 
XpHhjk92Ez0M7fN5y0DiYmIm+LbdA+hrgCHnnHx6arAsht9H
```