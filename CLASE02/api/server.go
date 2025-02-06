package main

import "github.com/gofiber/fiber/v2" // Importa el framework web Fiber

func main() {
	// Crea una nueva instancia de la aplicación Fiber
	app := fiber.New()

	// Define una ruta GET para el endpoint raíz "/"
	app.Get("/", func(c *fiber.Ctx) error {
		// Responde con un mensaje simple
		return c.SendString("Hello, main!")
	})

	// Inicia el servidor en el puerto 3000
	app.Listen(":3000")
}
