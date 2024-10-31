-- name: ListClubsOrderByNameAsc :many
SELECT *
FROM "club"
ORDER BY "club".name ASC;

-- name: FindClubByID :one
SELECT *
FROM "club"
WHERE "club".id = $1
LIMIT 1;