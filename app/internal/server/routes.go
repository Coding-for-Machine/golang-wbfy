package server

import (
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.Home)

	s.App.Get("/health", s.healthHandler)
	// WebSocket endpoint
	s.App.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// 1. Docker konteynerni terminal bilan ishga tushiramiz
		cmd := exec.Command("docker", "run", "-it", "--rm", "fiber-terminal")

		ptmx, err := pty.Start(cmd)
		if err != nil {
			log.Println("PTY error:", err)
			return
		}
		defer func() {
			ptmx.Close()
			cmd.Process.Kill()
		}()

		// PTY → WebSocket
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := ptmx.Read(buf)
				if err != nil {
					break
				}
				if err := c.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					break
				}
			}
		}()

		// WebSocket → PTY
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			ptmx.Write(msg)
		}
	}))

}

func (s *FiberServer) Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
