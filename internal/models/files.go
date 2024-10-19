package models

type File struct {
	ID             int    `json:"id"`
	ConversationId int    `json:"conversation_id,omitempty"`
	MessageId      int    `json:"message_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Contents       string `json:"contents,omitempty"`
}
