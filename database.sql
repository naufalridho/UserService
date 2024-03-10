/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE TABLE users (
	id UUID PRIMARY KEY,
	full_name VARCHAR ( 50 ) NOT NULL,
    phone_number VARCHAR ( 13 ) NOT NULL,
    hashed_password VARCHAR (64) NOT NULL,
    login_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_phone_number ON users ("phone_number");
