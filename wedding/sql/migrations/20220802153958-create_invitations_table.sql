
-- +migrate Up
CREATE TABLE invitations (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    invitees varchar NOT NULL,
    status int8 NOT NULL default 0,
    created_at timestamp default now()
);

-- +migrate Down
DROP TABLE invitations;