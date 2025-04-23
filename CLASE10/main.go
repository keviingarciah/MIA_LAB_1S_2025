package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Mensaje struct {
	Title string
	Name  string
}

func sendMensaje(c *fiber.Ctx) error {
	msg := Mensaje{}
	if err := c.BodyParser(&msg); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Gracias por estar en '%s', te esperamos en la proxima %s", msg.Title, msg.Name),
	})
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("SI SALE ARCHIVOS MUCHACHOS!!!")
	})

	app.Post("/mensaje", sendMensaje)
	app.Listen(":3000")
}

// IP:300/API
