-- name: GetInvitation :one
SELECT * FROM invitations
WHERE wedding_id = $1 AND id = $2 LIMIT 1;

-- name: GetInvitations :many
SELECT * FROM invitations
WHERE wedding_id = $1;

-- name: CreateInvitation :one
INSERT INTO invitations (
    id, invitees, wedding_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateInvitation :exec
UPDATE invitations
SET
	status = CASE WHEN @status_do_update::boolean
		THEN @status::int ELSE status END,
	invitees = CASE WHEN @invitees_do_update::boolean
		THEN @invitees::varchar ELSE invitees END
WHERE
	id = @id;

-- name: DeleteInvitation :exec
DELETE FROM invitations WHERE id = $1;
