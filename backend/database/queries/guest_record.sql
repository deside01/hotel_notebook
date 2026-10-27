-- name: CreateGuestRecord :one
INSERT INTO guest_records (
    user_id,
    full_name,
    phone,
    room,
    check_in,
    check_out,
    note
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetGuestRecordByID :one
SELECT * FROM guest_records
WHERE id = $1;

-- name: GetGuestRecordsByUserID :many
SELECT * FROM guest_records
WHERE user_id = $1;

-- name: UpdateGuestRecord :exec
UPDATE guest_records
SET check_out = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeleteGuestRecordByID :exec
DELETE FROM guest_records
WHERE id = $1;

-- name: SearchGuests :many
SELECT * FROM guest_records
WHERE 
    (sqlc.arg(full_name)::text = '' OR full_name ILIKE CONCAT('%', sqlc.arg(full_name)::text, '%'))
    AND (sqlc.arg(phone)::text = '' OR phone LIKE CONCAT('%', sqlc.arg(phone)::text, '%'))
    AND (sqlc.arg(room)::text = '' OR room ILIKE CONCAT(sqlc.arg(room)::text, '%'))
ORDER BY created_at DESC
LIMIT sqlc.arg(req_limit)::int
OFFSET sqlc.arg(req_offset)::int;

