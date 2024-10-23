package middleware

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Middle struct {
	store *sessions.CookieStore
}

func NewMiddle(store *sessions.CookieStore) *Middle {
	return &Middle{
		store: store,
	}
}

func (m *Middle) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, "auth")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
