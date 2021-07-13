package middleware

import (
	"fmt"
	"net/http"

	"github.com/shipsaw/scenario/context"
	"github.com/shipsaw/scenario/models"
)

type RequireUser struct {
	models.UserService
}

// Middleware to verify that the user is logged in
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			fmt.Println("BREAK 1")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			fmt.Println("BREAK 2")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		fmt.Println("User found: ", user)

		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
