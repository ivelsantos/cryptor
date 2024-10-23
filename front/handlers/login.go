package handlers

import (
	"github.com/ivelsantos/cryptor/front/views"
	"github.com/ivelsantos/cryptor/models"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		accounts, err := models.GetAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		set := views.LoginSettings{Accounts: accounts}
		views.Base(views.User(set)).Render(r.Context(), w)
		return
	}
	name := r.PathValue("name")

	session, _ := h.store.Get(r, "auth")
	session.Values["authenticated"] = true
	session.Values["user"] = name
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth")

	session.Values["authenticated"] = false
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	apikey := r.FormValue("apikey")
	secretkey := r.FormValue("secretkey")

	account := models.Account{Name: name, ApiKey: apikey, SecretKey: secretkey}

	err := models.InsertAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func (h *Handler) LoginLogin(w http.ResponseWriter, r *http.Request) {
	accounts, err := models.GetAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	set := views.LoginSettings{Accounts: accounts}
	views.Login(set).Render(r.Context(), w)
}

func (h *Handler) LoginRegister(w http.ResponseWriter, r *http.Request) {
	views.Register().Render(r.Context(), w)
}
