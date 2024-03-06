package server

import (
	"fmt"
	"runtime/debug"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/routes"
)

func SetupFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(
		logger.New(logger.Config{
			Format:     "${time} | ${green} ${status} ${white} | ${latency} | ${ip} | ${green} ${method} ${white} | ${path} | ${yellow} ${body} ${reset} | ${magenta} ${resBody} ${reset}\n",
			TimeFormat: "02 January 2006 15:04:05",
			TimeZone:   "Asia/Jakarta",
		}),

		func(c *fiber.Ctx) error {
			defer func() {
				if r := recover(); r != nil {
					golog.Slack.ErrorWithData("Server panic Occured", debug.Stack(), fmt.Errorf("%s", r))

					c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"message": "Server Error",
						"data":    nil,
					})
				}
			}()
			return c.Next()
		},
	)
	routes.RegisterApiRoute(app)
	wghttp.NewHttpWg()
	return app
}
