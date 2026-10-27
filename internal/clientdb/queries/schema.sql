CREATE TABLE IF NOT EXISTS clients(
    id INTEGER NOT NULL PRIMARY KEY ,
    client_id TEXT NOT NULL,
    rps REAL NOT NULL,
    burst INTEGER NOT NULL,
    UNIQUE(client_id)
)