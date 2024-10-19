package routes

import (
	"encoding/json"
	"jukebox/internal/db"
	"jukebox/internal/models"
	"net/http"
	"strconv"
)

func GetConversations(database *db.DB, w http.ResponseWriter, r *http.Request) {
	conversations, err := database.GetConversations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetConversation(database *db.DB, w http.ResponseWriter, r *http.Request) {

	conversationId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conversations, err := database.GetConversationById(conversationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateConversation(database *db.DB, w http.ResponseWriter, r *http.Request) {
	var c models.Conversation
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordId, err := database.CreateConversation(c.Goal, c.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ID = int(recordId)
	err = json.NewEncoder(w).Encode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
