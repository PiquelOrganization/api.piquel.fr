package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/markbates/goth/gothic"
)

const RedirectSession = "redirect_to"

func HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	saveRedirectURL(w, r)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	_, err = users.VerifyUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

    redirectUser(w, r)
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	username, err := users.VerifyUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.StoreUserSession(w, r, username, types.UserSessionFromGothUser(&user))
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	redirectUser(w, r)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Error authencticating", http.StatusInternalServerError)
		panic(err)
	}

	err = auth.RemoveUserSession(w, r)
	if err != nil {
		http.Error(w, "Error removing cookies", http.StatusInternalServerError)
		panic(err)
	}

	redirectURL := getRedirectURL(r)
    http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func getRedirectURL(r *http.Request) string {
	redirectTo := r.URL.Query().Get("redirectTo")

    return fmt.Sprintf("https://%s/%s", config.Envs.Domain, strings.Trim(redirectTo, "/"))
}

func saveRedirectURL(w http.ResponseWriter, r *http.Request) {
	redirectURL := getRedirectURL(r)

	session, err := gothic.Store.Get(r, RedirectSession)
	if err != nil {
		panic(err)
	}

	session.Values["redirectTo"] = redirectURL

    err = session.Save(r, w)
    if err != nil {
        panic(err)
    }
}

func redirectUser(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, RedirectSession)
	if err != nil {
		panic(err)
	}

    redirectURL := session.Values["redirectTo"]
    session.Values["redirectTo"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

    if redirectURL == nil || redirectURL == "" {
        redirectURL = fmt.Sprintf("https://%s", config.Envs.Domain)
    }

    http.Redirect(w, r, redirectURL.(string), http.StatusTemporaryRedirect)
}
