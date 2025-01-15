CREATE TABLE menu_items (
    menu_item_id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    description VARCHAR(100),
    category VARCHAR(50),
    price DECIMAL(10,2)
)

CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id),
    old_price DECIMAL(10,2),
    new_price DECIMAL(10,2),
    change_time TIMESTAMP WITH TIME ZONE DEFAULT
)

CREATE TABLE inventory (
    ingredient_id SERIAL PRIMARY KEY,
    quantity DECIMAL(10,2),
    unit VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT
)

CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    ingredient_id INT REFERENCES inventory(ingredient_id),
    old_quantity DECIMAL(10,2),
    new_quantity DECIMAL(10,2),
    unit VARCHAR(20),
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT
)