package db

import (
	"fmt"

	"github.com/matthewkappus/MagicCard/src/comment"
)

// starstrike table
const (
	createStarStrike                 = `CREATE TABLE IF NOT EXISTS starstrike (id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertStarStrike                 = `INSERT INTO starstrike(perm_id, teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,?,?);`
	selectStarStrikeByPermID         = `SELECT * FROM starstrike WHERE perm_id = ?`
	selectNewestStarStrikesByTeacher = `SELECT * FROM starstrike WHERE teacher=? LIMIT ?`
	// select count(cat) from starstrike where cat=1 and perm_id="980016917"
)

// mystarstrike table holds teacher's saved starstrikes
const (
	createMystarStrike = `CREATE TABLE IF NOT EXISTS mystarstrike (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)`
	// insertMystarStrike = `INSERT INTO mystarstrike(teacher, comment, title, icon, cat, isActive) VALUES("Kappus, Matthew D.", "The first sample strike", "Strike 1", "bi-clock", 1, true);`
	// selecting teacher="all" gets schoolwide starstrikes
	selectMyStarStrikes = `SELECT * FROM mystarstrike WHERE teacher="all" OR teacher = ?`
	insertMyStarStrike  = `INSERT INTO mystarstrike(teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,true)`
)

// InsertAllMyStarStrikes adds starstrikes for teacher="all" for school-wide use
func (s *Store) InsertAllMyStarStrikes() {

}

// CreateMystarStrike for teacher-create star strikes
func (s *Store) InsertMyStarStrike(teacher, comment, title, icon string, cat int) error {
	_, err := s.db.Exec(insertMyStarStrike, teacher, comment, title, icon, cat)
	return err
}

func (s *Store) GetMyStarStrikes(teacher string) ([]*comment.StarStrike, error) {
	rows, err := s.db.Query(selectMyStarStrikes, teacher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := make([]*comment.StarStrike, 0)
	for rows.Next() {
		strstr := new(comment.StarStrike)
		// id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true
		if err := rows.Scan(&strstr.ID, &strstr.Teacher, &strstr.Comment, &strstr.Title, &strstr.Icon, &strstr.Created, &strstr.Cat, &strstr.IsActive); err != nil {
			fmt.Printf("ss scan err: %v\n", err)
			continue
		}
		ss = append(ss, strstr)
	}
	return ss, nil
}

func (s *Store) GetTeacherStarStrikes(teacher string) (stars, strikes []*comment.StarStrike, err error) {
	rows, err := s.db.Query(selectMyStarStrikes, teacher)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stars = make([]*comment.StarStrike, 0)
	strikes = make([]*comment.StarStrike, 0)
	for rows.Next() {
		sb := new(comment.StarStrike)

		if err = rows.Scan(&sb.ID, &sb.Teacher, &sb.Comment, &sb.Title, &sb.Icon, &sb.Created, &sb.Cat, &sb.IsActive); err != nil {
			continue
		}
		if sb.Cat == comment.Star {
			stars = append(stars, sb)

		} else {

			strikes = append(strikes, sb)
		}
	}
	return stars, strikes, nil
}

// AddStarStrike into the store
func (s *Store) AddStarStrike(perm_id, teacher, comment, title, icon, cat string) error {
	_, err := s.db.Exec(insertStarStrike, perm_id, teacher, comment, title, icon, cat, true)
	return err
}

// GetStarStrikesByPerm returns all starstrikes for a given perm_id
func (s *Store) GetStarStrikesByPerm(id string) ([]*comment.StarStrike, error) {
	rows, err := s.db.Query(selectStarStrikeByPermID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ss := make([]*comment.StarStrike, 0)
	for rows.Next() {
		strstr := new(comment.StarStrike)
		// id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true
		if err := rows.Scan(&strstr.ID, &strstr.PermID, &strstr.Teacher, &strstr.Comment, &strstr.Title, &strstr.Icon, &strstr.Created, &strstr.Cat, &strstr.IsActive); err != nil {
			fmt.Printf("ss scan err: %v\n", err)
			continue
		}
		ss = append(ss, strstr)
	}
	return ss, nil

}
