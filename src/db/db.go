package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/matthewkappus/Roster/src/synergy"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

// OpenStore creates store from provided sqllite db path
func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// todo: check for tables
	return &Store{db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()

}

// staff table
const (
	// teacher is the s415 full name and name is the Mr/Mrs version. Email is their aps gmail
	createStaff              = `CREATE TABLE IF NOT EXISTS staff(teacher NOT NULL UNIQUE, full_name, staff_email TEXT)`
	insertStaff              = `INSERT INTO staff(teacher, full_name, staff_email) VALUES(?,?,?)`
	selectTeacherNameByEmail = `SELECT teacher, full_name FROM staff WHERE staff_email=?`
)

func (s *Store) SelectStu415s() ([]*synergy.Stu415, error) {
	rows, err := s.db.Query(selectStudentList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := make([]*synergy.Stu415, 0)
	for rows.Next() {
		s415 := new(synergy.Stu415)
		// 	selectStudentList     = `SELECT DISTINCT perm_id, student_name FROM stu415;`
		err = rows.Scan(&s415.PermID, &s415.StudentName)
		if err != nil {
			log.Printf("sql couldn't scan student: %v", err)
			continue
		}
		students = append(students, s415)
	}
	return students, rows.Err()
}

// TeacherNameFromEmail returns the "teacher" associated with stu415s and their formal name
func (s *Store) TeacherNameFromEmail(email string) (teacher, name string, err error) {
	err = s.db.QueryRow(selectTeacherNameByEmail, strings.ToLower(email)).Scan(&teacher, &name)
	return teacher, name, err
}

func (s *Store) StudentFromEmail(email string) (*synergy.Stu415, error) {
	perm := email[:9]
	return s.SelectStu415(perm)

}

func (s *Store) CreateStaff(stu415CSV string) error {
	f, err := os.Open(stu415CSV)
	if err != nil {
		return err
	}
	s415s, err := synergy.ReadStu415sFromCSV(f)
	if err != nil {
		return err
	}
	if len(s415s) < 2 {
		return fmt.Errorf("%s CSV empty (%d students)", stu415CSV, len(s415s))
	}

	// create table if not exists
	_, err = s.db.Exec(createStaff)
	if err != nil {
		return err
	}

	// insert staff names and generic email (first.last@aps.edu)
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStaff)
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, stu := range s415s {
		name, email := toNameEmail(stu.Teacher)
		// INSERT INTO staff(teacher, name, staff_email) VALUES(?,?,?)
		_, err = stmt.Exec(stu.Teacher, name, email)
		if err != nil {
			log.Printf("error inserting %s: %v", stu.Teacher, err)
		}
	}

	return tx.Commit()

}

// toNameEmail takes L, F MI. teacher name and returns [fl, f.l@aps.com]
func toNameEmail(teacherName string) (name, email string) {
	fl := strings.Split(teacherName, " ")

	if len(fl) < 2 {
		log.Printf("%s does not have first and last name", teacherName)
		return
	}
	// remove , after first name
	fl[0] = strings.TrimSuffix(fl[0], ",")

	name = fmt.Sprintf("%s %s", fl[1], fl[0])

	email = strings.ToLower(fmt.Sprintf("%s.%s@aps.edu", fl[1], fl[0]))

	return name, email
}
