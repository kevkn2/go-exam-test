-- Table: online_tests
CREATE TABLE online_tests (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    test_id VARCHAR(255) NOT NULL,
    duration INT NOT NULL
);

-- Table: questions
CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    type VARCHAR(50) NOT NULL,
    text TEXT NOT NULL,
    answer TEXT,
    meta JSONB,
    "order" INT NOT NULL,

    online_test_id BIGINT,
    CONSTRAINT fk_questions_online_test FOREIGN KEY (online_test_id) 
        REFERENCES online_tests (id) ON DELETE CASCADE
);

-- Index for soft delete
CREATE INDEX idx_online_tests_deleted_at ON online_tests (deleted_at);
CREATE INDEX idx_questions_deleted_at ON questions (deleted_at);

-- Optional: faster query for online_test_id lookup
CREATE INDEX idx_questions_online_test_id ON questions (online_test_id);
