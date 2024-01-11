package httpServer

import (
	"encoding/json"
	"iotTester/database"
	"iotTester/dto"
	"iotTester/models"
	"log"
	"net/http"
)

func ReceiveDataFromMqtt(topic string, payload []byte) {
	var data dto.DeviceDataDto
	err := json.Unmarshal(payload, &data)

	if err != nil {
		log.Println(err.Error())
		return
	}

	jsonBody, _ := json.Marshal(data.Data)

	database.DB.Create(&models.Log{
		Path:     topic,
		DeviceID: data.DeviceID,
		Data:     jsonBody,
	})

	json, _ := json.Marshal(map[string]interface{}{
		"path": topic,
		"body": data,
		"mqtt": true,
	})

	log.Println(string(json))

	webSocket.Broadcast(json)
}

func receiveData(w http.ResponseWriter, r *http.Request) {
	var body dto.DeviceDataDto
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(body.Data)

	database.DB.Create(&models.Log{
		Path:     r.URL.Path,
		DeviceID: body.DeviceID,
		Data:     jsonBody,
	})

	json, _ := json.Marshal(map[string]interface{}{
		"path": r.URL.Path,
		"body": body,
		"mqtt": false,
	})

	log.Println(string(json))

	webSocket.Broadcast(json)
}

func webSocketReq(w http.ResponseWriter, r *http.Request) {
	webSocket.HandleRequest(w, r)
}

// func middleware(r *http.Request) (string, error) {
// 	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
// 	if len(authHeader) != 2 {
// 		return "", errors.New("not bearer token received")
// 	}
// 	token := authHeader[1]
// 	id, err := GetIdFromJwt(token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }
// func login(w http.ResponseWriter, r *http.Request) {
// 	var data map[string]string
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	token, err := GenerateJwt(data["username"])
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// 	bearerData := map[string]interface{}{
// 		"token_type": "Bearer",
// 		"token":      token,
// 		"expires":    28800,
// 	}
// 	jData, _ := json.Marshal(bearerData)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jData)
// 	json, _ := json.Marshal(map[string]interface{}{
// 		"username": data["username"],
// 		"password": data["password"],
// 		"logged":   false,
// 	})
// 	log.Println(string(json))
// 	m.Broadcast(json)
// }
