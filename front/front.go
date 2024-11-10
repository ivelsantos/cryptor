package front

import (
	"github.com/ivelsantos/cryptor/front/handlers"
	"github.com/ivelsantos/cryptor/front/middle"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func Front() {
	// Creating a storage to store sessions
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options = &sessions.Options{
		Path: "/",

		// Setting sessions cookies to last only a week. Default is a month
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	handler := handlers.NewHandler(store)
	middle := middleware.NewMiddle(store)
	_ = middle

	m := http.NewServeMux()

	fs := http.FileServer(http.Dir("./assets"))
	m.Handle("GET /assets/", http.StripPrefix("/assets", fs))

	// Giving the handlers the http signature
	handlerHome := http.HandlerFunc(handler.Home)
	handlerEditor := http.HandlerFunc(handler.AlgoEditor)
	handlerEditorSave := http.HandlerFunc(handler.EditorSave)
	handlerAlgoDelete := http.HandlerFunc(handler.AlgoDelete)
	handlerAlgoStateUpdate := http.HandlerFunc(handler.AlgoStateUpdate)

	// Home related
	m.Handle("GET /", middle.CheckAuth(handlerHome))

	//Login related
	m.HandleFunc("GET /user/login", handler.Login)
	m.HandleFunc("POST /user/login/{name}", handler.Login)
	m.HandleFunc("GET /user/login/login", handler.LoginLogin)
	m.HandleFunc("GET /user/login/register", handler.LoginRegister)
	m.HandleFunc("GET /user/logout", handler.Logout)
	m.HandleFunc("POST /user/register", handler.Register)

	// Editor related
	m.Handle("GET /editor", middle.CheckAuth(handlerEditor))
	m.Handle("POST /editor/save", middle.CheckAuth(handlerEditorSave))
	m.Handle("DELETE /editor/delete/{id}", middle.CheckAuth(handlerAlgoDelete))
	m.Handle("POST /editor/update/state/{id}", middle.CheckAuth(handlerAlgoStateUpdate))

	log.Printf("Starting server at port :1234\n\n")
	log.Fatal(http.ListenAndServe(":1234", m))
}
