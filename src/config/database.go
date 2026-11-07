package config

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

func ConnectDatabase() {
	once.Do(func() {
		_ = godotenv.Load()

		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = os.Getenv("POSTGRES_URL")
		}
		if dsn == "" {
			log.Println("DATABASE_URL not set, running with fallback store")
			return
		}

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Database connection error: %v\n", err)
			return
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Printf("Database ping error: %v\n", err)
			return
		}

		DB = db
		runSchema(db)
	})
}

func runSchema(db *sql.DB) {
	paths := []string{
		"src/database/schema.sql",
		"database/schema.sql",
		"../src/database/schema.sql",
		"/var/task/src/database/schema.sql",
	}
	for _, p := range paths {
		if content, err := os.ReadFile(p); err == nil {
			if _, err := db.Exec(string(content)); err != nil {
				log.Printf("Schema init warning: %v\n", err)
			}
			return
		}
	}

	if _, err := db.Exec(defaultSchemaSQL); err != nil {
		log.Printf("Fallback schema init warning: %v\n", err)
	}
}

const defaultSchemaSQL = `
CREATE TABLE IF NOT EXISTS categories (
    category_id BIGINT PRIMARY KEY,
    category_name VARCHAR(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS clients (
    client_id BIGINT PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);
ALTER TABLE clients ADD COLUMN IF NOT EXISTS username VARCHAR(50);
ALTER TABLE clients ADD COLUMN IF NOT EXISTS password VARCHAR(255);
CREATE TABLE IF NOT EXISTS services (
    service_id BIGINT PRIMARY KEY,
    category_id BIGINT REFERENCES categories(category_id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    page_range VARCHAR(120),
    base_price NUMERIC(10,2) NOT NULL
);
CREATE TABLE IF NOT EXISTS service_deliverables (
    deliverable_id BIGINT PRIMARY KEY,
    service_id BIGINT REFERENCES services(service_id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL
);
CREATE TABLE IF NOT EXISTS past_works (
    work_id BIGINT PRIMARY KEY,
    service_id BIGINT REFERENCES services(service_id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    image_url VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS service_orders (
    order_id VARCHAR(20) PRIMARY KEY,
    client_id BIGINT REFERENCES clients(client_id) ON DELETE CASCADE,
    service_id BIGINT REFERENCES services(service_id) ON DELETE RESTRICT,
    project_name VARCHAR(120) NOT NULL,
    scale VARCHAR(30) NOT NULL,
    comments TEXT,
    srs_file_url VARCHAR(255),
    status VARCHAR(30) NOT NULL,
    estimated_delivery DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO categories (category_id, category_name) VALUES
    (1, 'Portofolio Page'), (2, 'ERP'), (3, 'Mobile App'), (4, 'Games')
ON CONFLICT (category_id) DO NOTHING;
INSERT INTO clients (client_id, full_name, email, username, password) VALUES
    (1, 'M Brahmana P', 'brahmana@flexoo.id', 'brahmana', 'password123'),
    (2, 'Demo Client', 'demo@flexoo.id', 'admin', 'admin123')
ON CONFLICT (client_id) DO NOTHING;
INSERT INTO services (service_id, category_id, name, description, page_range, base_price) VALUES
    (1, 1, 'Portofolio Web', 'lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et.', '1-3 Pages', 199.00),
    (2, 3, 'Mobile App', 'Custom cross-platform mobile application built for iOS and Android with modern UI design and scalable cloud backend.', '1-4 Pages', 299.00),
    (3, 2, 'ERP', 'Comprehensive Enterprise Resource Planning platform tailored to automate business operations, finance, inventory, and reporting.', '5-20 Pages', 9999.00),
    (4, 4, 'Custom', 'Chat with customer support for custom package inquiries.', 'Chat with customer support for custom package inquiries.', 9999.00)
ON CONFLICT (service_id) DO NOTHING;
INSERT INTO service_deliverables (deliverable_id, service_id, label) VALUES
    (1, 1, 'HTML, CSS, JS Files'), (2, 1, 'Free Vercel Domain'), (3, 1, 'Static Hosting'),
    (4, 2, 'iOS & Android Builds'), (5, 2, 'API & Push Notifications'), (6, 2, 'Store Deployment Guide'),
    (7, 3, 'Full-Stack Source Code'), (8, 3, 'Neon PostgreSQL Database'), (9, 3, 'Role-Based Access Control'),
    (10, 4, 'Custom System Architecture'), (11, 4, 'Dedicated Support Channel'), (12, 4, 'Flexible Milestone Delivery')
ON CONFLICT (deliverable_id) DO NOTHING;
INSERT INTO past_works (work_id, service_id, title, image_url) VALUES
    (1, 1, 'SaaS Platform', '/static/img/saas-preview.png'),
    (2, 2, 'SaaS Platform', '/static/img/saas-preview.png'),
    (3, 3, 'SaaS Platform', '/static/img/saas-preview.png'),
    (4, 4, 'SaaS Platform', '/static/img/saas-preview.png')
ON CONFLICT (work_id) DO NOTHING;
INSERT INTO service_orders (order_id, client_id, service_id, project_name, scale, comments, srs_file_url, status, estimated_delivery, created_at) VALUES
    ('ORDER-5123', 1, 1, 'HIMASAKTA ITS WEB', 'Standard', 'Official web portal for HIMASAKTA ITS', 'SRS_HIMASAKTA_ITS.pdf', 'In Progress', '2026-11-18', CURRENT_TIMESTAMP)
ON CONFLICT (order_id) DO NOTHING;
`
