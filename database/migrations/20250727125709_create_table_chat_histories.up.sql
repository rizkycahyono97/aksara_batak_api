CREATE TABLE chat_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    message TEXT NOT NULL,
    reply TEXT NOT NULL,
    message_type VARCHAR(30) DEFAULT 'text',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        CONSTRAINT fk_users
        FOREIGN KEY(user_id)
        REFERENCES users(uuid)
        ON DELETE CASCADE
);

CREATE INDEX idx_chat_histories_user_id ON chat_histories(user_id);