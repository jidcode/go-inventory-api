-- +goose Up

-- Custom ENUM types
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'shipped', 'delivered', 'cancelled');

-- users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    avatar TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- inventories table
CREATE TABLE IF NOT EXISTS inventories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    code VARCHAR(50),
    quantity INTEGER NOT NULL DEFAULT 0,
    cost DECIMAL(10,2) NOT NULL DEFAULT 0,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    inventory_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_products_inventory FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE RESTRICT
);

-- storages table
CREATE TABLE IF NOT EXISTS storages (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    capacity INTEGER,
    occupied_space INTEGER DEFAULT 0,
    inventory_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_storages_inventory FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE RESTRICT
);

-- units table
CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY,
    label VARCHAR(50) NOT NULL,
    storage_id UUID NOT NULL,
    is_occupied BOOLEAN DEFAULT false,
    capacity INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_units_storage FOREIGN KEY (storage_id) REFERENCES storages(id) ON DELETE RESTRICT,
    CONSTRAINT unique_unit_label UNIQUE (storage_id, label)
);

-- categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    inventory_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_categories_inventory FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE RESTRICT
);

-- junction tables 
CREATE TABLE IF NOT EXISTS product_categories (
    product_id UUID NOT NULL,
    category_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, category_id),
    CONSTRAINT fk_product_categories_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT,
    CONSTRAINT fk_product_categories_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS product_storages (
    product_id UUID NOT NULL,
    storage_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, storage_id),
    CONSTRAINT fk_product_storages_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT,
    CONSTRAINT fk_product_storages_storage FOREIGN KEY (storage_id) REFERENCES storages(id) ON DELETE RESTRICT
);

-- product images table 
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    product_id UUID NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_images_products FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    inventory_id UUID NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    order_date TIMESTAMP WITH TIME ZONE,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_address TEXT,
    recipient_phone VARCHAR(50),
    tracking_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_inventory FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE RESTRICT
);

-- order items
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT,
    CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- optimized indexes for common queries
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_products_inventory_search ON products (inventory_id, name);
CREATE INDEX idx_storages_inventory ON storages (inventory_id);
CREATE INDEX idx_categories_inventory ON categories (inventory_id);
CREATE INDEX idx_orders_status ON orders (status, inventory_id);
CREATE INDEX idx_order_items_composite ON order_items (order_id, product_id);

-- +goose Down

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS units;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS product_storages;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS storages;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS inventories;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS order_status;
