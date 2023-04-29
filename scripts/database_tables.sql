CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY,
    account_number VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY,
    account_id INTEGER NOT NULL,
    operation_type INTEGER,
    amount NUMERIC(10,2),
    created_at TEXT,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
);
