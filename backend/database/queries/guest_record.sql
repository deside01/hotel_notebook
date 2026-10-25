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
SELECT 
    id,
    user_id,
    phone
FROM guest_records
WHERE id = $1;