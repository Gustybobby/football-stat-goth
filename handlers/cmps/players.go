package cmps

import (
	"football-stat-goth/handlers"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"football-stat-goth/views/components/table_components"
	"net/http"
	"strconv"
)

func HandlePlayersTable(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	page_size := 30

	players, err := repo.Queries.ListPlayerLikeFullname(repo.Ctx, queries.ListPlayerLikeFullnameParams{
		FullnameLike: "%" + r.FormValue("fullname") + "%",
		Offset:       int32((page - 1) * page_size),
		Limit:        int32(page_size),
	})
	if err != nil {
		return err
	}

	return handlers.Render(w, r, table_components.PlayersTable(players, page))
}
