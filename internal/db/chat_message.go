package db

import "jukebox/internal/models"

func (d *DB) CreateChatMessage(conversationId int, role string, content string, artifact *string, files []models.File) error {
	res, err := d.db.Exec("INSERT INTO chat_messages (conversation_id, role, content, artifact) VALUES (?, ?, ?, ?)", conversationId, role, content, artifact)
	if err != nil {
		return err
	}

	msgId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, f := range files {
		_, err := d.db.Exec("INSERT INTO files (name, contents, message_id) VALUES (?, ?, ?)", f.Name, f.Contents, msgId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) GetChatMessages(conversationId int) ([]models.ChatMessage, error) {
	messages := []models.ChatMessage{}
	rows, err := d.db.Query("SELECT id, role, content, artifact FROM chat_messages WHERE conversation_id = ?", conversationId)
	if err != nil {
		return messages, err
	}

	for rows.Next() {
		m := models.ChatMessage{}
		err = rows.Scan(&m.ID, &m.Role, &m.Content, &m.Artifact)
		if err != nil {
			return messages, err
		}

		fileRows, err := d.db.Query("SELECT id, name, contents, message_id FROM files WHERE message_id = ?", m.ID)
		if err != nil {
			return messages, err
		}

		for fileRows.Next() {
			f := models.File{}
			err = fileRows.Scan(&f.ID, &f.Name, &f.Contents, &f.MessageId)
			if err != nil {
				return nil, err
			}

			m.Files = append(m.Files, f)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
