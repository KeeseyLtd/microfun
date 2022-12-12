
-- +migrate Up
ALTER TABLE invitations
ADD wedding_id uuid NOT NULL,
ADD CONSTRAINT fk_invitations_weddings FOREIGN KEY (wedding_id) REFERENCES weddings(id)
ON DELETE CASCADE;

-- +migrate Down
ALTER TABLE invitations
DROP COLUMN wedding_id;