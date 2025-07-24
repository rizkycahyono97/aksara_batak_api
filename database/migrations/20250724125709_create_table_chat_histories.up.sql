CREATE TABLE chat_histories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id VARCHAR(36) NOT NULL,
    role VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
        CONSTRAINT fk_users
        FOREIGN KEY(user_id)
        REFERENCES users(uuid)
        ON DELETE CASCADE
);

CREATE INDEX idx_chat_histories_user_id ON chat_histories(user_id);