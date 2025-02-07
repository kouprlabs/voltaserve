// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/webdav/client/idp_client"
	"github.com/kouprlabs/voltaserve/webdav/config"
	"github.com/kouprlabs/voltaserve/webdav/handler"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

var (
	tokens   = make(map[string]*infra.Token)
	expiries = make(map[string]time.Time)
	mu       sync.RWMutex
)

func startTokenRefresh(idpClient *idp_client.TokenClient) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-ticker.C
			mu.Lock()
			for username, token := range tokens {
				expiry := expiries[username]
				if time.Now().After(expiry.Add(-1 * time.Minute)) {
					newToken, err := idpClient.Exchange(idp_client.TokenExchangeOptions{
						GrantType:    idp_client.GrantTypeRefreshToken,
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

func basicAuthMiddleware(next http.Handler, idpClient *idp_client.TokenClient) http.Handler {
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
			token, err = idpClient.Exchange(idp_client.TokenExchangeOptions{
				GrantType: idp_client.GrantTypePassword,
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
// @version	3.0.0
// @BasePath	/v3
//
// .
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

	idpClient := idp_client.NewTokenClient()

	h := handler.NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/v3/health", h.Health)
	mux.HandleFunc("/version", h.Version)
	mux.HandleFunc("/", h.Dispatch)

	startTokenRefresh(idpClient)

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		ReadHeaderTimeout: 30 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v3/health") || strings.HasPrefix(r.URL.Path, "/version") {
				mux.ServeHTTP(w, r)
			} else {
				basicAuthMiddleware(mux, idpClient).ServeHTTP(w, r)
			}
		}),
	}

	log.Printf("Listening on %s", server.Addr)

	log.Fatal(server.ListenAndServe())
}
