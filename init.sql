CREATE TYPE order_status AS ENUM ('active', 'inactive', 'completed', 'cancelled');
CREATE TYPE measurement_units AS ENUM ('kg', 'g', 'l');

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name VARCHAR(50) NOT NULL,
    status order_status NOT NULL,
    order_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_status_change TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10,2) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_items (
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(100),
    category VARCHAR(50),
    price DECIMAL(10,2) NOT NULL
);

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id),
    menu_item_id INT NOT NULL REFERENCES menu_items(menu_item_id),
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    customization JSONB
);

CREATE TABLE order_status_history (
    order_status_history_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id),
    old_status order_status NOT NULL,
    new_status order_status NOT NULL,
    change_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory (
    ingredient_id SERIAL PRIMARY KEY,
    quantity DECIMAL(10,2) NOT NULL,
    unit measurement_units NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES inventory(ingredient_id),
    old_quantity DECIMAL(10,2) NOT NULL,
    new_quantity DECIMAL(10,2) NOT NULL,
    unit measurement_units NOT NULL,
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT NOT NULL REFERENCES menu_items(menu_item_id),
    old_price DECIMAL(10,2) NOT NULL,
    new_price DECIMAL(10,2) NOT NULL,
    change_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO menu_items (name, description, category, price) VALUES
('Spaghetti Bolognese', 'Classic Italian pasta with meat sauce', 'Pasta', 12.50),
('Cheeseburger', 'Juicy beef burger with cheese, lettuce, and tomato', 'Burger', 9.00),
('Margherita Pizza', 'Tomato, mozzarella, and fresh basil pizza', 'Pizza', 10.00),
('Chicken Caesar Salad', 'Fresh salad with grilled chicken and Caesar dressing', 'Salad', 8.50),
('Veggie Pizza', 'Vegetarian pizza with mushrooms, peppers, and onions', 'Pizza', 11.00),
('Fish and Chips', 'Crispy battered fish with fries', 'Seafood', 13.00),
('Beef Tacos', 'Soft tacos filled with seasoned beef, lettuce, and cheese', 'Mexican', 7.50),
('Lasagna', 'Layers of pasta, beef, and cheese', 'Pasta', 14.00),
('Penne Arrabbiata', 'Pasta in a spicy tomato sauce', 'Pasta', 11.50),
('Grilled Chicken Sandwich', 'Grilled chicken breast with lettuce and tomato', 'Sandwich', 8.00);


INSERT INTO inventory (quantity, unit, created_at, updated_at) VALUES
(100, 'kg', '2025-01-01 09:00:00', '2025-01-01 09:00:00'),
(200, 'kg', '2025-01-02 10:00:00', '2025-01-02 10:00:00'),
(50, 'l', '2025-01-03 11:00:00', '2025-01-03 11:00:00'),
(30, 'g', '2025-01-04 12:00:00', '2025-01-04 12:00:00'),
(120, 'kg', '2025-01-05 14:00:00', '2025-01-05 14:00:00'),
(80, 'kg', '2025-01-06 13:00:00', '2025-01-06 13:00:00'),
(150, 'l', '2025-01-07 09:30:00', '2025-01-07 09:30:00'),
(60, 'g', '2025-01-08 08:45:00', '2025-01-08 08:45:00'),
(90, 'kg', '2025-01-09 10:15:00', '2025-01-09 10:15:00'),
(110, 'kg', '2025-01-10 11:00:00', '2025-01-10 11:00:00'),
(75, 'l', '2025-01-11 12:30:00', '2025-01-11 12:30:00'),
(50, 'kg', '2025-01-12 13:30:00', '2025-01-12 13:30:00'),
(200, 'g', '2025-01-13 14:00:00', '2025-01-13 14:00:00'),
(125, 'kg', '2025-01-14 15:00:00', '2025-01-14 15:00:00'),
(300, 'l', '2025-01-15 16:00:00', '2025-01-15 16:00:00'),
(180, 'g', '2025-01-16 17:00:00', '2025-01-16 17:00:00'),
(90, 'kg', '2025-01-17 18:00:00', '2025-01-17 18:00:00'),
(60, 'kg', '2025-01-18 19:00:00', '2025-01-18 19:00:00'),
(200, 'l', '2025-01-19 20:00:00', '2025-01-19 20:00:00'),
(150, 'g', '2025-01-20 21:00:00', '2025-01-20 21:00:00');


INSERT INTO orders (customer_name, status, order_date, last_status_change, total_amount) VALUES
('Alice', 'active', '2025-01-01 12:00:00', '2025-01-01 12:00:00', 50.00),
('Bob', 'inactive', '2025-01-02 13:30:00', '2025-01-02 13:30:00', 75.00),
('Charlie', 'active', '2025-01-03 14:00:00', '2025-01-03 14:00:00', 30.00),
('David', 'completed', '2025-01-04 15:00:00', '2025-01-04 15:00:00', 45.00),
('Eve', 'cancelled', '2025-01-05 16:00:00', '2025-01-05 16:00:00', 60.00),
('Frank', 'active', '2025-01-06 17:30:00', '2025-01-06 17:30:00', 80.00),
('Grace', 'inactive', '2025-01-07 18:00:00', '2025-01-07 18:00:00', 25.00),
('Hannah', 'completed', '2025-01-08 19:00:00', '2025-01-08 19:00:00', 65.00),
('Ivy', 'active', '2025-01-09 20:00:00', '2025-01-09 20:00:00', 55.00),
('Jack', 'cancelled', '2025-01-10 21:00:00', '2025-01-10 21:00:00', 70.00),
('Kathy', 'completed', '2025-01-11 22:00:00', '2025-01-11 22:00:00', 60.00),
('Leo', 'active', '2025-01-12 23:00:00', '2025-01-12 23:00:00', 45.00),
('Mona', 'inactive', '2025-01-13 14:30:00', '2025-01-13 14:30:00', 55.00),
('Nina', 'completed', '2025-01-14 15:30:00', '2025-01-14 15:30:00', 65.00),
('Oscar', 'cancelled', '2025-01-15 16:45:00', '2025-01-15 16:45:00', 60.00),
('Paul', 'active', '2025-01-18 18:30:00', '2025-01-18 18:30:00', 75.00),
('Quinn', 'completed', '2025-01-19 19:00:00', '2025-01-19 19:00:00', 85.00),
('Rita', 'inactive', '2025-01-20 19:30:00', '2025-01-20 19:30:00', 40.00),
('Sam', 'cancelled', '2025-01-21 20:30:00', '2025-01-21 20:30:00', 50.00),
('Tina', 'active', '2025-01-22 21:00:00', '2025-01-22 21:00:00', 60.00),
('Ursula', 'completed', '2025-01-23 22:30:00', '2025-01-23 22:30:00', 75.00);

INSERT INTO price_history (menu_item_id, old_price, new_price, change_time) VALUES
(1, 12.50, 13.00, '2025-01-01 10:00:00'),
(2, 9.00, 9.50, '2025-01-05 12:00:00'),
(3, 10.00, 10.50, '2025-01-10 13:00:00'),
(4, 8.50, 9.00, '2025-01-15 14:00:00'),
(5, 11.00, 11.50, '2025-01-20 15:00:00'),
(6, 13.00, 13.50, '2025-01-25 16:00:00'),
(7, 7.50, 8.00, '2025-02-01 09:00:00'),
(8, 14.00, 14.50, '2025-02-10 11:00:00'),
(9, 11.50, 12.00, '2025-02-15 13:00:00'),
(10, 8.00, 8.50, '2025-02-20 14:00:00');


INSERT INTO order_status_history (order_id, old_status, new_status, change_time) VALUES
(1, 'active', 'completed', '2025-01-01 12:30:00'),
(2, 'inactive', 'active', '2025-01-02 14:00:00'),
(3, 'active', 'cancelled', '2025-01-03 15:00:00'),
(4, 'completed', 'inactive', '2025-01-04 16:30:00'),
(5, 'cancelled', 'active', '2025-01-05 17:00:00'),
(6, 'active', 'completed', '2025-01-06 18:00:00'),
(7, 'inactive', 'active', '2025-01-07 19:30:00');


INSERT INTO inventory_transactions (ingredient_id, old_quantity, new_quantity, unit, modified_at) VALUES
(1, 100, 90, 'kg', '2025-01-02 09:00:00'),
(2, 200, 180, 'kg', '2025-01-04 10:30:00'),
(3, 50, 45, 'l', '2025-01-06 11:00:00'),
(4, 30, 25, 'g', '2025-01-08 12:00:00'),
(5, 120, 100, 'kg', '2025-01-10 13:30:00'),
(6, 80, 60, 'kg', '2025-01-12 14:00:00'),
(7, 150, 140, 'l', '2025-01-14 15:00:00');
