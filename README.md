# fiber prometheus-middleware
[Fiber](https://gofiber.io/) middleware for Prometheus

Export metrics for request duration ```request_duration``` and request count ```request_count```

## Example
using fiber

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"log"

	fiberprom "github.com/carousell/fiber-prometheus-middleware"
)

func main() {

	r := fiber.New()
	p := fiberprom.NewPrometheus("")
	p.Use(r)

	r.Get("/health", func(ctx *fiber.Ctx) error {
		ctx.Status(200)
		log.Println(string(ctx.Request().URI().Path()))
		return ctx.JSON(map[string]string{"status": "pass"})
	})

	log.Println("main is listening on ", "8081")
	log.Fatal(r.Listen(":8081"))

}

```