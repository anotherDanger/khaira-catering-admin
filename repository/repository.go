package repository

import (
	"context"
	"database/sql"
	"khaira-admin/domain"

	"github.com/google/uuid"
)

type Repository interface {
	Login(ctx context.Context, db *sql.DB, entity *domain.Admin) (*domain.Admin, error)
	AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error)
	GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Domain, error)
	DeleteProduct(ctx context.Context, tx *sql.Tx, id string) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error)
	GetOrders(ctx context.Context, db *sql.DB) ([]*domain.Orders, error)
	AddOrders(ctx context.Context, tx *sql.Tx, entity *domain.Orders, id uuid.UUID) error
	UpdateOrder(ctx context.Context, tx *sql.Tx, entity *domain.Orders, id string) error
	DeleteOrder(ctx context.Context, tx *sql.Tx, id string) error
	GetOrderByUsername(ctx context.Context, db *sql.DB, username string) ([]*domain.Orders, error)
	GetOrderById(ctx context.Context, db *sql.DB, id string) (*domain.Orders, error)
	GetUsers(ctx context.Context, db *sql.DB) ([]*domain.Users, error)
	GetUserByUsername(ctx context.Context, db *sql.DB, username string) (*domain.Users, error)
	DeleteUserById(ctx context.Context, db *sql.DB, id string) error
	GetLog(ctx context.Context) ([]*domain.Hit, error)
}
