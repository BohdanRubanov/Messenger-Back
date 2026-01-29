-- Drop an existing table 'TableName'
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;


CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL Unique,
    name VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL
);



