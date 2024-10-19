package models

type Conversation struct {
	ID       int     `json:"id"`
	Goal     string  `json:"goal"`
	Artifact *string `json:"artifact,omitempty"`
	Files    []File  `json:"files,omitempty"`
}

type ChatMessage struct {
	ID             int     `json:"id"`
	ConversationId int     `json:"conversation_id"`
	Role           string  `json:"role"`
	Content        string  `json:"content"`
	Artifact       *string `json:"artifact,omitempty"`
	Files          []File  `json:"files,omitempty"`
}
