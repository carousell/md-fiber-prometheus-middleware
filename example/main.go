package main

import (
	"github.com/gofiber/fiber/v2"
	"log"

	prom "github.com/carousell/fiber-prometheus-middleware"
)

func main() {

	r := fiber.New()
	p := prom.NewPrometheus("")
	p.Use(r)

	r.Get("/health", func(ctx *fiber.Ctx) error {
		ctx.Status(200)
		log.Println(string(ctx.Request().URI().Path()))
		return ctx.JSON(map[string]string{"status": "pass"})
	})

	log.Println("main is listening on ", "8081")
	log.Fatal(r.Listen(":8081"))

}
