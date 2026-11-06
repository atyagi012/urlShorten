CREATE TABLE IF NOT EXISTS urls (
    code         VARCHAR(10) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);
