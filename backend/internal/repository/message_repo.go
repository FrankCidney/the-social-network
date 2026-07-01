package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/models"
)

type MessageRepository interface {
	CreateMessage(msg *models.Message) error
	GetPrivateMessages(user1ID, user2ID string, limit, offset int) ([]*models.Message, error)
	GetGroupMessages(groupID string, limit, offset int) ([]*models.Message, error)
}

type sqliteMessageRepo struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &sqliteMessageRepo{db: db}
}

func (r *sqliteMessageRepo) CreateMessage(msg *models.Message) error {
	const query = `
		INSERT INTO messages (id, sender_id, receiver_id, group_id, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, msg.ID, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Content, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func (r *sqliteMessageRepo) GetPrivateMessages(user1ID, user2ID string, limit, offset int) ([]*models.Message, error) {
	const query = `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE group_id IS NULL AND 
		      ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, user1ID, user2ID, user2ID, user1ID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get private messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *sqliteMessageRepo) GetGroupMessages(groupID string, limit, offset int) ([]*models.Message, error) {
	const query = `
		SELECT id, sender_id, group_id, content, created_at
		FROM messages
		WHERE group_id = ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get group messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.GroupID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
