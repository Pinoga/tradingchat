# TradingChat

## Requirements

- Docker v20+
- Go v1.16+

## Starting

```bash
docker run -d -p 27017:27017 mongo
docker run -d -p 15672:15672 -p 5672:5672 rabbitmq:3-management
cd $root/cmd/bot && go run .
cd $root/cmd/web && export $(cat .env) && go run main.go
```

## Hope you enjoy!
