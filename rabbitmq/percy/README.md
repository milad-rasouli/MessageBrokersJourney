
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
