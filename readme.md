# message-responder

## Что делает

1. Подписывается на Kafka-топик `KAFKA_TOPIC_NAME_NORMALIZED_MSG`, получает `NormalizedRequest` (см. `internal/contract`), логирует метаданные и пробует обработать сообщение последовательно всеми хендлерами.
2. Если сообщение содержит изображение, `HasFileHandler` пересылает десериализованный `OcrRequest` в `KAFKA_TOPIC_NAME_OCR_REQUEST` и возвращает пользователю подтверждение.
3. Если ни один обработчик не подошёл, `DefaultHandler` отвечает шаблонным текстом про ожидание изображения.
4. Ответы сериализуются как `NormalizedResponse`, добавляют префикс исходного `source` к `KAFKA_TOPIC_NAME_TG_RESPONSE_PREPARER` и отправляются туда через Kafka-продюсер.

## Запуск

1. Задайте окружение (см. список ниже).
2. Соберите и запустите локально:
   ```bash
   go run ./cmd/message-responder
   ```
3. Или соберите Docker-образ и запустите его с переменными:
   ```bash
   docker build -t message-responder .
   docker run --rm -e ... message-responder
   ```

## Переменные окружения

Все переменные обязательны, кроме `SASL_USERNAME`/`SASL_PASSWORD`, если кластер Kafka работает без SASL.

- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — список брокеров (`host:port[,host:port]`).
- `KAFKA_GROUP_ID_MESSAGE_RESPONDER` — `consumer group` для подписки на запросы.
- `KAFKA_TOPIC_NAME_NORMALIZED_MSG` — входной топик с нормализованными событиями Telegram.
- `KAFKA_TOPIC_NAME_TG_RESPONSE_PREPARER` — суффикс ответного топика (`source + response`, где `source` берётся из запроса).
- `KAFKA_CLIENT_ID_MESSAGE_RESPONDER` — идентификатор Kafka-клиента (продюсер и консьюмер).
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — если нуждается Kafka в SASL/PLAIN.
- `KAFKA_TOPIC_NAME_OCR_REQUEST` — топик с запросами на OCR (передаётся внешнему OCR-процессу).

## Примечания

- Переход к OCR идёт только при наличии хотя бы одного медиа-объекта, у которого MIME-тип начинается с `image/` или тип явно указывает на изображение (`image`, `photo`, `picture`, `img`).
- Ответы публикуются в топик с суффиксом `KAFKA_TOPIC_NAME_TG_RESPONSE_PREPARER`, причём `source` из запроса добавляется в начало названия топика.
- Логика отправки на OCR и обработки ответа сосредоточена в `internal/responder`, Kafka-коммуникация в `internal/messaging`, конфигурация в `internal/config`.
