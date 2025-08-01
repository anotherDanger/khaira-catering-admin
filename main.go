package main

import (
	"khaira-admin/controller"
	"khaira-admin/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(handler controller.Controller) *fiber.App {
	app := fiber.New(fiber.Config{

		ProxyHeader: "X-Forwarded-Proto",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://catering-admin.netlify.app",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Static("/images", "/app/uploads")

	app.Post("/v1/login", handler.Login)

	protectedRoute := app.Group("/api")
	protectedRoute.Use(middleware.MyMiddleware)
	protectedRoute.Get("/v1/orders", handler.GetOrders)
	protectedRoute.Post("/v1/orders", handler.AddOrders)
	protectedRoute.Put("/v1/orders/:id", handler.UpdateOrder)
	protectedRoute.Delete("/v1/orders/:id", handler.DeleteOrder)
	protectedRoute.Get("/v1/orders/user/:username", handler.GetOrdersByUsername)
	protectedRoute.Get("/v1/orders/:id", handler.GetOrderById)

	protectedRoute.Post("/v1/products", handler.AddProduct)
	protectedRoute.Get("/v1/products", handler.GetProducts)
	protectedRoute.Delete("/v1/products/:id", handler.DeleteProduct)
	protectedRoute.Put("/v1/products/:id", handler.UpdateProduct)

	protectedRoute.Get("/v1/users", handler.GetUsers)
	protectedRoute.Get("/v1/users/:username", handler.GetUserByUsername)
	protectedRoute.Delete("/v1/users/delete/:id", handler.DeleteUserById)

	protectedRoute.Get("/v1/logs", handler.GetLog)

	return app
}

func main() {
	app, cleanup, err := InitServer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	app.Listen(":8082")

}
