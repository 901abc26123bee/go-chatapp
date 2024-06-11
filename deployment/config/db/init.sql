-- Ensure the postgres user is created
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
    CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'postgres';
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(255) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255),
  encrypted_password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS group_chats (
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  user_counts INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS group_chat_users (
  id SERIAL NOT NULL PRIMARY KEY,
  group_chat_id INTEGER REFERENCES group_chats(id),
  user_id VARCHAR(255) REFERENCES users(id),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS chat_users (
  id SERIAL NOT NULL PRIMARY KEY,
  user1_id VARCHAR(255) REFERENCES users(id),
  user2_id VARCHAR(255) REFERENCES users(id),
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS group_chat_messages (
  id SERIAL NOT NULL PRIMARY KEY,
  group_chat_id INTEGER REFERENCES group_chats(id),
  sent_at INTEGER NOT NULL,
  text VARCHAR(255) DEFAULT '',
  storage_path VARCHAR(255) DEFAULT '',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS user_chat_messages (
  id SERIAL NOT NULL PRIMARY KEY,
  user_chat_id INTEGER REFERENCES chat_users(id),
  sender_id VARCHAR(255) REFERENCES users(id),
  sent_at TIMESTAMPTZ DEFAULT now(),
  text VARCHAR(255) DEFAULT '',
  storage_path VARCHAR(255) DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  deleted BOOLEAN DEFAULT false
);
