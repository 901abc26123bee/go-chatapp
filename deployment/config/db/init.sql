-- Ensure the postgres user is created
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
    CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'postgres';
  END IF;
END
$$;
