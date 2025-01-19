-- ENUM for transaction types
CREATE TYPE transaction_type_enum AS ENUM ('top_up', 'transfer');

-- Users table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  fullname VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  pin VARCHAR(255),
  phone VARCHAR(16) UNIQUE,
  image VARCHAR(255),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);


-- Wallets table
CREATE TABLE wallets (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  balance DECIMAL(10, 2) NOT NULL CHECK (balance >= 0),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

-- Transactions table
CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
  transaction_type transaction_type_enum NOT NULL,
  notes VARCHAR(255),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

-- Payment Methods table
CREATE TABLE payment_methods (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  tax DECIMAL(10, 2) NOT NULL CHECK (tax >= 0),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

-- Transfer Transactions table
CREATE TABLE transfer_transactions (
  id SERIAL PRIMARY KEY,
  transaction_id INT REFERENCES transactions(id),
  target_user_id INT REFERENCES users(id),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

-- Top-Up Transactions table
CREATE TABLE topup_transactions (
  id SERIAL PRIMARY KEY,
  transaction_id INT REFERENCES transactions(id),
  payment_method_id INT REFERENCES payment_methods(id),
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);