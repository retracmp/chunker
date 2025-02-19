package hoster

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func Start() error {
	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	r.Static("/", "./builds")

	r.Hooks().OnListen(func(ld fiber.ListenData) error {
		fmt.Printf("development file hoster service started on %s\n", ld.Port)
		return nil
	})

	if err := r.Listen(":80"); err != nil {
		return err
	}
	
	return nil
}