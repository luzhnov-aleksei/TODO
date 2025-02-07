## Инструкция
Создать БД и применить скрипт init_database.sql
Создать файл .env с полями 
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres 
DB_PASSWORD=YOUR_PASS
DB_NAME=todo_db
```
Для запуска введите в терминал `go run main.go`  
API доступно по адресу `http://localhost:3000`


## Эндпоинты:
- `POST /tasks`
- `GET /tasks`
- `PUT /tasks/:id`
- `DELETE /tasks/:id`
