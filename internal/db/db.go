package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDB(driver string, connectionString string) (*DB, error) {
	database, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	createConversationsTable := `
	CREATE TABLE IF NOT EXISTS conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		goal TEXT,
		artifact TEXT
	);`

	_, err = database.Exec(createConversationsTable)
	if err != nil {
		return nil, err
	}

	createChatMessagesTable := `
	CREATE TABLE IF NOT EXISTS chat_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INT,
		role TEXT,
		content TEXT,
		artifact TEXT,
		CONSTRAINT fk_conversation
      		FOREIGN KEY(conversation_id) 
        	REFERENCES conversations(id)
	);`

	_, err = database.Exec(createChatMessagesTable)
	if err != nil {
		return nil, err
	}

	createFilesTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INT,
		message_id INT,
		name TEXT,
		contents TEXT,
		CONSTRAINT fk_conversation
      		FOREIGN KEY(conversation_id) 
        	REFERENCES conversations(id)
		CONSTRAINT fk_message
      		FOREIGN KEY(message_id) 
        	REFERENCES chat_messages(id)
	);`

	_, err = database.Exec(createFilesTable)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: database,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
