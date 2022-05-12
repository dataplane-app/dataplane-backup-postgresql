package main

import (
	"dataplane-backup/routes"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	app := fiber.New()

	//recover from panic
	app.Use(recover.New())

	// add timer field to response header
	app.Use(Timer())

	app.Use(logger.New(
		logger.Config{
			Format: "Latency: ${latency} Time:${time} Method:${method} Path:${path} Body:${body} Host:${host} UA:${ua} Header:${header} Query:${query}",
		}))

	app.Post("/postgres-backup", routes.RunPostgresBackup)

	// ------- HEALTH CHECKS ------
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("Hello 👋! Healthy 🍏")
	})

	log.Fatal(app.Listen(":8099"))
}

/* Add timer to header */
func Timer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		ms := float32(stop.Sub(start)) / float32(time.Millisecond)
		c.Append("Server-Timing", fmt.Sprintf("Dataplane;dur=%f", ms))

		return err
	}
}
