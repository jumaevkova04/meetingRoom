package models

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response ...
type Response struct {
	Code    int    `json:"-"`
	Message string `json:"-"`

	Payload interface{} `json:"payload"`
}

// Send - метод отправки ответа
func (res *Response) Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(res.Code)

	if res.Payload == nil && res.Code != http.StatusOK {
		res.Payload = struct {
			Error   bool   `json:"error,omitempty"`
			Message string `json:"message,omitempty"`
		}{
			Error:   true,
			Message: res.Message,
		}
	}
	if len(res.Message) == 0 {
		res.Message = http.StatusText(res.Code)
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println("Sending response failed:", err)
	}
}
