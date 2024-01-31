CREATE TABLE IF NOT EXISTS file_event (
    path TEXT PRIMARY KEY,
    event INT CHECK (event IN (0, 1, 2, 3, 4, 5)),
    timestamp DATETIME
);