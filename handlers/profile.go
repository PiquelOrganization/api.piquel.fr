package handlers

import (
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleProfileQuery(w http.ResponseWriter, r *http.Request) {
    // Get username from query params. Should look likes "GET api.piquel.fr/profile?[username]
    username := ""

    writeProfile(w, r, username)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]

    writeProfile(w, r, username)
}

func writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := database.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	}

	profile := &types.UserProfile{User: user}

	group, err := database.Queries.GetGroupInfo(r.Context(), user.Group)
	if err != nil {
		panic(err)
	}

	profile.UserColor = group.Color
	profile.UserGroup = group.Displayname.String

    // Write json of the object to the requesting client
}
