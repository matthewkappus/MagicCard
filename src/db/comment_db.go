package db

import (
	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

// contact table
const (
	createComment       = `CREATE TABLE IF NOT EXISTS contact(id INTEGER PRIMARY KEY, sender_name STRING, sender_email TEXT, respondent TEXT, starstrike_id INT, message TEXT, is_closed BOOL)`
	selectContactByPerm = `SELECT * FROM contact WHERE perm_id = ?`
	insertContact       = `INSERT INTO contact( sender_name, sender_email, respondent, starstrike_id, message, is_closed) VALUES(?, ?,?, ?, ?, ?)`

	selectContactByID         = `SELECT * FROM contact WHERE id = ?`
	selectContactByStarStrike = `SELECT * FROM contact WHERE starstrike = ?`
)

func (s *Store) CreateCommentTable() {
	s.db.Exec(createComment)
}

func (s *Store) InsertContact(c *comment.Contact) error {

	// `id INTEGER PRIMARY KEY, sender_name STRING, sender_email TEXT, respondent TEXT, starstrike_id INT, message TEXT, is_closed BOOL
	id, err := s.db.Exec(insertContact, &c.Sender.Name, &c.Sender.Email, &c.Respondent, &c.StarStrike.ID, &c.Message, &c.IsClosed)
	if err != nil {
		return err
	}
	ii, err := id.LastInsertId()
	c.ID = int(ii)
	return err
}

// GetContact takes a strike id and returns a strike or error
func (s *Store) GetContact(strikeID int) (*comment.Contact, error) {
	c := &comment.Contact{
		Sender:     new(synergy.Staff),
		StarStrike: new(comment.StarStrike),
	}
	// (id INTEGER PRIMARY KEY, sender_name STRING, sender_email TEXT, respondent TEXT, starstrike_id INT, message TEXT, is_closed BOOL)
	err := s.db.QueryRow(selectContactByID, strikeID).Scan(&c.ID, &c.Sender.Name, &c.Sender.Email, &c.Respondent, &c.StarStrike.ID, &c.Message, &c.IsClosed)
	return c, err
}
