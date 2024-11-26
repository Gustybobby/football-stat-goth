package api

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views/admin/admin_components/admin_lineup_components"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreateLineupPlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	lineupID, err := strconv.Atoi(chi.URLParam(r, "lineupID"))
	if err != nil {
		return err
	}

	no, err := strconv.Atoi(r.FormValue("no"))
	if err != nil {
		return err
	}

	player_id, err := repo.Queries.FindPlayerIDByClubNoSeason(repo.Ctx, queries.FindPlayerIDByClubNoSeasonParams{
		ClubID: r.FormValue("club_id"),
		No:     int16(no),
		Season: pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	position_no, err := strconv.Atoi(r.FormValue("position_no"))
	if err != nil {
		return err
	}

	lineupPlayer, err := repo.Queries.CreateLineupPlayer(repo.Ctx, queries.CreateLineupPlayerParams{
		LineupID:   int32(lineupID),
		PlayerID:   player_id,
		PositionNo: int16(position_no),
		Position:   queries.PlayerPosition(r.FormValue("position")),
	})
	if err != nil {
		return err
	}

	return handleFormResponse(lineupPlayer.LineupID, w, r, repo)
}

func HandleUpdateLineupPlayer(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	lineupID, err := strconv.Atoi(chi.URLParam(r, "lineupID"))
	if err != nil {
		return err
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	var position_no = pgtype.Int2{
		Int16: 0,
		Valid: false,
	}
	if r.FormValue("position_no") != "" {
		pos_no_value, err := strconv.Atoi(r.FormValue("position_no"))
		if err != nil {
			return err
		}
		position_no.Int16 = int16(pos_no_value)
		position_no.Valid = true
	}

	var position = queries.NullPlayerPosition{
		PlayerPosition: queries.PlayerPositionGK,
		Valid:          false,
	}
	if r.FormValue("position") != "" {
		position.PlayerPosition = queries.PlayerPosition(r.FormValue("position"))
		position.Valid = true
	}

	lineupPlayer, err := repo.Queries.UpdateLineupPlayer(repo.Ctx, queries.UpdateLineupPlayerParams{
		LineupID:   int32(lineupID),
		PlayerID:   int32(playerID),
		PositionNo: position_no,
		Position:   position,
	})
	if err != nil {
		return err
	}

	return handleFormResponse(lineupPlayer.LineupID, w, r, repo)
}

func handleFormResponse(lineupID int32, w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	matchID, err := repo.Queries.FindMatchIDFromLineupID(repo.Ctx, lineupID)
	if err != nil {
		return err
	}

	match, err := repo.Queries.FindMatchByID(repo.Ctx, matchID)
	if err != nil {
		return err
	}

	events, err := repo.Queries.ListLineupEventsByMatchID(repo.Ctx, int32(matchID))
	if err != nil {
		return err
	}

	homeLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.HomeLineupID)
	if err != nil {
		return err
	}

	awayLineupPlayers, err := repo.Queries.ListLineupPlayersByLineupID(repo.Ctx, match.AwayLineupID)
	if err != nil {
		return err
	}

	club_id := match.HomeClubID
	club_name := match.HomeClubName
	lineupPlayers := homeLineupPlayers
	mirror := false

	if lineupID == match.AwayLineupID {
		club_id = match.AwayClubID
		club_name = match.AwayClubName
		lineupPlayers = awayLineupPlayers
		mirror = true
	}

	return handlers.Render(w, r, admin_lineup_components.LineupPlayerFormResponse(admin_lineup_components.AddLineupPlayerParams{
		ClubID:        club_id,
		ClubName:      club_name,
		LineupID:      int32(lineupID),
		LineupPlayers: lineupPlayers,
		Mirror:        mirror,
	}, match, events, homeLineupPlayers, awayLineupPlayers))
}
