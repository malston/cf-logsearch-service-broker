package handlers

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

func HandleAuthCheck() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authHeader := parseAuthzHeader(req)
		authEnvString := createAuthzStringFromEnv()
		if authHeader != authEnvString {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}

func parseAuthzHeader(req *http.Request) string {
	authzHeader := req.Header.Get("Authorization")

	parts := strings.Split(authzHeader, " ")

	if len(parts) != 2 {
		return ""
	}

	if strings.ToLower(parts[0]) != "basic" {
		return ""
	}

	return "basic " + parts[1]
}

func createAuthzStringFromEnv() string {
	username := os.Getenv("LOGSEARCH_BROKER_USERNAME")
	password := os.Getenv("LOGSEARCH_BROKER_PASSWORD")
	data := []byte(username + ":" + password)
	return "basic " + base64.StdEncoding.EncodeToString(data)
}
