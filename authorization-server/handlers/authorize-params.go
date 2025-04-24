package handlers

import (
	"net/http"
	"net/url"
)

type AuthorizeParams struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	Scope        string
	State        string
	Error        string
}

func AuthorizeParamsFromQuery(r *http.Request) *AuthorizeParams {
	return &AuthorizeParams{
		ClientID:     r.URL.Query().Get("client_id"),
		RedirectURI:  r.URL.Query().Get("redirect_uri"),
		ResponseType: r.URL.Query().Get("response_type"),
		Scope:        r.URL.Query().Get("scope"),
		State:        r.URL.Query().Get("state"),
		Error:        r.URL.Query().Get("error"),
	}
}

func AuthorizeParamsFromForm(r *http.Request) *AuthorizeParams {
	return &AuthorizeParams{
		ClientID:     r.FormValue("client_id"),
		RedirectURI:  r.FormValue("redirect_uri"),
		ResponseType: r.FormValue("response_type"),
		Scope:        r.FormValue("scope"),
		State:        r.FormValue("state"),
		Error:        r.FormValue("error"),
	}
}

func (params *AuthorizeParams) toQuery() url.Values {
	query := url.Values{}
	query.Set("client_id", params.ClientID)
	query.Set("redirect_uri", params.RedirectURI)
	query.Set("response_type", params.ResponseType)
	query.Set("scope", params.Scope)
	query.Set("state", params.State)
	query.Set("error", params.Error)
	return query
}
