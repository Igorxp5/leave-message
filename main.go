package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
	"github.com/gofiber/websocket/v2"
	"github.com/igorxp5/leave-message/routes"
)

func main() {
	engine := pug.New("./views", ".pug")
	engine.Reload(true) //Dynamic loading the templates

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public") //Serve assets files

	app.Get("/", routes.RouteIndex)
	app.Get("/queue", routes.RouteQueue)
	app.Get("/message", routes.RouteMessage)

	app.Get("/api/message", routes.RouteGetMessage)
	app.Post("/api/message", routes.RoutePostMessage)

	app.Get("/ws", websocket.New(routes.WebSocket))

	host := "127.0.0.1"
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	go routes.StartQueueManager()
	log.Fatal(app.Listen(host + ":" + port))
}
