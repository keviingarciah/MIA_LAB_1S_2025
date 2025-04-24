package main

import (
	analyzer "backend/analyzer"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{}))

	app.Post("/execute", func(c *fiber.Ctx) error {
		var req CommandRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(CommandResponse{
				Output: "Error: Petición inválida",
			})
		}

		commands := strings.Split(req.Command, "\n")
		output := ""

		for _, cmd := range commands {
			if strings.TrimSpace(cmd) == "" {
				continue
			}

			result, err := analyzer.Analyzer(cmd)
			if err != nil {
				output += fmt.Sprintf("Error: %s\n", err.Error())
			} else {
				output += fmt.Sprintf("%s\n", result)
			}
		}

		if output == "" {
			output = "No se ejecutó ningún comando"
		}

		return c.JSON(CommandResponse{
			Output: output,
		})
	})

	app.Listen(":3001")
}
