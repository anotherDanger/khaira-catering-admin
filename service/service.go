package service

import (
	"context"
	"khaira-admin/domain"
	"khaira-admin/web"
	"mime/multipart"
)

type Service interface {
	Login(ctx context.Context, request *domain.Admin) (*web.AdminResponse, error)
	AddProduct(ctx context.Context, request *web.Request, file *multipart.FileHeader) (*domain.Domain, error)
	GetProducts(ctx context.Context) ([]*domain.Domain, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, request *web.Request, id string) (*domain.Domain, error)
	GetOrders(ctx context.Context) ([]*domain.Orders, error)
	UpdateOrder(ctx context.Context, entity *domain.Orders, id string) error
	DeleteOrder(ctx context.Context, id string) error
	GetUsers(ctx context.Context) ([]*domain.Users, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.Users, error)
	GetOrdersByUsername(ctx context.Context, username string) ([]*domain.Orders, error)
	GetOrderById(ctx context.Context, id string) (*domain.Orders, error)
}
