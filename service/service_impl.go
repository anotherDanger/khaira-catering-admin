package service

import (
	"context"
	"database/sql"
	"khaira-admin/domain"
	"khaira-admin/helper"
	"khaira-admin/logger"
	"khaira-admin/repository"
	"khaira-admin/web"
	"mime/multipart"
	"os"
	"path/filepath"

	"time"
)

type ServiceImpl struct {
	repo repository.Repository
	db   *sql.DB
}

func NewServiceImpl(repo repository.Repository, db *sql.DB) Service {
	return &ServiceImpl{
		repo: repo,
		db:   db,
	}
}

func (svc *ServiceImpl) Login(ctx context.Context, request *domain.Admin) (*web.AdminResponse, error) {

	result, err := svc.repo.Login(ctx, svc.db, request)
	if err != nil {
		return nil, err
	}

	response := &web.AdminResponse{
		Username: result.Username,
	}

	return response, nil
}

func (svc *ServiceImpl) AddProduct(ctx context.Context, request *web.Request, file *multipart.FileHeader) (data *domain.Domain, err error) {
	tempDir := "/tmp/uploads"
	finalDir := "/app/uploads"

	filename, err := helper.SaveFile(file, tempDir)
	if err != nil {
		return nil, err
	}

	tempPath := filepath.Join(tempDir, filename)
	finalPath := filepath.Join(finalDir, filename)

	err = helper.MoveFile(tempPath, finalPath)
	if err != nil {
		return nil, err
	}

	request.ImageMetadata = filename
	now := time.Now()
	request.CreatedAt = &now

	tx, err := svc.db.Begin()
	if err != nil {
		os.Remove(finalPath)
		return nil, err
	}
	defer helper.WithTransaction(tx, &err)

	data, err = svc.repo.AddProduct(ctx, tx, (*domain.Domain)(request))
	if err != nil {
		os.Remove(finalPath)
		return nil, err
	}

	return data, nil
}

func (svc *ServiceImpl) GetProducts(ctx context.Context) (data []*domain.Domain, err error) {

	products, err := svc.repo.GetProducts(ctx, svc.db)
	if err != nil {
		logger.GetLogger("service-log").Log("get product", "error", err.Error())
		return nil, err
	}

	return products, nil
}

func (svc *ServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("delete product", "error", err.Error())
		return err
	}

	defer helper.WithTransaction(tx, &err)

	err = svc.repo.DeleteProduct(ctx, tx, id)
	if err != nil {
		logger.GetLogger("service-log").Log("delete product", "error", err.Error())
		return err
	}

	return nil
}

func (svc *ServiceImpl) UpdateProduct(ctx context.Context, request *web.Request, id string) (data *domain.Domain, err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("update product", "error", err.Error())
		return nil, err
	}

	date := time.Now()
	request.ModifiedAt = &date
	defer helper.WithTransaction(tx, &err)
	data, err = svc.repo.UpdateProduct(ctx, tx, (*domain.Domain)(request), id)
	if err != nil {
		logger.GetLogger("service-log").Log("update product", "error", err.Error())
		return nil, err
	}

	return data, nil
}

func (svc *ServiceImpl) GetOrders(ctx context.Context) (orders []*domain.Orders, err error) {

	orders, err = svc.repo.GetOrders(ctx, svc.db)
	if err != nil {
		logger.GetLogger("service-log").Log("get orders", "error", err.Error())
		return nil, err
	}

	return orders, nil
}

func (svc *ServiceImpl) UpdateOrder(ctx context.Context, entity *domain.Orders, id string) (err error) {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("update order", "error", err.Error())
		return err
	}

	defer helper.WithTransaction(tx, &err)

	err = svc.repo.UpdateOrder(ctx, tx, entity, id)
	if err != nil {
		logger.GetLogger("service-log").Log("update order", "error", err.Error())
		return err
	}

	return nil
}

func (svc *ServiceImpl) DeleteOrder(ctx context.Context, id string) error {
	tx, err := svc.db.Begin()
	if err != nil {
		logger.GetLogger("service-log").Log("delete order", "error", err.Error())
		return err
	}

	defer helper.WithTransaction(tx, &err)

	err = svc.repo.DeleteOrder(ctx, tx, id)
	if err != nil {
		logger.GetLogger("service-log").Log("delete order", "error", err.Error())
		return err
	}

	return nil
}

func (svc *ServiceImpl) GetUsers(ctx context.Context) (data []*domain.Users, err error) {
	result, err := svc.repo.GetUsers(ctx, svc.db)
	if err != nil {
		logger.GetLogger("service-log").Log("get users", "error", err.Error())
		return nil, err
	}

	return result, nil
}

func (svc *ServiceImpl) GetUserByUsername(ctx context.Context, username string) (*domain.Users, error) {
	result, err := svc.repo.GetUserByUsername(ctx, svc.db, username)
	if err != nil {
		logger.GetLogger("service-log").Log("get user by username", "error", err.Error())
		return nil, err
	}

	return result, nil

}

func (svc *ServiceImpl) GetOrdersByUsername(ctx context.Context, username string) ([]*domain.Orders, error) {
	result, err := svc.repo.GetOrderByUsername(ctx, svc.db, username)
	if err != nil {
		logger.GetLogger("service-log").Log("get orders by username", "error", err.Error())
		return nil, err
	}

	return result, nil
}

func (svc *ServiceImpl) GetOrderById(ctx context.Context, id string) (*domain.Orders, error) {
	result, err := svc.repo.GetOrderById(ctx, svc.db, id)
	if err != nil {
		logger.GetLogger("service-log").Log("get order by id", "error", err.Error())
		return nil, err
	}

	return result, nil
}
