package routes

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func RouteQueue(c *fiber.Ctx) error {
	var text string
	if message != nil {
		text = message.Text
	}
	return c.Render("queue", fiber.Map{
		"Text": text,
	})
}

func RouteIndex(c *fiber.Ctx) error {
	var text string
	if message != nil {
		text = message.Text
	}
	return c.Render("index", fiber.Map{
		"Text": text,
	})
}

func RouteMessage(c *fiber.Ctx) error {
	clientId := c.Query("clientId", "")
	if matched, err := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", clientId); !matched || err != nil {
		return c.Redirect("queue")
	}

	if currentClient, _ := clientQueue.First(); clientId == "" || clientId != currentClient {
		return c.Redirect("queue")
	}

	var text string
	if message != nil {
		text = message.Text
	}

	return c.Render("message", fiber.Map{
		"Text":     text,
		"ClientId": clientId,
	})
}

func RouteGetMessage(c *fiber.Ctx) error {
	out := make(map[string](map[string]string))
	out["message"] = nil
	if message != nil {
		out["message"] = make(map[string]string)
		out["message"]["text"] = message.Text
	}
	return c.JSON(out)
}

func RoutePostMessage(c *fiber.Ctx) error {
	payload := struct {
		ClientId string `json:"clientId"`
		Text     string `json:"text"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Could not parse the payload",
		})
	}
	if payload.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Your message cannot be blank",
		})
	}
	if currentClient, _ := clientQueue.First(); message == nil || currentClient == "" || currentClient == payload.ClientId {
		defer func() {
			postMessageChannel <- true
		}()
		message = &Message{payload.Text}
		return c.JSON(&fiber.Map{
			"text": message.Text,
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
		"error": "You aren't the first of the queue",
	})
}
