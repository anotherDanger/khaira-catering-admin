package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"khaira-admin/domain"
	"khaira-admin/logger"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

type RepositoryImpl struct {
	elastic *elasticsearch.Client
}

func NewRepositoryImpl(elastic *elasticsearch.Client) Repository {
	return &RepositoryImpl{
		elastic: elastic,
	}
}

func (repo *RepositoryImpl) Login(ctx context.Context, db *sql.DB, entity *domain.Admin) (*domain.Admin, error) {
	query := "SELECT id, username, password FROM admin WHERE username = ?"
	row := db.QueryRowContext(ctx, query, entity.Username)
	var response domain.Admin
	err := row.Scan(&response.Id, &response.Username, &response.Password)
	if err != nil {
		logger.GetLogger("repository-log").Log("login", "error", err.Error())
		return nil, err
	}
	if response.Password != entity.Password {
		return nil, errors.New("invalid credentials")
	}
	return &response, nil
}

func (repo *RepositoryImpl) AddProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain) (*domain.Domain, error) {
	query := "INSERT INTO products(id, name, description, stock, price, image_metadata, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, entity.Id, entity.Name, entity.Description, entity.Stock, entity.Price, entity.ImageMetadata, entity.CreatedAt)
	if err != nil {
		logger.GetLogger("repository-log").Log("add product", "error", err.Error())
		return nil, err
	}
	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		logger.GetLogger("repository-log").Log("add product", "error", "failed to insert row")
		return nil, errors.New("failed to insert product")
	}
	return entity, nil
}

func (repo *RepositoryImpl) GetProducts(ctx context.Context, db *sql.DB) ([]*domain.Domain, error) {
	query := "SELECT id, name, description, stock, price, image_metadata, created_at, modified_at FROM products"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger("repository-log").Log("get products", "error", err.Error())
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Domain
	for rows.Next() {
		var product domain.Domain
		var description sql.NullString
		var imageMetadata sql.NullString
		err := rows.Scan(&product.Id, &product.Name, &description, &product.Stock, &product.Price, &imageMetadata, &product.CreatedAt, &product.ModifiedAt)
		if err != nil {
			logger.GetLogger("repository-log").Log("get products", "error", err.Error())
			return nil, err
		}
		if description.Valid {
			product.Description = description.String
		}
		if imageMetadata.Valid {
			product.ImageMetadata = imageMetadata.String
		}
		products = append(products, &product)
	}
	return products, nil
}

func (repo *RepositoryImpl) DeleteProduct(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("delete product", "error", err.Error())
		return err
	}
	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		return errors.New("product not found or already deleted")
	}
	return nil
}

func (repo *RepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, entity *domain.Domain, id string) (*domain.Domain, error) {
	query := "UPDATE products SET name = ?, description = ?, stock = ?, price = ?, modified_at = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, entity.Name, entity.Description, entity.Stock, entity.Price, entity.ModifiedAt, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}
	rowAff, err := result.RowsAffected()
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}
	if rowAff == 0 {
		return nil, errors.New("no rows updated")
	}
	var product domain.Domain
	row := tx.QueryRowContext(ctx, "SELECT id, name, description, stock, price, created_at, modified_at FROM products WHERE id = ?", id)
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Stock, &product.Price, &product.CreatedAt, &product.ModifiedAt)
	if err != nil {
		logger.GetLogger("repository-log").Log("update product", "error", err.Error())
		return nil, err
	}
	return &product, nil
}

func (repo *RepositoryImpl) GetOrders(ctx context.Context, db *sql.DB) ([]*domain.Orders, error) {
	query := "SELECT id, product_id, product_name, username, name, phone, alamat, kecamatan, desa ,quantity, total, status, created_at, modified_at FROM orders"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
		return nil, err
	}
	defer rows.Close()
	var orders []*domain.Orders
	for rows.Next() {
		var order domain.Orders
		err := rows.Scan(&order.Id, &order.ProductId, &order.ProductName,
			&order.Username, &order.Name, &order.Phone,
			&order.Alamat, &order.Kecamatan, &order.Desa, &order.Quantity,
			&order.Total, &order.Status, &order.CreatedAt, &order.ModifiedAt)
		if err != nil {
			logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger("repository-log").Log("get orders", "error", err.Error())
		return nil, err
	}
	return orders, nil
}

