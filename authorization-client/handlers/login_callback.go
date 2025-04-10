package handlers

import (
	"authorization-client/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func LoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	body := fmt.Sprintf(`{"code": "%s"}`, code)
	res, err := http.Post(config.AuthorizationServerUrl("token", true), "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("Error getting token: %s\n", err)
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		http.Error(w, "Error getting token", res.StatusCode)
		return
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		http.Error(w, "Error parsing token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := responseBody["access_token"].(string)
	if !ok {
		http.Error(w, "Access token not found in response", http.StatusInternalServerError)
		return
	}

	refreshToken, ok := responseBody["refresh_token"].(string)
	if !ok {
		http.Error(w, "Refresh token not found in response", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		// For production, set Secure and HttpOnly
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		// For production, set Secure and HttpOnly
	})

	w.Header().Set("Location", config.Url(""))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
