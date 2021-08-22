package db

import (
	"fmt"

	"github.com/matthewkappus/MagicCard/src/comment"
)

// contact table
const (
	createComment       = `CREATE TABLE IF NOT EXISTS contact(id INTEGER PRIMARY KEY, sender_name STRING, sender_fullname TEXT, sender_email TEXT, student_name TEXT, perm_id TEXT, sent DATETIME, respondent TEXT, starstrike INT, message TEXT)`
	selectContactByPerm = `SELECT * FROM contact WHERE perm_id = ?`
	insertContact       = `INSERT INTO contact(sender_name, sender_fullname, sender_email, student_name, perm_id, sent, respondent, starstrike, message) VALUES(?, ?,?, ?, ?, ?, ?, ?, ?)`
)

func (s *Store) CreateCommentTable() {
	s.db.Exec(createComment)
}

func (s *Store) InsertContact(c *comment.Contact) error {

	// `INSERT INTO comment(sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	id, err := s.db.Exec(insertContact, c.Sender.Teacher, c.Sender.FullName, c.Sender.StaffEmail, c.StudentName, c.PermID, c.Sent, c.Respondent, c.StarStrike.ID, c.Message)
	if err != nil {
		return err
	}
	ii, err := id.LastInsertId()
	c.ID = int(ii)
	return err
}

func (s *Store) GetContacts(perm_id string) ([]comment.Contact, error) {
	rows, err := s.db.Query(selectContactByPerm, perm_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []comment.Contact
	for rows.Next() {
		c := comment.Contact{
			StarStrike: &comment.StarStrike{},
		}

		err := rows.Scan(&c.ID,
			&c.Sender.Teacher,
			&c.Sender.FullName,
			&c.Sender.StaffEmail,
			&c.StudentName,
			&c.PermID,
			&c.Sent,
			&c.Respondent,
			&c.StarStrike.ID,
			&c.Message)
		if err != nil {
			continue
		}
		contacts = append(contacts, c)
	}

	for id, c := range contacts {
		if c.StarStrike.ID != 0 {
			contacts[id].StarStrike, err = s.GetStarStrike(c.StarStrike.ID)
			if err != nil {
				fmt.Println("ss scann err", err.Error())
			}
		} else {
			contacts[id].StarStrike = nil
		}
	}
	return contacts, nil
}
