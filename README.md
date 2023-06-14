# TikTok Tech Immersion Assignment Submission
![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

## Installation

Requirement:

- golang 1.18+
- docker

To install dependency tools:

```bash
make pre
```

## Run

```bash
docker compose -f "docker-compose.yml" up -d --build
```


## Architecture
Remains as-is from the demo template. IDLs are unchanged.

## Data Storage
Handled by Redis. Sorted sets are used to store the messages.  

Key: RoomID  
RoomID is a combination of the 2 users names in the chatroom separated by a colon.  
e.g: A RoomID can be "user1:user2"  

Value: Message struct  
A Message struct has the following 3 fields:  
* Sender (string)
* Message (string)
* Timestamp (int64)  

The timestamp is used as the score in the Redis sorted set to store the messages in chronological order, as would be preferable in an instant messaging application

## Message Delivery  
Both Send and Pull APIs were implemented in rpc-server/handler.go.  
As "user1:user2" and "user2:user1" essentially refer to the same room, "user2:user1" is reconciled to "user1:user2" within the API logic when working with roomIDs.  
Received expected responses when tested with Postman while the servers ran on Docker locally.  

## Performance and Scalability
Able to handle a throughput of more than 20 QPS in JMeter testing.

## Others
Unit test in handler_test.go will not work as intended due to the new requirement of Redis.  
I am not familiar with mocking in Golang and haven't been able to implement it properly in time.
Not able to implement Kubernetes for better scaling capabilities in time either.
