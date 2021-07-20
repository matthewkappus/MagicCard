package db

import (
	"fmt"
	"log"
	"os"

	"github.com/matthewkappus/Roster/src/synergy"
)

// stu415 table
const (
	createStu415          = `CREATE TABLE IF NOT EXISTS stu415 (organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled text)`
	dropStu415            = `DROP TABLE IF EXISTS stu415`
	insertStu415          = `INSERT INTO stu415(organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	selectStudentList     = `SELECT DISTINCT perm_id, student_name FROM stu415;`
	selectStu415ByPermID  = `SELECT organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled FROM stu415 WHERE perm_id=? `
	selectStu415BySection = `SELECT organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled FROM stu415 WHERE section_id=? `

	//    SELECT  DISTINCT * from stu415 WHERE teacher="Susco Taylor, Kevin R.";
	selectDistinctSectionsByTeacher = `SELECT DISTINCT section_id,  course_id_and_title from stu415 WHERE teacher=?`
)

func (s *Store) SelectStu415(permid string) (s415 *synergy.Stu415, err error) {
	s415 = new(synergy.Stu415)
	row := s.db.QueryRow(selectStu415ByPermID, permid)
	// organization_name, school_year, student_name, perm_id, gender, grade, term_name, per, term, section_id, course_id_and_title, meet_days, teacher, room, pre_scheduled
	err = row.Scan(&s415.OrganizationName, &s415.SchoolYear, &s415.StudentName, &s415.PermID, &s415.Gender, &s415.Grade, &s415.TermName, &s415.Per, &s415.Term, &s415.SectionID, &s415.CourseIDAndTitle, &s415.MeetDays, &s415.Teacher, &s415.Room, &s415.PreScheduled)

	return
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

// ListStudents returns a list of distinct stu415s
func (s *Store) ListAllStudents() ([]*synergy.Stu415, error) {
	rows, err := s.db.Query(selectStudentList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := make([]*synergy.Stu415, 0)
	for rows.Next() {
		s415 := new(synergy.Stu415)
		// SELECT DISTINCT perm_id, student_name FROM stu415
		err = rows.Scan(&s415.PermID, &s415.StudentName)
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
