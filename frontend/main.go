package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine
	engine := html.New("./templates", ".gohtml")

	// Create new Fiber instance with custom config
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base", // specify the layout template
	})

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("test.page", fiber.Map{
			"Title": "Test Page",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
