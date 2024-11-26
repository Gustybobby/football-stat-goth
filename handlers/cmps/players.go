package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views/components/table_components"
	"net/http"
)

func HandlePlayersTable(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	players, err := repo.Queries.ListPlayerLikeFullname(repo.Ctx, queries.ListPlayerLikeFullnameParams{
		FullnameLike: "%" + r.FormValue("fullname") + "%",
		Offset:       0,
		Limit:        30,
	})
	if err != nil {
		return err
	}

	return handlers.Render(w, r, table_components.PlayersTable(players))
}
