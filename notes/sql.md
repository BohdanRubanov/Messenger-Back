
```sql
-- ----------------------------
-- Common PostgreSQL Data Types
-- ----------------------------

-- Integer types
CREATE TABLE numbers_example (
    small_num SMALLINT,       -- small integer (-32,768 … 32,767)
    normal_num INTEGER,       -- standard integer (-2,147,483,648 … 2,147,483,647)
    big_num BIGINT,           -- big integer, for IDs or counters
    auto_id SERIAL,           -- auto-increment integer ID
    big_auto_id BIGSERIAL     -- auto-increment big integer ID
);

-- Floating point and numeric types
CREATE TABLE money_example (
    price NUMERIC(10,2),     -- fixed-precision number (e.g., 12345678.90)
    rating REAL,              -- single precision float (4 bytes)
    score DOUBLE PRECISION    -- double precision float (8 bytes)
);

-- String types
CREATE TABLE text_example (
    short_text VARCHAR(255),  -- short string, e.g., username
    long_text TEXT,           -- long text, e.g., message content
    uuid_id UUID              -- universally unique identifier
);

-- Boolean type
CREATE TABLE boolean_example (
    is_active BOOLEAN         -- true / false flag
);

-- Date and time types
CREATE TABLE datetime_example (
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  -- timestamp with timezone
    birthday DATE,                                     -- only date
    login_time TIME                                    -- only time
);

-- JSON type
CREATE TABLE json_example (
    settings JSONB            -- JSON with indexing support, e.g., user settings or media metadata
);

-- ----------------------------
-- SQL Commands Cheat Sheet (PostgreSQL)
-- ----------------------------

-- CREATE TABLE
CREATE TABLE IF NOT EXISTS TableName (
    id SERIAL PRIMARY KEY,       -- auto-increment integer ID
    name VARCHAR(255),           -- short text
    description TEXT,            -- long text
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- timestamp with timezone
);

-- INSERT INTO
INSERT INTO TableName (name, description)
VALUES ('Sample Name', 'Sample description');

-- SELECT (Read data)
SELECT * FROM TableName;                   -- read all columns
SELECT id, name FROM TableName;           -- read specific columns
SELECT * FROM TableName WHERE id = 1;     -- read with condition
SELECT * FROM TableName ORDER BY id DESC; -- sort results
SELECT * FROM TableName LIMIT 10;         -- limit results
SELECT * FROM TableName OFFSET 10;        -- skip first N rows

-- UPDATE (Modify data)
UPDATE TableName
SET name = 'New Name'
WHERE id = 1;

-- DELETE (Remove data)
DELETE FROM TableName
WHERE id = 1;

-- ALTER TABLE (Change table structure)
ALTER TABLE TableName ADD COLUMN price INT;           -- add new column
ALTER TABLE TableName DROP COLUMN description;        -- remove column
ALTER TABLE TableName ALTER COLUMN name TYPE TEXT;    -- change column type

-- DROP TABLE (Remove table completely)
DROP TABLE IF EXISTS TableName;

-- CREATE INDEX (Speed up queries)
CREATE INDEX idx_name ON TableName(name);    -- simple index
CREATE UNIQUE INDEX idx_unique_name ON TableName(name);  -- unique index

-- Transactions (for safety)
BEGIN;                                       -- start transaction
INSERT INTO TableName (name) VALUES ('Test');
UPDATE TableName SET name = 'Updated' WHERE id = 1;
COMMIT;                                      -- commit changes
ROLLBACK;                                    -- cancel transaction

-- Special / Advanced
-- UPSERT (Insert or Update)
INSERT INTO TableName (id, name)
VALUES (1, 'Name')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- JSONB operations
CREATE TABLE json_example (
    data JSONB
);
INSERT INTO json_example (data)
VALUES ('{"theme":"dark","notifications":true}');
SELECT data->>'theme' FROM json_example;  -- extract JSON field

-- ARRAY example
CREATE TABLE array_example (
    tags TEXT[]
);
INSERT INTO array_example (tags) VALUES (ARRAY['go','postgres','docker']);
SELECT * FROM array_example WHERE 'go' = ANY(tags);

-- ----------------------------
-- Meta-commands in psql
-- ----------------------------
\l           -- list all databases
\c dbname    -- connect to a database
\dt          -- list tables
\d TableName -- describe table structure
\q           -- quit psql


-- ----------------------------
-- Relationships / Joins Cheat Sheet
-- ----------------------------

-- One-to-Many (1:N)
-- Example: One product can have many users
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    product_id INT,  -- foreign key to products
    CONSTRAINT fk_product
        FOREIGN KEY (product_id) REFERENCES products(id)
        ON DELETE CASCADE  -- if product deleted, related users are deleted
);

-- SELECT with JOIN
SELECT u.id AS user_id, u.name AS user_name, p.name AS product_name
FROM users u
JOIN products p ON u.product_id = p.id;

-- One-to-One (1:1)
-- Example: One user has one profile
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,  -- one-to-one relationship
    bio TEXT,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);

-- Many-to-Many (M:N)
-- Example: Users belong to many chats
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE chat_members (
    chat_id INT,
    user_id INT,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (chat_id, user_id),  -- composite primary key
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- SELECT with Many-to-Many
SELECT u.name AS user_name, c.name AS chat_name
FROM chat_members cm
JOIN users u ON cm.user_id = u.id
JOIN chats c ON cm.chat_id = c.id;

-- Different JOIN types
-- INNER JOIN: only matching rows
SELECT * FROM users u
INNER JOIN products p ON u.product_id = p.id;

-- LEFT JOIN: all rows from left, matching from right
SELECT * FROM users u
LEFT JOIN products p ON u.product_id = p.id;

-- RIGHT JOIN: all rows from right, matching from left
SELECT * FROM users u
RIGHT JOIN products p ON u.product_id = p.id;

-- FULL OUTER JOIN: all rows from both tables
SELECT * FROM users u
FULL OUTER JOIN products p ON u.product_id = p.id;

-- CROSS JOIN: Cartesian product (all combinations)
SELECT * FROM users u
CROSS JOIN products p;

-- SELF JOIN: joining table with itself
SELECT u1.name AS user1, u2.name AS user2
FROM users u1
JOIN users u2 ON u1.id < u2.id;

-- Foreign key options
-- ON DELETE CASCADE      -> delete child rows if parent deleted
-- ON DELETE SET NULL     -> set foreign key to NULL if parent deleted
-- ON UPDATE CASCADE      -> update child rows if parent ID changes

-- Composite keys
CREATE TABLE user_chats (
    user_id INT,
    chat_id INT,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, chat_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES chats(id)
);

-- Check relationships
-- Ensure column satisfies some condition
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    sender_id INT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_text_length CHECK (char_length(text) > 0)
);

-- Sample SELECT with multiple joins (messages + users + chats)
SELECT m.id AS message_id,
       u.name AS sender_name,
       c.name AS chat_name,
       m.text,
       m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
JOIN chats c ON m.chat_id = c.id
ORDER BY m.created_at DESC;

-- ----------------------------
-- Notes
-- ----------------------------
-- 1. Define foreign keys in the child table, pointing to parent table
-- 2. Use CASCADE / SET NULL / NO ACTION as needed
-- 3. Composite keys are useful for M:N relationships
-- 4. Always use JOINs to read data across related tables
