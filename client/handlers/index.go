package handlers

import (
	"authorization-client/config"
	"authorization-client/store"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Event struct {
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var html strings.Builder

	html.WriteString("<h1>Client</h1>")

	tokens, err := getTokens(r)
	if err != nil {
		html.WriteString(`<p>Welcome, to access the Resource Server, click <a href="/auth_server/initiate">here</a> to authorize.</p>`)
	} else {
		html.WriteString(`<p>Welcome, you are authorized to access the Resource Server.</p>`)
		html.WriteString(getEventTable(tokens))
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}

func getEventTable(tokens *store.Tokens) string {
	resourceUrl := url.URL{
		Scheme: "http",
		Host:   config.ResourceServerHost(),
		Path:   "/resource",
	}

	req, err := http.NewRequest("GET", resourceUrl.String(), nil)
	if err != nil {
		log.Printf("error: failed to create request, error: %v", err)
		return "<p>Failed to access the resource server.</p>"
	}

	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error: failed to get resource, error: %v", err)
		return "<p>Failed to access the resource server.</p>"
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("error: resource server returned status code %d", res.StatusCode)
		return "<p>Failed to access the resource server.</p>"
	}

	defer res.Body.Close()

	var events []Event
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		log.Printf("error: failed to decode response, error: %v", err)
		return "<p>Failed to decode response from resource server.</p>"
	}

	var html strings.Builder
	html.WriteString("<table border='1'>")
	html.WriteString("<tr><th>Title</th><th>Start</th><th>End</th></tr>")
	for _, event := range events {
		html.WriteString("<tr>")
		html.WriteString("<td>" + event.Title + "</td>")
		html.WriteString("<td>" + event.Start + "</td>")
		html.WriteString("<td>" + event.End + "</td>")
		html.WriteString("</tr>")
	}
	html.WriteString("</table>")

	return html.String()
}

func getTokens(r *http.Request) (*store.Tokens, error) {
	session, err := r.Cookie(config.SESSION_COOKIE_NAME)
	if err != nil || session == nil || session.Value == "" {
		log.Printf("error: session cookie not found, error: %v", err)
		return nil, errors.New("session cookie not found")
	}

	tokens, err := store.GetToken(session.Value)
	if err != nil {
		log.Printf("error: failed to get token from store, error: %v", err)
		return nil, errors.New("failed to get token from store")
	}

	return tokens, nil
}
