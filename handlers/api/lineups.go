package api

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
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

	club_id := r.FormValue("club_id")

	mirror, err := strconv.ParseBool(r.FormValue("mirror"))
	if err != nil {
		return err
	}

	player_id, err := repo.Queries.FindPlayerIDByClubAndNo(repo.Ctx, queries.FindPlayerIDByClubAndNoParams{
		ClubID: pgtype.Text{String: club_id, Valid: true},
		No:     int16(no),
	})
	if err != nil {
		return err
	}

	position_no, err := strconv.Atoi(r.FormValue("position_no"))
	if err != nil {
		return err
	}

	repo.Queries.CreateLineupPlayer(repo.Ctx, queries.CreateLineupPlayerParams{
		LineupID:   int32(lineupID),
		PlayerID:   player_id,
		PositionNo: int16(position_no),
		Position:   queries.PlayerPosition(r.FormValue("position")),
	})

	matchID, err := repo.Queries.FindMatchIDFromLineupID(repo.Ctx, int32(lineupID))
	if err != nil {
		return err
	}

	match, err := repo.Queries.FindMatchByID(repo.Ctx, matchID)
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

	club_name := match.HomeClubName
	lineupPlayers := homeLineupPlayers
	if club_id == match.AwayClubID {
		club_name = match.AwayClubName
		lineupPlayers = awayLineupPlayers
	}

	return handlers.Render(w, r, admin_lineup_components.AddLineupPlayerFormResponse(admin_lineup_components.AddLineupPlayerParams{
		ClubID:        club_id,
		ClubName:      club_name,
		LineupID:      int32(lineupID),
		LineupPlayers: lineupPlayers,
		Mirror:        mirror,
	}, match, homeLineupPlayers, awayLineupPlayers))
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
	if r.FormValue("position_no") == "" {
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
	if r.FormValue("position") == "" {
		position.PlayerPosition = queries.PlayerPosition(r.FormValue("position"))
		position.Valid = true
	}

	var goals = pgtype.Int2{
		Int16: 0,
		Valid: false,
	}
	if r.FormValue("goals") == "" {
		goals_value, err := strconv.Atoi(r.FormValue("goals"))
		if err != nil {
			return err
		}
		goals.Int16 = int16(goals_value)
		goals.Valid = true
	}

	var yellow_cards = pgtype.Int2{
		Int16: 0,
		Valid: false,
	}
	if r.FormValue("yellow_cards") == "" {
		yellow_cards_value, err := strconv.Atoi(r.FormValue("yellow_cards"))
		if err != nil {
			return err
		}
		yellow_cards.Int16 = int16(yellow_cards_value)
		yellow_cards.Valid = true
	}

	var red_cards = pgtype.Int2{
		Int16: 0,
		Valid: false,
	}
	if r.FormValue("red_cards") == "" {
		red_cards_value, err := strconv.Atoi(r.FormValue("red_cards"))
		if err != nil {
			return err
		}
		red_cards.Int16 = int16(red_cards_value)
		red_cards.Valid = true
	}

	repo.Queries.UpdateLineupPlayer(repo.Ctx, queries.UpdateLineupPlayerParams{
		LineupID:    int32(lineupID),
		PlayerID:    int32(playerID),
		PositionNo:  position_no,
		Position:    position,
		Goals:       goals,
		YellowCards: yellow_cards,
		RedCards:    red_cards,
	})

	return nil
}
