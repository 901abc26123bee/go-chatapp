-- Ensure the postgres user is created
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
    CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'postgres';
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(256),
  uid VARCHAR(256),
  email VARCHAR(256),
  password VARCHAR(256)
);

USER_RELATIONSHIP {
    user_first_id,
    user_second_id,
    type

    primary key(user_first_id, user_second_id)
}