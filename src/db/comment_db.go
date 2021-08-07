package db

import "github.com/matthewkappus/MagicCard/src/comment"

// contact table
const (
	createComment = `CREATE TABLE IF NOT EXISTS contact(id INTEGER PRIMARY KEY, sender_name STRING, sender_fullname TEXT, sender_email TEXT, student_name TEXT, sent DATETIME, respondent TEXT, starstrike INT, message TEXT)`
	insertContact = `INSERT INTO contact(sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
)

func (s *Store) CreateCommentTable() {
	s.db.Exec(createComment)
}

func (s *Store) InsertContact(c *comment.Contact) error {

	// `INSERT INTO comment(sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	id, err := s.db.Exec(insertContact, c.Sender.Teacher, c.Sender.FullName, c.Sender.StaffEmail, c.StudentName, c.Sent, c.Respondent, c.StarStrike.ID, c.Message)
	if err != nil {
		return err
	}
	ii, err := id.LastInsertId()
	c.ID = int(ii)
	return err
}
