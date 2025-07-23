CREATE TABLE lessons (
                         id INT UNSIGNED NOT NULL AUTO_INCREMENT,
                         title VARCHAR(255) NOT NULL,
                         description TEXT NULL,
                         icon_url VARCHAR(255) NULL,
                         order_index INT NOT NULL DEFAULT 0,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         PRIMARY KEY (id)
);

ALTER TABLE quizzes
    ADD COLUMN lesson_id INT UNSIGNED NULL AFTER id,
    ADD CONSTRAINT fk_quizzes_lesson
    FOREIGN KEY (lesson_id)
    REFERENCES lessons(id)
    ON DELETE SET NULL;
