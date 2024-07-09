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
	"github.com/kouprlabs/voltaserve/webdav/client"
	"github.com/kouprlabs/voltaserve/webdav/handler"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kouprlabs/voltaserve/webdav/config"
)

var (
	tokens   = make(map[string]*infra.Token)
	expiries = make(map[string]time.Time)
	mu       sync.Mutex
)

func startTokenRefresh(idpClient *client.IdPClient) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			mu.Lock()
			for username, token := range tokens {
				expiry := expiries[username]
				if time.Now().After(expiry.Add(-1 * time.Minute)) {
					newToken, err := idpClient.Exchange(client.TokenExchangeOptions{
						GrantType:    client.GrantTypeRefreshToken,
						RefreshToken: token.RefreshToken,
					})
					if err == nil {
						tokens[username] = newToken
						expiries[username] = helper.NewExpiry(newToken)
					}
				}
			}
			mu.Unlock()
		}
	}()
}

func basicAuthMiddleware(next http.Handler, idpClient *client.IdPClient) http.Handler {
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
			var err error
			token, err = idpClient.Exchange(client.TokenExchangeOptions{
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
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "token", token)))
	})
}

// @title		Voltaserve WebDAV
// @version	2.0.0
// @BasePath	/v2
func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		if err := godotenv.Load(".env.local"); err != nil {
			panic(err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}

	cfg := config.GetConfig()

	idpClient := client.NewIdPClient()

	h := handler.NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/health", h.Health)
	mux.HandleFunc("/", h.Dispatch)

	startTokenRefresh(idpClient)

	log.Printf("Listening on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v2/health") {
			mux.ServeHTTP(w, r)
		} else {
			basicAuthMiddleware(mux, idpClient).ServeHTTP(w, r)
		}
	})))
}
