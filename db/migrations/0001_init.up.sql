CREATE TABLE products (
    id VARCHAR(6) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL,
    image_metadata VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE admin (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    id CHAR(36) PRIMARY KEY,
    product_id VARCHAR(6) NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone CHAR(12) NOT NULL,
    alamat VARCHAR(255) NOT NULL,
    kecamatan VARCHAR(100) NOT NULL,
    desa VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    total DOUBLE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users(username)
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_users_username ON users(username);

INSERT INTO admin (id, username, password) VALUES
('e7b8a9d4-3f5a-4c82-b7e2-2c3f49b0e9c1', 'admin', 'hashed_admin_password');