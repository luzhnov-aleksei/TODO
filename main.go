package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

func connectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@"
	connStr += os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	conn, err := pgx.Connect(context.Background(), "postgres://"+connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	db = conn
	log.Println("Подключение к БД успешно")
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	connectDB()
	app := fiber.New()

	app.Post("/tasks", createTask)
	app.Get("/tasks", getTasks)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)

	log.Fatal(app.Listen(":3000"))
}

func createTask(c *fiber.Ctx) error {
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	err := db.QueryRow(context.Background(), query, task.Title, task.Description, task.Status).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка добавления задачи"})
	}

	return c.JSON(task)
}

func getTasks(c *fiber.Ctx) error {
	rows, err := db.Query(context.Background(), "SELECT id, title, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения списка задач"})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println(err)
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка чтения данных"})
			
		}
		tasks = append(tasks, t)
	}

	return c.JSON(tasks)
}

func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	query := "UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=now() WHERE id=$4"
	status, err := db.Exec(context.Background(), query, task.Title, task.Description, task.Status, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления задачи"})
	}
	if status.RowsAffected() == 0{
		return c.Status(404).JSON(fiber.Map{"error": "Такой задачи нет"})

	}

	return c.JSON(fiber.Map{"message": "Задача обновлена"})
}

func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	status, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1", id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Ошибка удаления задачи"})
	}
	if status.RowsAffected() == 0{
		return c.Status(404).JSON(fiber.Map{"error": "Такой задачи нет"})

	}
	return c.JSON(fiber.Map{"message": "Задача удалена"})
}
