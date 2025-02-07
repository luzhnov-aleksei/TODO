## Инструкция
Создать в корне файл .env с полями
```
DB_HOST=postgres #Имя сервиса в docker-compose
DB_PORT=5432
DB_USER=postgres 
DB_PASSWORD=YOUR_PASS
DB_NAME=todo_db
```
Для запуска используется docker-compose, введите в терминал `docker compose up -d`  
API доступно по адресу `http://localhost:3000`


## Эндпоинты:
- `POST /tasks`
- `GET /tasks`
- `PUT /tasks/:id`
- `DELETE /tasks/:id`
