package handle_graphql_gqlgen

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

var userCtxKey = &contextKey{"userClaims"}

type contextKey struct {
	name string
}

// GetAuthID is middleware that gets authenticated claim if present and valid.  This doesn't enforce authentication.
func GetAuthID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("no Authorization present")
			next.ServeHTTP(w, r)
			return
		}
		log.Printf("auth token is %+v\n", token)
		tokenSlice := strings.Split(token, " ")
		var bearerToken string
		if len(tokenSlice) > 1 {
			bearerToken = tokenSlice[len(tokenSlice)-1]
		}

		var keys jwt.KeyRegister
		keyCount, err := keys.LoadPEM([]byte(os.Getenv("AUTH0_CLIENT_PUBLIC_KEY")), nil)
		if err != nil {
			log.Fatal("JWT key import: ", err)
		}
		log.Print(keyCount, " JWT key(s) ready")
		claims, err := keys.Check([]byte(bearerToken))
		if err != nil {
			log.Print("credentials denied")
			// if credential bad, just proceed without setting authentication.
			// later on we check and only allow and set context on a valid token.
			next.ServeHTTP(w, r)
			return
		}
		if !claims.Valid(time.Now()) {
			log.Print("time constraints exceeded")
			http.Error(w, "Expired Token", http.StatusForbidden)
			return
		}
		log.Print("hello ", claims.Audiences)
		log.Printf("Claims %+v", claims)

		_, ok := claims.String("email")
		if !ok {
			log.Print("email not present")
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		// put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.

	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *jwt.Claims {
	raw, _ := ctx.Value(userCtxKey).(*jwt.Claims)
	return raw
}
