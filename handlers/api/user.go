package api

import (
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleUpdateUser(w http.ResponseWriter, r *http.Request, repo *repos.Repository) error {
	param_username := chi.URLParam(r, "username")

	r.ParseForm()
	firstname := pgtype.Text{String: r.Form.Get("firstname"), Valid: r.Form.Has("firstname")}
	lastname := pgtype.Text{String: r.Form.Get("lastname"), Valid: r.Form.Has("lastname")}

	err := repo.Queries.UpdateUser(repo.Ctx, queries.UpdateUserParams{
		Username:  param_username,
		Firstname: firstname,
		Lastname:  lastname,
	})
	if err != nil {
		return err
	}

	w.Header().Add("HX-Refresh", "true")
	return nil
}
