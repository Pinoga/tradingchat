# TradingChat

## Requirements

- Docker v20+
- Go v1.16+

## Starting

```bash
docker run -d -p 27017:27017 mongo
docker run -d -p 15672:15672 -p 5672:5672 rabbitmq:3-management
cd $root/cmd/web && export $(cat .env) && go run main.go
cd $root/cmd/bot && go run .

```

## Considerations

- Didn't have the time to finish some features/code, namely:
  - Unit tests
  - Show only first 50 messages in the chat
- Bugs:
  - Chatrooms behavior sometimes don't work as expected, the connection closes everytime a client enters a room
  - Bot message error handling stopped working in the final ajustments
  - And probably many others lol...
- Other:
  - UI code is bad. I didn't put any effort into making it scalable
  - The 'internal' directory was left to refactoring
  - Lack of comments
  - The bot and RabbitMQ parts were particularly rushed, MUCH room for improvement on scalability/modularity

## Hope you enjoy!
