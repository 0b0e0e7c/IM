.PHONY: all user-service friend-service message-service api

all: user-service friend-service message-service api

user-service:
	go build -o bin/user-service.exe ./service/user-service/user.go 

friend-service:
	go build -o bin/friend-service.exe ./service/friend-service/friend.go 

message-service:
	go build -o bin/message-service.exe ./service/message-service/message.go 

api:
	go build -o bin/api.exe ./gin/. 