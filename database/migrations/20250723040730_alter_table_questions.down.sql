ALTER TABLE questions
    MODIFY COLUMN question_type ENUM(
    'multiple_choice',
    'image_choice',
    'drawing',
    'listen_and_type'
    ) NOT NULL DEFAULT 'multiple_choice';