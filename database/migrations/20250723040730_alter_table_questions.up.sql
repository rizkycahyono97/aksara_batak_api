-- Step 1: Ubah question_type menjadi VARCHAR sementara (hindari error enum)
ALTER TABLE questions MODIFY COLUMN question_type VARCHAR(50) NOT NULL;

-- Step 2: Update semua nilai ke value valid
UPDATE questions SET question_type = 'pilihan_ganda_aksara' WHERE question_type IN ('multiple_choice');
UPDATE questions SET question_type = 'pilihan_ganda_batak' WHERE question_type IN ('image_choice', 'listen_and_type');
UPDATE questions SET question_type = 'nulis_aksara' WHERE question_type = 'drawing';

-- Optional: Tangani NULL atau kosong
UPDATE questions SET question_type = 'pilihan_ganda_batak' WHERE question_type IS NULL OR question_type = '';

-- Optional: Hapus data yang masih tidak valid (safety check)
DELETE FROM questions
WHERE question_type NOT IN ('pilihan_ganda_aksara', 'pilihan_ganda_batak', 'nulis_aksara');

-- Step 3: Ubah kembali menjadi ENUM
ALTER TABLE questions MODIFY COLUMN question_type ENUM(
    'pilihan_ganda_aksara',
    'pilihan_ganda_batak',
    'nulis_aksara'
    ) NOT NULL DEFAULT 'pilihan_ganda_batak';
