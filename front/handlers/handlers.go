package handlers

import (
	"github.com/ivelsantos/cryptor/front/views"
	"github.com/ivelsantos/cryptor/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

type Handler struct {
	store *sessions.CookieStore
}

func NewHandler(store *sessions.CookieStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth")
	name := session.Values["user"].(string)
	algos, err := models.GetAlgos(name)
	if err != nil {
		log.Fatal(err)
	}

	sets := views.HomeSettings{User: name, Algos: algos}

	views.Base(views.Home(sets)).Render(r.Context(), w)
}

func (h *Handler) AlgoEditor(w http.ResponseWriter, r *http.Request) {
	views.Algo().Render(r.Context(), w)
}

func (h *Handler) AlgoDelete(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth")
	owner := session.Values["user"].(string)
	id, ok := strconv.Atoi(r.PathValue("id"))
	if ok != nil {
		log.Print(ok)
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	err := models.DeleteAlgo(id, owner)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditorSave(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "auth")
	owner, _ := session.Values["user"].(string)
	name := r.FormValue("namecode")
	created := int(time.Now().Unix())
	buycode := r.FormValue("buycode")
	state := "new"
	baseAsset := r.FormValue("baseAsset")
	quoteAsset := r.FormValue("quoteAsset")

	algo := models.Algor{Owner: owner, Name: name, Created: created, Buycode: buycode, State: state, BaseAsset: baseAsset, QuoteAsset: quoteAsset}

	err := models.InsertAlgo(algo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
