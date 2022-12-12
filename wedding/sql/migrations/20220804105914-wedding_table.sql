
-- +migrate Up
CREATE TABLE weddings (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    names varchar NOT NULL,
    status int8 NOT NULL default 0,
    user_id uuid NOT NULL,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    wedding_date timestamp NOT NULL

);

-- +migrate Down
DROP table weddings;