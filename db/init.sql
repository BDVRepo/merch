CREATE DOMAIN ID AS INTEGER;
CREATE DOMAIN CODE AS VARCHAR(50);

CREATE TABLE doc_merchs (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code CODE UNIQUE NOT NULL,
    name TEXT NOT NULL,
    price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE auth_users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE info_users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance INT NOT NULL CHECK (balance >= 0)
);

CREATE TABLE doc_merch_transactions (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    merch_id INTEGER NOT NULL REFERENCES doc_merchs(id) ON DELETE CASCADE,
    receiver_id INTEGER NOT NULL REFERENCES info_users(id) ON DELETE CASCADE,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT now()
);

-- Наполнение таблицы doc_merchs данными из задания
INSERT INTO doc_merchs (code, name, price) VALUES
    ('T_SHIRT', 't-shirt', 80),
    ('CUP', 'cup', 20),
    ('BOOK', 'book', 50),
    ('PEN', 'pen', 10),
    ('POWERBANK', 'powerbank', 200),
    ('HOODY', 'hoody', 300),
    ('UMBRELLA', 'umbrella', 200),
    ('SOCKS', 'socks', 10),
    ('WALLET', 'wallet', 50),
    ('PINK_HOODY', 'pink-hoody', 500);

-- Тестовые данные для пользователей
INSERT INTO auth_users (login, password) VALUES
    ('user1', 'hashed_password1'),
    ('user2', 'hashed_password2');

INSERT INTO info_users (user_id, name, balance) VALUES
    (1, 'User One', 1000),
    (2, 'User Two', 1000);
