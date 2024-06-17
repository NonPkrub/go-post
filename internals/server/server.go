package server

import (
	"fmt"
	"go-test/internals/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Server struct {
	postHandler ports.PostHandler
}

func NewServer(postHandler ports.PostHandler) *Server {
	return &Server{
		postHandler: postHandler,
	}
}

func (s *Server) Initialize() {
	app := fiber.New()
	v1 := app.Group("/v1")

	postGroup := v1.Group("/posts")
	{
		postGroup.Post("create", s.postHandler.Create)
		postGroup.Get("all", s.postHandler.GetAll)
		postGroup.Get(":id", s.postHandler.GetByID)
		postGroup.Patch(":id", s.postHandler.UpdateByID)
		postGroup.Delete(":id", s.postHandler.DeleteByID)
	}

	app.Listen(fmt.Sprintf(":%v", viper.GetInt("app.port")))

}
