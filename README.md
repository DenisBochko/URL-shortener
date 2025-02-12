# URL-shortener
Сервис сокращения ссылок

Ручка 1, создание ссылки: POST http://localhost:8082/url с Basic Auth (myuser, mypass)
Запрос:
body:
{
    "url": "https://example.com",
    "alias": "example"
}
Ответ:
{
    "status": "OK",
    "alias": "example"
}

Ручка 2, Удаление ссылки: DELETE http://localhost:8082/url/{alias} с Basic Auth (myuser, mypass)
Запрос: http://localhost:8082/url/{alias}
Ответ:
{
    "status": "OK"
}

Ручка 3, Редирект ссылки: GET http://localhost:8082/{alias}
Запрос: http://localhost:8082/{alias}
Ответ: Редирект по alias

