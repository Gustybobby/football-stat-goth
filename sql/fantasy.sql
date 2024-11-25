-- name: GetPlayerInfoForFantasy :many
SELECT "player".lastname, "player".position, "player".image, "club".id as club_id
FROM "player" 
inner join "club_player" on "player".id = "club_player".player_id
inner join "club" on "club_player".club_id = "club".id
ORDER BY "player".position ASC, "player".lastname ASC;

-- name: GetFantasy_PlayerInfoForFantasy :many
SELECT "player".lastname, "player".position, "player".image, "club".id as club_id, "fantasy_player".cost, "fantasy_player".points, "fantasy_player".rating
FROM "fantasy_player" 
inner join "player" on "fantasy_player".player_id = "player".id
inner join "club" on "fantasy_player".club_id = "club".id
ORDER BY "player".position ASC, "player".lastname ASC;