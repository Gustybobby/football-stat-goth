package pages

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/services/plauth"
	"football-stat-goth/services/pltime"
	"football-stat-goth/views"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HandlePlayerPage(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	user := plauth.GetContextUser(r)

	player_id, err := strconv.Atoi(chi.URLParam(r, "playerID"))
	if err != nil {
		return err
	}

	player, err := repo.Queries.FindPlayerByID(repo.Ctx, int32(player_id))
	if err != nil {
		return err
	}

	club_players, err := repo.Queries.ListClubPlayerByPlayerID(repo.Ctx, player.ID)
	if err != nil {
		return err
	}

	season := pltime.GetCurrentSeasonString()

	performance, err := repo.Queries.ListPlayerSeasonPerformance(repo.Ctx, queries.ListPlayerSeasonPerformanceParams{
		Season:         season,
		FilterPlayerID: true,
		PlayerID:       player.ID,
		FilterClubID:   false,
		ClubID:         "",
		Limit:          1,
	})
	if err != nil {
		return err
	}

	matches, err := repo.Queries.ListPlayerMatchHistory(repo.Ctx, player.ID)
	if err != nil {
		return err
	}

	return handlers.Render(w, r, views.Player(user, player, club_players, performance[0], matches))
}
