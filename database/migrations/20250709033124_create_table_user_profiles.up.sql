CREATE TABLE user_profiles (
    user_id VARCHAR(36) NOT NULL,
    total_xp INT NOT NULL DEFAULT 0,
    current_streak INT NOT NULL DEFAULT 0,
    last_active_at DATE NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(uuid) ON DELETE CASCADE
);