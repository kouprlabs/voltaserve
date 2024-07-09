package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"voltaserve/config"
	"voltaserve/infra"
)

type IdPErrorResponse struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	MoreInfo    string `json:"moreInfo"`
}

type IdPError struct {
	Value IdPErrorResponse
}

func (e *IdPError) Error() string {
	return fmt.Sprintf("IdPError: %v", e.Value)
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

const (
	GrantTypePassword     = "password"
	GrantTypeRefreshToken = "refresh_token"
)

type TokenExchangeOptions struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Locale       string `json:"locale,omitempty"`
}

type IdPClient struct {
	config *config.Config
}

func NewIdPClient() *IdPClient {
	return &IdPClient{
		config: config.GetConfig(),
	}
}

func (cl *IdPClient) Exchange(options TokenExchangeOptions) (*Token, error) {
	form := url.Values{}
	form.Set("grant_type", options.GrantType)
	if options.Username != "" {
		form.Set("username", options.Username)
	}
	if options.Password != "" {
		form.Set("password", options.Password)
	}
	if options.RefreshToken != "" {
		form.Set("refresh_token", options.RefreshToken)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/token", cl.config.IdPURL), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	return cl.jsonResponseOrThrow(resp)
}

func (cl *IdPClient) jsonResponseOrThrow(resp *http.Response) (*Token, error) {
	if resp.Header.Get("Content-Type") == "application/json" {
		var jsonResponse Token
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &jsonResponse)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var idpError IdPErrorResponse
			err = json.Unmarshal(body, &idpError)
			if err != nil {
				return nil, err
			}
			return nil, &IdPError{Value: idpError}
		}
		return &jsonResponse, nil
	} else {
		if resp.StatusCode > 299 {
			return nil, fmt.Errorf(resp.Status)
		}
	}
	return nil, fmt.Errorf("unexpected response format")
}

func (cl *IdPClient) GetHealth() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/health", cl.config.IdPURL))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			infra.GetLogger().Error(err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
