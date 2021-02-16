package handlers

import (
	"encoding/json"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jumaevkova04/meetingRoom/models"
)

// Server ...
type Server struct {
	mux      *mux.Router
	usersSvc *models.Users
}

// NewServer ...
func NewServer(mux *mux.Router, usersSvc *models.Users) *Server {
	return &Server{mux: mux, usersSvc: usersSvc}
}

// ServeHTTP ...
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	// GET ...
	GET = "GET"
	// PUT ...
	PUT = "PUT"
	// POST ...
	POST = "POST"
	// DELETE ...
	DELETE = "DELETE"
)

// Init ...
func (s *Server) Init() {
	s.mux.HandleFunc("/register", s.RegisterUsers).Methods(POST)
}

// RegisterUsers ...
func (s *Server) RegisterUsers(w http.ResponseWriter, r *http.Request) {

	var (
		user     models.User
		response = models.Response{
			Code: http.StatusOK,
		}
		err error
	)
	defer response.Send(w, r)

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Code = http.StatusBadRequest
		log.Println("handlers: RegisterUsers parsing request failed:", err)
		return
	}
	s.usersSvc.Register(r.Context(), &user)
	response.Payload = user
}
