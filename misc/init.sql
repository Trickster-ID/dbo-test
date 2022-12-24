CREATE TABLE user (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "order" (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES user(user_id),
    product_id INTEGER REFERENCES product(product_id),
    quantity INTEGER NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    shipping_address VARCHAR(255),
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO user (username, password, email, first_name, last_name, is_admin) VALUES
('user1', 'password1', 'user1@example.com', 'User', 'One', FALSE),
('user2', 'password2', 'user2@example.com', 'User', 'Two', FALSE),
('user3', 'password3', 'user3@example.com', 'User', 'Three', FALSE),
('user4', 'password4', 'user4@example.com', 'User', 'Four', FALSE),
('user5', 'password5', 'user5@example.com', 'User', 'Five', FALSE),
('admin1', 'password1', 'admin1@example.com', 'Admin', 'One', TRUE),
('admin2', 'password2', 'admin2@example.com', 'Admin', 'Two', TRUE),
('admin3', 'password3', 'admin3@example.com', 'Admin', 'Three', TRUE),
('admin4', 'password4', 'admin4@example.com', 'Admin', 'Four', TRUE),
('admin5', 'password5', 'admin5@example.com', 'Admin', 'Five', TRUE);

INSERT INTO product (name, price, description, image_url, quantity) VALUES
('Keyboard', 59.99, 'Mechanical keyboard with backlighting', 'https://example.com/keyboard.jpg', 10),
('Monitor', 199.99, '27-inch 1080p monitor with thin bezels', 'https://example.com/monitor.jpg', 5),
('Mouse', 39.99, 'Wireless gaming mouse with high DPI sensor', 'https://example.com/mouse.jpg', 15),
('Headphones', 79.99, 'Over-ear headphones with noise cancelling', 'https://example.com/headphones.jpg', 8),
('Speakers', 99.99, 'Bluetooth speaker with 360-degree sound', 'https://example.com/speakers.jpg', 3),
('Laptop', 999.99, '15-inch gaming laptop with NVIDIA graphics', 'https://example.com/laptop.jpg', 2),
('Smartphone', 499.99, 'Latest model smartphone with large AMOLED display', 'https://example.com/smartphone.jpg', 10),
('Tablet', 299.99, '10-inch tablet with long battery life', 'https://example.com/tablet.jpg', 7),
('Printer', 149.99, 'Inkjet printer with wireless connectivity', 'https://example.com/printer.jpg', 4),
('External hard drive', 99.99, '1TB external hard drive with fast transfer speeds', 'https://example.com/hard_drive.jpg', 6);

INSERT INTO "order" (user_id, product_id, quantity, price, shipping_address, status) VALUES
(1, 1, 1, 59.99, 'Jl. Sudirman No.1, Jakarta', 'pending'),
(2, 2, 1, 199.99, 'Jl. Gatot Subroto No.45, Jakarta', 'pending'),
(3, 3, 1, 39.99, 'Jl. HR Rasuna Said No.12, Jakarta', 'pending'),
(4, 4, 1, 79.99, 'Jl. Diponegoro No.22, Semarang', 'pending'),
(5, 5, 1, 99.99, 'Jl. Sudirman No.55, Surabaya', 'pending'),
(6, 6, 1, 999.99, 'Jl. Jenderal Sudirman No.66, Bandung', 'pending'),
(7, 7, 1, 499.99, 'Jl. Pahlawan No.77, Medan', 'pending'),
(8, 8, 1, 299.99, 'Jl. Mayor Jenderal Sutoyo No.88, Yogyakarta', 'pending'),
(9, 9, 1, 149.99, 'Jl. Kapten Tendean No.99, Denpasar', 'pending');