# Http downloader:

## ТЗ
- может выполнять одновременно до N запросов
- бесконечно можно набивать очередь новыми задачами (без блокировки потока управления). То есть могу добавить сразу мульт запросов и продолжить работу
- умеет graceful shutdown. При этом, после получения сигнала остановки - задачи в обработку больше не берутся.
- после получения сигнала остановки надо дождаться обработки всей накопившейся очереди
- в качестве задачи на отработку принимает строку с адресом
- результат выполнения запросов можно как то получать. На свой вкус. Может быть это канал с ответами, может ещё что то
- со звездочкой (можно устно проговорить, не писать код): экономия памяти для ситуации, где нам накидывают очень много задач.


## Описание

При запуске приложения настройки подхватываются из переменных окружения. Считывания постановку задач можно осуществить двумя способами:
- при запуске указать из какого файла их необходимо считать
- в виде http запроса на ручку `http://127.0.0.1:8080/api/loader`


## Настройки

Считывания настроек происходит из переменных окружения:
- `FILE` - имя файла из которого будут считаны url и добавлены в очередь заданий
- `REQUEST_LIMIT` - кол-во одновременных выполняемых запросов, если не установлено то берется значение по умолчанию (5 запросов)

## Пример использования

Запускаем приложение:
```
go run ./cmd/app/main.go
```

запускается сервер, на который мы можем отправлять очередь задач на загрузку

```
curl -X POST \
     -H "Content-Type: application/json" \
     -d '["https://yandex.ru", "https://habr.ru", "https://google.com", "http://fe.tt", "we", "https://ffff"]'  \
     http://127.0.0.1:8080/api/loader 
```
