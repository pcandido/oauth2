package handlers

import "net/http"

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token":"example-access-token","refresh_token":"example-refresh-token","token_type":"Bearer","expires_in":3600}`))
}
