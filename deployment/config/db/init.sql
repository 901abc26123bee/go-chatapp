-- Ensure the postgres user is created
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
    CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'postgres';
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS TABLE users (
  id SERIAL PRIMARY KEY,
  uid VARCHAR(256),
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  password VARCHAR,
  encrypted_password VARCHAR NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS group_chats (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  user_counts INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS group_chat_users (
  id SERIAL PRIMARY KEY,
  group_chat_id INTEGER REFERENCES group_chats(id),
  user_id INTEGER REFERENCES users(id),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS chat_users (
  id SERIAL PRIMARY KEY,
  user1_id INTEGER REFERENCES users(id),
  user2_id INTEGER REFERENCES users(id),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS group_chat_messages (
  id SERIAL PRIMARY KEY,
  group_chat_id INTEGER REFERENCES group_chats(id),
  sent_at INTEGER NOT NULL,
  text VARCHAR,
  storage_path VARCHAR,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS user_chat_messages (
  id SERIAL PRIMARY KEY,
  user_chat_id INTEGER REFERENCES chat_users(id),
  sender_id INTEGER REFERENCES users(id),
  sent_at INTEGER,
  text VARCHAR,
  storage_path VARCHAR,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted BOOLEAN
);
