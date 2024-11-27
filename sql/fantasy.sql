-- name: GetFantasy_PlayerInfoForFantasy :many
SELECT "player".lastname, "player".position, "player".image, "club".id as club_id, "fantasy_player".cost, "fantasy_player".points, "fantasy_player".rating
FROM "fantasy_player" 
JOIN "player" on "fantasy_player".player_id = "player".id
JOIN "club" on "fantasy_player".club_id = "club".id
ORDER BY "player".position ASC, "player".lastname ASC;