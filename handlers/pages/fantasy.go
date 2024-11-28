package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/handlers/api"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/plconstant"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views"
	"net/http"
	"net/url"
)

func HandleFantasyPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)
	if user == nil {
		http.Redirect(w, r, "/signin?redirectUrl="+url.QueryEscape("/fantasy"), http.StatusFound)
		return nil
	}

	fixtures, err := repo.Queries.ListMatchesWithClubsAndGoals(repo.Ctx, queries.ListMatchesWithClubsAndGoalsParams{
		FilterClubID: false,
		FilterWeek:   false,
		IsFinished:   false,
		Order:        "ASC",
	})
	if err != nil {
		return err
	}

	players, err := repo.Queries.ListFantasyPlayers(repo.Ctx, queries.ListFantasyPlayersParams{
		FilterFantasyPlayerID: false,
		MinCost:               plconstant.FantasyPlayerMinCost,
		AvgCost:               plconstant.FantasyPlayerAverageCost,
		Season:                pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	fantasy_team_player_refs, err := repo.Queries.ListFantasyTeamPlayersByUsernameSeason(repo.Ctx, queries.ListFantasyTeamPlayersByUsernameSeasonParams{
		Username: user.Username,
		Season:   pltime.GetCurrentSeasonString(),
	})
	if err != nil {
		return err
	}

	var fantasy_team_players []queries.ListFantasyPlayersRow
	for _, player := range players {
		for _, team_player_ref := range fantasy_team_player_refs {
			if player.ID == team_player_ref.FantasyPlayerID {
				fantasy_team_players = append(fantasy_team_players, player)
				break
			}
		}
	}

	players_params, _, err := api.GetFantasyTeamFieldParams(fantasy_team_players)
	if err != nil {
		return err
	}

	budget := plconstant.FantasyTeamMaxBudget
	if len(fantasy_team_player_refs) > 0 {
		budget = int(fantasy_team_player_refs[0].Budget)
	}

	return handlers.Render(w, r, views.Fantasy(user, fixtures, players, *players_params, budget))
}
