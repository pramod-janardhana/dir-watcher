CREATE TABLE IF NOT EXISTS job_runs (
    run_id TEXT PRIMARY KEY,
    job_type INT CHECK (job_type IN (1)),
    start_time DATETIME,
    end_time DATETIME,
    status INT CHECK (status IN (1, 2, 3)),
    details TEXT
);