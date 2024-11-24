-- name: GetPlayerInfoForFantasy :many
SELECT "player".lastname, "player".position, "player".image, "club".id as club_id
FROM "player" 
inner join "club_player" on "player".id = "club_player".player_id
inner join "club" on "club_player".club_id = "club".id
ORDER BY "player".position ASC;