CREATE TYPE order_status as ENUM ('active', 'inactive');
CREATE TYPE measurement_units as ENUM ('kg', 'g', 'l');

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name VARCHAR(50),
    status order_status,
    order_date TIMESTAMP WITH TIME ZONE  DEFAULT CURRENT_TIMESTAMP,
    last_status_change TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10,2),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id),
    menu_item_id INT REFERENCES menu_items(menu_item_id),
    quantity INT,
    price DECIMAL(10,2),
    customization JSONB
);


CREATE TABLE order_status_history (
    order_status_history_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id),
    old_status order_status,
    new_status order_status,
    change_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_item_ingredients (
    menu_item_id INT REFERENCES menu_items(menu_item_id),
    ingredient_id INT REFERENCES inventory(ingredient_id),
    quantity DECIMAL(10,2)
);

CREATE TABLE menu_items (
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    description VARCHAR(100),
    category VARCHAR(50),
    price DECIMAL(10,2)
);

CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id),
    old_price DECIMAL(10,2),
    new_price DECIMAL(10,2),
    change_time TIMESTAMP WITH TIME ZONE DEFAULT
);

CREATE TABLE inventory (
    ingredient_id SERIAL PRIMARY KEY,
    quantity DECIMAL(10,2),
    unit VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT
);

CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    ingredient_id INT REFERENCES inventory(ingredient_id),
    old_quantity DECIMAL(10,2),
    new_quantity DECIMAL(10,2),
    unit VARCHAR(20),
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT
);
