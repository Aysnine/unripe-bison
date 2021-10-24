SET sql_safe_updates = FALSE;

USE defaultdb;

DROP DATABASE IF EXISTS unripe_bison CASCADE;
CREATE DATABASE IF NOT EXISTS unripe_bison;

USE unripe_bison;

CREATE TABLE books (
  id UUID PRIMARY KEY,
  name TEXT
);
