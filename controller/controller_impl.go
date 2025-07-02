package controller

import (
	"context"
	"khaira-admin/domain"
	"khaira-admin/helper"
	"khaira-admin/service"
	"khaira-admin/web"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ControllerImpl struct {
	svc service.Service
}

func NewControllerImpl(svc service.Service) Controller {
	return &ControllerImpl{svc: svc}
}

func (ctrl *ControllerImpl) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var reqBody domain.Admin
	if err := c.BodyParser(&reqBody); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Invalid request body")
	}
	if err := helper.ValidateStruct(reqBody); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Invalid login data")
	}
	result, err := ctrl.svc.Login(ctx, &reqBody)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Invalid username or password")
	}
	return web.SuccessResponse[*web.AdminResponse](c, fiber.StatusOK, "Login successful", result)
}

func (ctrl *ControllerImpl) AddProduct(c *fiber.Ctx) error {
	var reqBody web.Request
	reqBody.Id = c.FormValue("id")
	reqBody.Name = c.FormValue("name")
	reqBody.Description = c.FormValue("description")

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Price must be a number")
	}
	reqBody.Price = price

	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Stock must be a number")
	}
	reqBody.Stock = stock

	if err := helper.ValidateStruct(reqBody); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Incomplete product data")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Image is required")
	}

	result, err := ctrl.svc.AddProduct(c.Context(), &reqBody, file)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to add product")
	}

	return web.SuccessResponse[*domain.Domain](c, fiber.StatusCreated, "Product added successfully", result)
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	products, err := ctrl.svc.GetProducts(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to load products")
	}
	return web.SuccessResponse[[]*domain.Domain](c, fiber.StatusOK, "Products loaded successfully", products)
}

func (ctrl *ControllerImpl) DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	err := ctrl.svc.DeleteProduct(ctx, id)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to delete product")
	}
	return web.SuccessResponse[interface{}](c, fiber.StatusNoContent, "Product deleted successfully", nil)
}

func (ctrl *ControllerImpl) UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	stockStr := c.FormValue("stock")
	priceStr := c.FormValue("price")

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Stock must be a number")
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Price must be a number")
	}

	reqBody := &web.Request{
		Name:        name,
		Description: description,
		Stock:       stock,
		Price:       price,
	}

	response, err := ctrl.svc.UpdateProduct(ctx, reqBody, id)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to update product")
	}
	return web.SuccessResponse(c, fiber.StatusOK, "Product updated successfully", response)
}

func (ctrl *ControllerImpl) GetOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	orders, err := ctrl.svc.GetOrders(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to load orders")
	}
	return web.SuccessResponse[[]*domain.Orders](c, fiber.StatusOK, "Orders loaded successfully", orders)
}

func (ctrl *ControllerImpl) UpdateOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var reqBody domain.Orders
	if err := c.BodyParser(&reqBody); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Invalid request body")
	}

	id := c.Params("id")
	if err := ctrl.svc.UpdateOrder(ctx, &reqBody, id); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to update order")
	}
	return web.SuccessResponse[interface{}](c, fiber.StatusOK, "Order updated successfully", nil)
}

func (ctrl *ControllerImpl) DeleteOrder(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	if err := ctrl.svc.DeleteOrder(ctx, id); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to delete order")
	}
	return web.SuccessResponse[interface{}](c, fiber.StatusNoContent, "Order deleted successfully", nil)
}

func (ctrl *ControllerImpl) GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := ctrl.svc.GetUsers(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Failed to load users")
	}
	return web.SuccessResponse[[]*domain.Users](c, fiber.StatusOK, "Users loaded successfully", result)
}

func (ctrl *ControllerImpl) GetUserByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	username := c.Params("username")
	result, err := ctrl.svc.GetUserByUsername(ctx, username)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "User not found")
	}
	return web.SuccessResponse[*domain.Users](c, fiber.StatusOK, "User found", result)
}

func (ctrl *ControllerImpl) GetOrdersByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	username := c.Params("username")
	result, err := ctrl.svc.GetOrdersByUsername(ctx, username)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Orders not found")
	}
	return web.SuccessResponse[[]*domain.Orders](c, fiber.StatusOK, "Orders found", result)
}

func (ctrl *ControllerImpl) GetOrderById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	result, err := ctrl.svc.GetOrderById(ctx, id)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Bad Request", "Order not found")
	}
	return web.SuccessResponse[*domain.Orders](c, fiber.StatusOK, "Order found", result)
}

func (ctrl *ControllerImpl) GetLog(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	result, err := ctrl.svc.GetLog(ctx)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadGateway, "Bad Gateway", "Client error")
	}

	return web.SuccessResponse[[]*domain.Hit](c, fiber.StatusOK, "OK", result)
}

func (ctrl *ControllerImpl) AddOrders(c *fiber.Ctx) error {
	var order domain.Orders
	err := c.BodyParser(&order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	err = ctrl.svc.AddOrders(c.Context(), &order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "ok", "Success")
}

func (ctrl *ControllerImpl) DeleteUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := ctrl.svc.DeleteUserById(c.Context(), userId); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "cannot delete user")
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "ok", "success delete user")
}
