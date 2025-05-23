-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(16)  PRIMARY KEY, 
    password CHAR(64) NOT NULL,       
    coins BIGINT NOT NULL   
);

-- Таблица переводов монет (история транзакций)
CREATE TABLE IF NOT EXISTS coin_transactions (
    id SERIAL PRIMARY KEY,
    from_user VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
    to_user VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
    amount BIGINT NOT NULL                
);

-- Индексы для ускорения поиска по переводам
CREATE INDEX IF NOT EXISTS idx_coin_transactions_from ON coin_transactions(from_user);
CREATE INDEX IF NOT EXISTS idx_coin_transactions_to ON coin_transactions(to_user);

-- Таблица инвентаря пользователя
CREATE TABLE IF NOT EXISTS inventory_items (
    id SERIAL PRIMARY KEY,
    user_username VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
    item_name VARCHAR(16) NOT NULL,      
    quantity INT NOT NULL DEFAULT 0,
    UNIQUE (user_username, item_name)            -- Гарантируем уникальность записи для каждого товара у пользователя
);

CREATE INDEX IF NOT EXISTS idx_inventory_items_user ON inventory_items(user_username);

-- Таблица товаров мерча
CREATE TABLE IF NOT EXISTS merch_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE,   
    price BIGINT NOT NULL                  
);

-- Заполнение таблицы товаров мерча
INSERT INTO merch_items (name, price) VALUES 
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500)
ON CONFLICT DO NOTHING;