func (repo *RepositoryImpl) UpdateOrder(ctx context.Context, tx *sql.Tx, entity *domain.Orders, id string) error {
	query := "UPDATE orders SET status = ? WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, entity.Status, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("update order", "error", err.Error())
		return err
	}
	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (repo *RepositoryImpl) DeleteOrder(ctx context.Context, tx *sql.Tx, id string) error {
	query := "DELETE FROM orders WHERE id = ?"
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		logger.GetLogger("repository-log").Log("delete order", "error", err.Error())
		return err
	}
	rowAff, err := result.RowsAffected()
	if err != nil || rowAff == 0 {
		return errors.New("order not found or already deleted")
	}
	return nil
}

func (repo *RepositoryImpl) GetUsers(ctx context.Context, db *sql.DB) ([]*domain.Users, error) {
	query := "SELECT id, username, first_name, last_name, last_accessed FROM users"
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger("repository-log").Log("get users", "error", err.Error())
		return nil, err
	}
	defer result.Close()
	var rows []*domain.Users
	for result.Next() {
		var row domain.Users
		if err := result.Scan(&row.Id, &row.Username, &row.FirstName, &row.LastName, &row.LastAccessed); err != nil {
			logger.GetLogger("repository-log").Log("get users", "error", err.Error())
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows, nil
}

func (repo *RepositoryImpl) GetUserByUsername(ctx context.Context, db *sql.DB, username string) (*domain.Users, error) {
	query := "SELECT id, username, first_name, last_name, last_accessed FROM users WHERE username = ?"
	result := db.QueryRowContext(ctx, query, username)
	var row domain.Users
	if err := result.Scan(&row.Id, &row.Username, &row.FirstName, &row.LastName, &row.LastAccessed); err != nil {
		return nil, err
	}
	return &row, nil
}

func (repo *RepositoryImpl) GetOrderByUsername(ctx context.Context, db *sql.DB, username string) ([]*domain.Orders, error) {
	query := "SELECT id, product_id, product_name, username, quantity, total, status, created_at, modified_at FROM orders WHERE username = ?"
	result, err := db.QueryContext(ctx, query, username)
	if err != nil {
		logger.GetLogger("repository-log").Log("get orders by username", "error", err.Error())
		return nil, err
	}
	defer result.Close()
	var rows []*domain.Orders
	for result.Next() {
		var row domain.Orders
		if err := result.Scan(&row.Id, &row.ProductId, &row.ProductName, &row.Username, &row.Quantity, &row.Total, &row.Status, &row.CreatedAt, &row.ModifiedAt); err != nil {
			logger.GetLogger("repository-log").Log("get orders by username", "error", err.Error())
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows, nil
}

func (repo *RepositoryImpl) GetOrderById(ctx context.Context, db *sql.DB, id string) (*domain.Orders, error) {
	query := "SELECT id, product_id, product_name, username, quantity, total, status, created_at, modified_at FROM orders WHERE id = ?"
	result := db.QueryRowContext(ctx, query, id)
	var order domain.Orders
	if err := result.Scan(&order.Id, &order.ProductId, &order.ProductName, &order.Username, &order.Quantity, &order.Total, &order.Status, &order.CreatedAt, &order.ModifiedAt); err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *RepositoryImpl) GetLog(ctx context.Context) ([]*domain.Hit, error) {

	query := `{"query": {"match_all": {}}}`

	res, err := repo.elastic.Search(
		repo.elastic.Search.WithContext(ctx),
		repo.elastic.Search.WithIndex("repository-log"),
		repo.elastic.Search.WithBody(strings.NewReader(query)),
		repo.elastic.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var searchResp domain.SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	var hits []*domain.Hit
	for i := range searchResp.Hits.Hits {
		hits = append(hits, &searchResp.Hits.Hits[i])
	}

	return hits, nil
}
