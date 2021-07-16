CREATE DATABASE card;
USE card;
CREATE TABLE  stu415 (organization_name TEXT, school_year TEXT, student_name TEXT, perm_id TEXT, gender TEXT, grade TEXT, term_name TEXT, per TEXT, term TEXT, section_id TEXT, course_id_and_title TEXT, meet_days TEXT, teacher TEXT, room TEXT, pre_scheduled TEXT)

-- CREATE TABLE  staff(teacher  TEXT, name TEXT, staff_email TEXT, key TEXT);

CREATE TABLE staff(teacher TEXT, full_name TEXT, staff_email TEXT, guid TEXT)
CREATE TABLE  starbar(id INTEGER AUTO_INCREMENT, teacher TEXT, title TEXT, comment TEXT, isStar BOOLEAN, PRIMARY KEY(id))
CREATE TABLE  comment (id INTEGER AUTO_INCREMENT, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, created TIMESTAMP DEFAULT CURRENT_TIMESTAMP, isStar BOOLEAN, isActive BOOLEAN DEFAULT true, PRIMARY KEY(id))