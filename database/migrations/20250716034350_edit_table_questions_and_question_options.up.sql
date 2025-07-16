ALTER TABLE questions
ADD COLUMN question_type  ENUM('multiple_choice', 'image_choice', 'listen_and_type', 'drawing') NOT NULL DEFAULT 'multiple_choice' AFTER quiz_id,
ADD COLUMN image_url VARCHAR(255) NULL AFTER question_text,
ADD COLUMN audio_url VARCHAR(255) NULL AFTER image_url,
ADD COLUMN lottie_url VARCHAR(255) NULL AFTER audio_url;

ALTER TABLE question_options
ADD COLUMN aksara_text TEXT NULL AFTER option_text,
ADD COLUMN image_url VARCHAR(255) NULL AFTER aksara_text,
ADD COLUMN audio_url VARCHAR(255) NULL AFTER image_url;