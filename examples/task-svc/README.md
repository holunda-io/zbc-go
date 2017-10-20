# task-svc

Simple microservice example on how to use task subscriptions.

# Note

* Credit system is not implemented so worker will eventually stop accepting work
* Broker is expected to run the client API on 51015


# Running

1.Run the broker

```
docker-compose up -d
```

2.Run the example with 

```
go run main.go
```

3.Create default-topic with the following command
```
zbctl create topic
```

# Using

There are 3 endpoints implemented.

```
localhost:3000/start
```

This will start worker and give us response of the worker ID.

```
localhost:3000/stats
```
Will give us basic runtime stats.

```
localhost:3000/stop
```
Will stop all running workers.


