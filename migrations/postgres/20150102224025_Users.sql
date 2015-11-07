
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users
(
  id uuid primary key default uuid_generate_v4(),
  lastname text NOT NULL,
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  tokenemail text NOT NULL,
  tokenexpires timestamp with time zone NOT NULL,
  logincount integer NOT NULL,
  lastlogin timestamp with time zone,
  isactive boolean NOT NULL,
  issystemadministrator boolean NOT NULL,
  datecreated timestamp with time zone NOT NULL,
  lastmodified timestamp with time zone NULL,
  CONSTRAINT users_email_key UNIQUE (email)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

drop TABLE users;