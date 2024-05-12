## Delete User
``` 
$ sudo docker exec 63191ca rabbitmqctl delete_user guest
Deleting user "guest" ...
```

## Add User
```
$ sudo docker exec 63191ca rabbitmqctl add_user ninja 1234qwer
Adding user "ninja" ...
Done. Don't forget to grant the user permissions to some virtual hosts! See 'rabbitmqctl help set_permissions' to learn more.

$ sudo docker exec 63191ca rabbitmqctl set_user_tags ninja administrator
Setting tags for user "ninja" to [administrator] 
```