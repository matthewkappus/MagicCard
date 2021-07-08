package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/matthewkappus/Roster/src/synergy"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// todo: check for tables
	return &Store{db}, nil
}

// stu415 table
const (
	createStu415         = `CREATE TABLE IF NOT EXISTS stu415 (organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled text)`
	dropStu415           = `DROP TABLE IF EXISTS stu415`
	insertStu415         = `INSERT INTO stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	selectStu415ByPermID = `SELECT organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled FROM stu415 WHERE perm_id=? `
)

// comment table
const (
	createComment = `CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY, perm_id, email, comment TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, is_merrit BOOLEAN, is_active BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	// createComment = `CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY, perm_id, email, comment TEXT, created DATETIME, is_merrit, is_active boolean, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertComment = `INSERT INTO comment(perm_id, email, comment, is_merrit) VALUES(?,?,?,?);`

	// INSERT INTO comment(perm_id, email, comment, created, is_merrit, is_active) VALUES(1, "2", "3", time('now'), true, true);

)

func (s *Store) CreateCommentTable() error {
	_, err := s.db.Exec(createComment)
	return err
}
func (s *Store) InsertComment(permID, email, comment string, isMerrit bool) error {
	_, err := s.db.Exec(insertComment, permID, email, comment, isMerrit)
	return err
}

func (s *Store) UpdateStu415(stu415CSV string) error {
	f, err := os.Open(stu415CSV)
	if err != nil {
		return err
	}
	defer f.Close()
	s415s, err := synergy.ReadStu415sFromCSV(f)
	if err != nil {
		return err
	}
	if len(s415s) < 2 {
		return fmt.Errorf("%s CSV empty (%d students)", stu415CSV, len(s415s))
	}

	// drop old stu415s
	_, err = s.db.Exec(dropStu415)
	if err != nil {
		return err
	}
	// create new one
	_, err = s.db.Exec(createStu415)
	if err != nil {
		return err
	}

	return s.insertStu415(s415s)
}

func (s *Store) insertStu415(stu415s []*synergy.Stu415) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(insertStu415)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, stu415 := range stu415s {
		_, err = stmt.Exec(stu415.OrganizationName, stu415.SchoolYear, stu415.StudentName, stu415.PermID, stu415.Gender, stu415.Grade, stu415.TermName, stu415.Per, stu415.Term, stu415.SectionID, stu415.CourseIDAndTitle, stu415.MeetDays, stu415.Teacher, stu415.Room, stu415.PreScheduled)
		if err != nil {
			log.Printf("insertStu415 sql error: %v", err)
		}
	}
	return tx.Commit()
}

// staff table
const (
	createStaff = `CREATE TABLE IF NOT EXISTS staff(email primary_key, teacher_name text)`
)

// comments table
