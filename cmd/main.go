package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Habit struct {
	ID int `json:"id"`
	User_ID string `json:"user_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Created_at string `json:"created_at"`
}

func main() {
	fmt.Println("Hello, Fiber-Go!")
	app := fiber.New()

	habits := []Habit{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg":"hello world"})
	})

	//Create habit
	app.Post("/api/habits", func(c *fiber.Ctx) error {
		habit := &Habit{}

		if err := c.BodyParser(habit); err != nil {
			return err
		}

		if habit.Name == "" {
			return c.Status(400).JSON(fiber.Map{"error":"Habit name is required!"})
		}

		habit.ID = len(habits) + 1
		habits = append(habits, *habit)

		return c.Status(201).JSON(habit)
	})

	log.Fatal(app.Listen(":4000"))
}