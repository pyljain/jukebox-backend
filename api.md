curl -v http://localhost:8080/0/api/v1/conversations

curl -v -X POST http://localhost:8080/api/v1/conversations -d @samples/api/post_conversation.json

curl -v http://localhost:8080/0/api/v1/conversations/3/messages