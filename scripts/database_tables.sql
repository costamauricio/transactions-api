CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY,
    document_number VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY,
    account_id INTEGER NOT NULL,
    operation_type INTEGER NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    created_at TEXT,
    FOREIGN KEY(account_id) REFERENCES accounts(id)
);
