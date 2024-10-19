package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"jukebox/internal/db"
	"jukebox/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

func GetChatMessages(database *db.DB, w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		http.Error(w, "Missing conversation id", http.StatusBadRequest)
		return
	}

	conversationIdPathParam := r.PathValue("id")

	conversationId, err := strconv.ParseInt(conversationIdPathParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := database.GetChatMessages(int(conversationId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateChatMessage(database *db.DB, llm llms.Model, w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		http.Error(w, "Missing conversation id", http.StatusBadRequest)
		return
	}

	conversationIdPathParam := r.PathValue("id")
	ucm := UserChatMessage{}

	err := json.NewDecoder(r.Body).Decode(&ucm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conversationId, err := strconv.ParseInt(conversationIdPathParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := database.GetChatMessages(int(conversationId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages = append(messages, models.ChatMessage{
		Role:           "user",
		ConversationId: int(conversationId),
		Content:        ucm.Message,
		Files:          ucm.Files,
	})

	err = database.CreateChatMessage(int(conversationId), "user", ucm.Message, nil, ucm.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	flusher := w.(http.Flusher)

	systemPrompt := `
	You are a helpful assistant. For the request made to you, please provide your
	response in the following format, where you are providing the contents of the artifact
	you are generating and your thought process in distinctly demarcated tags. Please ensure that
	you do not provide any content outside of these tags.

	If you are asked to generate a document, please think though the various sections of the document and put in as
	much detail as possible. Be thorough and detailed.
	<artifact>
	# Markdown of the document

	## Section 1
	...
	</artifact>
	<explanation>
		Users first need to install the pre-requisites because...
	</explanation>
	`
	llmMessages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
	}
	for _, m := range messages {
		llmRole := llms.ChatMessageTypeHuman
		switch m.Role {
		case "system":
			llmRole = llms.ChatMessageTypeSystem
		case "assistant":
			llmRole = llms.ChatMessageTypeAI
		}

		content := m.Content
		if m.Artifact != nil {
			content = fmt.Sprintf("<artifact>\n%s\n</artifact>\n<explanation>\n%s\n</explanation>", *m.Artifact, m.Content)
		}

		log.Printf("for conversation %d files are %v", conversationId, m.Files)
		if len(m.Files) > 0 {
			for _, f := range m.Files {
				content += fmt.Sprintf(`<context>\n<file_name>%s</file_name>\n<contents>\n%s\n</contents>\n</context>\n`, f.Name, f.Contents)
			}
		}
		llmMessages = append(llmMessages, llms.TextParts(llmRole, content))
	}
	log.Printf("for conversation %d llmMessages: %v", conversationId, llmMessages)

	streamMessage := &UserChatMessageResponse{
		Message:  "",
		Artifact: "",
	}

	inArtifact := false
	inExplanation := false
	collectedChunks := ""

	_, err = llm.GenerateContent(r.Context(), llmMessages, llms.WithMaxTokens(16384), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		collectedChunks += string(chunk)

		if strings.Contains(collectedChunks, "<artifact>") {
			inArtifact = true
			collectedChunks = strings.Replace(collectedChunks, "<artifact>", "", -1)
			return nil
		}

		if strings.Contains(collectedChunks, "</artifact>") {
			inArtifact = false
			collectedChunks = ""
			return nil
		}

		if inArtifact {
			streamMessage.Artifact += string(chunk)
			streamMessage.Artifact = strings.Replace(streamMessage.Artifact, "</artifact", "", -1)
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			flusher.Flush()
		}

		if strings.Contains(collectedChunks, "<explanation>") {
			inExplanation = true
			collectedChunks = strings.Replace(collectedChunks, "<explanation>", "", -1)
			return nil
		}

		if strings.Contains(collectedChunks, "</explanation>") {
			inExplanation = false
			collectedChunks = ""
			return nil
		}

		if inExplanation {
			streamMessage.Message += string(chunk)
			streamMessage.Message = strings.Replace(streamMessage.Message, "</explanation", "", -1)
			err := json.NewEncoder(w).Encode(streamMessage)
			if err != nil {
				return err
			}
			w.Write([]byte("\r\n"))
			flusher.Flush()
		}

		return nil
	}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = database.CreateChatMessage(int(conversationId), "assistant", streamMessage.Message, &streamMessage.Artifact, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UserChatMessage struct {
	Message string        `json:"message"`
	Files   []models.File `json:"files,omitempty"`
}

type UserChatMessageResponse struct {
	Message  string `json:"message"`
	Artifact string `json:"artifact,omitempty"`
}
