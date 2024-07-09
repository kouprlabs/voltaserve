// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"voltaserve/client"
	"voltaserve/handler"

	"github.com/joho/godotenv"
	"voltaserve/config"
)

var (
	tokens   = make(map[string]*client.Token)
	expiries = make(map[string]time.Time)
	api      = &client.IdPClient{}
	mu       sync.Mutex
)

func startTokenRefresh() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			mu.Lock()
			for username, token := range tokens {
				expiry := expiries[username]
				if time.Now().After(expiry.Add(-1 * time.Minute)) {
					newToken, err := api.Exchange(client.TokenExchangeOptions{
						GrantType:    client.GrantTypeRefreshToken,
						RefreshToken: token.RefreshToken,
					})
					if err == nil {
						tokens[username] = newToken
						expiries[username] = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
					}
				}
			}
			mu.Unlock()
		}
	}()
}

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV Server"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		token, exists := tokens[username]
		if !exists {
			token, err := api.Exchange(client.TokenExchangeOptions{
				GrantType: client.GrantTypePassword,
				Username:  username,
				Password:  password,
			})
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokens[username] = token
			expiries[username] = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// @title		Voltaserve WebDAV
// @version	2.0.0
// @BasePath	/v2
func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		err := godotenv.Load(".env.local")
		if err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	cfg := config.GetConfig()

	h := handler.NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/health", h.Health)
	mux.HandleFunc("/", h.Dispatch)

	startTokenRefresh()

	log.Printf("Listening on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v2/health") {
			mux.ServeHTTP(w, r)
		} else {
			basicAuthMiddleware(mux).ServeHTTP(w, r)
		}
	})))
}
