CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    name VARCHAR(255),
    age INTEGER,
    is_deleted BOOLEAN
);

INSERT INTO users (email, password, name, age, is_deleted)
VALUES ('admin@example.com', 'password', 'Admin', 34, FALSE),
       ('john.doe@gmail.com', 'password', 'John Doe', 27, FALSE),
       ('jen.star@gmail.com', 'password', 'Jennifer', 31, FALSE)
;

