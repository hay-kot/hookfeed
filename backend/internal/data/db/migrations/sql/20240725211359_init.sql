-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Source:
--  - https://gist.github.com/kjmph/5bd772b2c2df145aa645b837da7eca74
CREATE FUNCTION uuid_generate_v7() RETURNS uuid AS
$$
BEGIN
  -- use random v4 uuid as starting point (which has the same variant we need)
  -- then overlay timestamp
  -- then set version 7 by flipping the 2 and 1 bit in the version 4 string
  return encode(
    set_bit(
      set_bit(
        overlay(uuid_send(gen_random_uuid())
                placing substring(int8send(floor(extract(epoch from clock_timestamp()) * 1000)::bigint) from 3)
                from 1 for 6
        ),
        52, 1
      ),
      53, 1
    ), 'hex')::uuid;
END
$$
language plpgsql volatile;

CREATE TABLE IF NOT EXISTS users (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v7 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- user identity
    username VARCHAR(100) NOT NULL DEFAULT '',
    email VARCHAR(254) NOT NULL UNIQUE,
    password_hash VARCHAR(500) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    -- billing information
    stripe_customer_id VARCHAR(50),
    stripe_subscription_id VARCHAR(50),
    subscription_start_date TIMESTAMP,
    subscription_ended_date TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

DROP FUNCTION uuid_generate_v7;
-- +goose StatementEnd
