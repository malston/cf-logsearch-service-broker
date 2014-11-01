package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
)

func HandleAuthCheck() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authHeader := parseAuthzHeader(req)
		log.Printf("HTTP AUTH HEADER: %s", authHeader)
		authEnvString := createAuthzStringFromEnv()
		log.Printf("ENV AUTH HEADER: %s", authEnvString)
		if authHeader != authEnvString {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}

func parseAuthzHeader(req *http.Request) string {
	authzHeader := req.Header.Get("Authorization")
	log.Printf("FULL HTTP AUTHZ HEADER: %s", authzHeader)

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
	log.Printf("username: %s, password: %s", username, password)
	data := []byte(username + ":" + password)
	return "basic " + base64.StdEncoding.EncodeToString(data)
}
