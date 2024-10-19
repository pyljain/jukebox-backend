package db

import (
	"jukebox/internal/models"
)

func (d *DB) GetConversations() ([]*models.Conversation, error) {
	rows, err := d.db.Query("SELECT id, goal, artifact FROM conversations")
	if err != nil {
		return nil, err
	}

	var cs []*models.Conversation

	for rows.Next() {
		c := &models.Conversation{}
		err = rows.Scan(&c.ID, &c.Goal, &c.Artifact)
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
	}
	return cs, nil
}

func (d *DB) CreateConversation(goal string, files []models.File) (int64, error) {

	result, err := d.db.Exec("INSERT INTO conversations (goal) VALUES (?)", goal)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	for _, f := range files {
		_, err := d.db.Exec("INSERT INTO files (conversation_id, name, contents) VALUES (?, ?, ?)", id, f.Name, f.Contents)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (d *DB) GetConversationById(id int64) (*models.Conversation, error) {

	c := models.Conversation{}
	err := d.db.QueryRow("SELECT id, goal, artifact FROM conversations WHERE id = ?", id).Scan(&c.ID, &c.Goal, &c.Artifact)
	if err != nil {
		return nil, err
	}

	fileRows, err := d.db.Query("SELECT id, name, contents, conversation_id FROM files WHERE conversation_id = ?", id)
	if err != nil {
		return nil, err
	}

	for fileRows.Next() {
		f := models.File{}
		err = fileRows.Scan(&f.ID, &f.Name, &f.Contents, &f.ConversationId)
		if err != nil {
			return nil, err
		}

		c.Files = append(c.Files, f)
	}

	return &c, nil
}
