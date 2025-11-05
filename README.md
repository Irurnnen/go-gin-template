# go-gin-template

## Генерация OpenAPI документации

Для генерации документации по комментариям в коде скачайте утилиту `swaggo/swag` с помощью команды:

```bash
go install github.com/swaggo/swag/v2/cmd/swag@latest
```

После установки сгенерируйте документацию с помощью команды 


```bash
swag init -g main.go --ot json --v3.1
```

После выполнения команды новую документацию можно увидеть в файле по пути `docs/swagger.json`