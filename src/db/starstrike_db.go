package db

import (
	"fmt"

	"github.com/matthewkappus/MagicCard/src/comment"
)

// starstrike table
const (
	createStarStrike                 = `CREATE TABLE IF NOT EXISTS starstrike (id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertStarStrike                 = `INSERT INTO starstrike(perm_id, teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,?,?)`
	selectStarStrikeByPermID         = `SELECT * FROM starstrike WHERE perm_id = ?`
	selectNewestStarStrikesByTeacher = `SELECT * FROM starstrike WHERE teacher=? LIMIT ?`
	selectStarStrikeByID             = `SELECT * FROM starstrike WHERE id = ?`
	selectStarStrikeByPerm           = `SELECT * FROM starstrike WHERE perm_id = ?`
)

// mystarstrike table holds teacher's saved starstrikes
const (
	createMystarStrike = `CREATE TABLE IF NOT EXISTS mystarstrike (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)`
	// insertMystarStrike = `INSERT INTO mystarstrike(teacher, comment, title, icon, cat, isActive) VALUES("Kappus, Matthew D.", "The first sample strike", "Strike 1", "bi-clock", 1, true);`
	// selecting teacher="all" gets schoolwide starstrikes
	selectMyStarStrikes = `SELECT * FROM mystarstrike WHERE teacher="all" OR teacher = ?`
	insertMyStarStrike  = `INSERT INTO mystarstrike(teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,true)`
)

func (s *Store) CreateMyStarStrikeTable() error {
	_, err := s.db.Exec(createMystarStrike)
	return err
}

// SelectTeacherStarStrikes with teacher name up to query limmit
func (s *Store) SelectTeacherStarStrikes(teacher string, limit int) (stars, strikes []*comment.StarStrike, err error) {
	rows, err := s.db.Query(selectNewestStarStrikesByTeacher, teacher, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	stars = make([]*comment.StarStrike, 0)
	strikes = make([]*comment.StarStrike, 0)

	for rows.Next() {
		strstr := new(comment.StarStrike)
		// id  perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive
		if err := rows.Scan(&strstr.ID, &strstr.PermID, &strstr.Teacher, &strstr.Comment, &strstr.Title, &strstr.Icon, &strstr.Created, &strstr.Cat, &strstr.IsActive); err != nil {
			fmt.Printf("ss scan err: %v\n", err)
			continue
		}
		if strstr.Cat == 0 {
			stars = append(stars, strstr)
			continue
		}
		strikes = append(strikes, strstr)
	}

	return stars, strikes, nil

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

// AddStarStrike into the store
// Should call IsValid() before calling this
func (s *Store) AddStarStrike(ss *comment.StarStrike) error {
	// INSERT INTO starstrike(perm_id, teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,?,?)
	_, err := s.db.Exec(insertStarStrike, ss.PermID, ss.Teacher, ss.Comment, ss.Title, ss.Icon, ss.Cat, true)
	return err
}

// BatchAddStarStrikes into the store
func (s *Store) BatchAddStarStrikes(ss []*comment.StarStrike) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, starstrike := range ss {
		if starstrike == nil {
			continue
		}
		// INSERT INTO starstrike(perm_id, teacher, comment, title, icon, cat, isActive) VALUES(?,?,?,?,?,?,?)
		_, err := tx.Exec(insertStarStrike, starstrike.PermID, starstrike.Teacher, starstrike.Comment, starstrike.Title, starstrike.Icon, starstrike.Cat, true)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()

}

// GetStarStrikesByPerm returns all starstrikes for a given perm_id
func (s *Store) GetStarStrikesByPerm(id string) ([]*comment.StarStrike, error) {
	stu415, err := s.GetStu415ByPermID(id)

	if err != nil {
		return nil, err
	}
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
		strstr.Student = stu415
		ss = append(ss, strstr)
	}

	return ss, nil

}

func (s *Store) GetStarStrike(id int) (*comment.StarStrike, error) {

	strstr := new(comment.StarStrike)

	err := s.db.QueryRow(selectStarStrikeByID, id).Scan(&strstr.ID, &strstr.PermID, &strstr.Teacher, &strstr.Comment, &strstr.Title, &strstr.Icon, &strstr.Created, &strstr.Cat, &strstr.IsActive)
	if err != nil {
		return nil, err
	}
	strstr.Student, err = s.GetStu415ByPermID(strstr.PermID)

	return strstr, err

}

// GetStudentStrikes returns all starstrikes for a given student
func (s *Store) GetStudentStrikes(perm string) ([]*comment.StarStrike, error) {
	stu415, err := s.GetStu415ByPermID(perm)

	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(selectStarStrikeByPerm, perm)
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
		strstr.Student = stu415
		ss = append(ss, strstr)
	}
	return ss, nil

}
