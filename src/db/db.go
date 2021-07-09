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
	createStu415          = `CREATE TABLE IF NOT EXISTS stu415 (organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled text)`
	dropStu415            = `DROP TABLE IF EXISTS stu415`
	insertStu415          = `INSERT INTO stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	selectStu415ByPermID  = `SELECT organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled FROM stu415 WHERE perm_id=? `
	selectStu415BySection = `SELECT organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled FROM stu415 WHERE section_id=? `

	//    SELECT  DISTINCT * from stu415 WHERE teacher="Susco Taylor, Kevin R.";
	selectDistinctSectionsByTeacher = `SELECT DISTINCT section_id,  course_id_and_title from stu415 WHERE teacher=?`
)

// staff table
const (
	// teacher is the s415 full name and name is the Mr/Mrs version. Email is their aps gmail
	createStaff = `CREATE TABLE IF NOT EXISTS staff(teacher NOT NULL UNIQUE, name, staff_email, key NOT NULL UNIQUE, FOREIGN KEY(teacher) REFERENCES stu415(teacher))`
	// key is a uuid for session
	// insertStaff              = `INSERT INTO staff(teacher, name, staff_email, key) VALUES(?,?,?,?)`
	insertStaff              = `INSERT INTO staff(teacher, name, staff_email) VALUES(?,?,?)`
	selectTeacherNameByEmail = `SELECT teacher, name FROM staff where staff_email=?`
)

// TeacherNameFromEmail returns the "teacher" associated with stu415s and their formal name
func (s *Store) TeacherNameFromEmail(email string) (teacher, name string, err error) {
	err = s.db.QueryRow(selectTeacherNameByEmail, strings.ToLower(email)).Scan(&teacher, &name)
	return teacher, name, err
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
		// key := uuid.Must(uuid.NewRandom()).String()
		// fmt.Println(" made key ", key)
		// INSERT INTO staff(teacher, name, staff_email, key) VALUES(?,?,?,?)
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

// comment table
const (
	createComment = `CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY, perm_id, email, comment TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, is_merrit BOOLEAN, is_active BOOLEAN DEFAULT true, FOREIGN KEY(perm_id) REFERENCES stu415(perm_id))`
	insertComment = `INSERT INTO comment(perm_id, email, comment, is_merrit) VALUES(?,?,?,?);`

	// INSERT INTO comment(perm_id, email, comment, created, is_merrit, is_active) VALUES(1, "2", "3", time('now'), true, true);

)

func (s *Store) SelectStu415(permid string) (s415 *synergy.Stu415, err error) {
	s415 = new(synergy.Stu415)
	row := s.db.QueryRow(selectStu415ByPermID, permid)
	// organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled
	err = row.Scan(&s415.OrganizationName, &s415.SchoolYear, &s415.StudentName, &s415.PermID, &s415.Gender, &s415.Grade, &s415.TermName, &s415.Per, &s415.Term, &s415.SectionID, &s415.CourseIDAndTitle, &s415.MeetDays, &s415.Teacher, &s415.Room, &s415.PreScheduled)

	return
}
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

// ListStudents takes a stu415 section_id and returns the students
func (s *Store) ListStudents(section string) ([]*synergy.Stu415, error) {
	rows, err := s.db.Query(selectStu415BySection, section)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := make([]*synergy.Stu415, 0)
	for rows.Next() {
		s415 := new(synergy.Stu415)
		err = rows.Scan(&s415.OrganizationName, &s415.SchoolYear, &s415.StudentName, &s415.PermID, &s415.Gender, &s415.Grade, &s415.TermName, &s415.Per, &s415.Term, &s415.SectionID, &s415.CourseIDAndTitle, &s415.MeetDays, &s415.Teacher, &s415.Room, &s415.PreScheduled)
		if err != nil {
			log.Printf("sql couldn't scan student: %v", err)
			continue
		}
		students = append(students, s415)
	}
	return students, rows.Err()
}

// ListClasses takes a teacher name and returns stu415s with unique sections
func (s *Store) ListClasses(teacher string) ([]*synergy.Stu415, error) {
	rows, err := s.db.Query(selectDistinctSectionsByTeacher, teacher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	classes := make([]*synergy.Stu415, 0)
	for rows.Next() {
		// SELECT DISTINCT section_id,  course_id_and_title from stu415 WHERE teacher=?
		stu := new(synergy.Stu415)
		err = rows.Scan(&stu.SectionID, &stu.CourseIDAndTitle)
		if err != nil {
			log.Printf("error listing classes: %v", err)
			continue
		}
		classes = append(classes, stu)

	}

	return classes, nil
}
