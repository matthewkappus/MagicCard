-- https://console.cloud.google.com/sql/instances/magiccard/overview?folder&organizationId&project=mymagiccard&cloudshell=true
gcloud sql connect magiccard --user=root


CREATE DATABASE card;
USE card;
CREATE TABLE  stu415 (organization_name TEXT, school_year TEXT, student_name TEXT, perm_id TEXT, gender TEXT, grade TEXT, term_name TEXT, per TEXT, term TEXT, section_id TEXT, course_id_and_title TEXT, meet_days TEXT, teacher TEXT, room TEXT, pre_scheduled TEXT)

-- CREATE TABLE  staff(teacher  TEXT, name TEXT, staff_email TEXT, key TEXT);

CREATE TABLE staff(teacher TEXT, full_name TEXT, staff_email TEXT, guid TEXT)
CREATE TABLE starstrike (id INTEGER PRIMARY KEY, perm_id TEXT, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)
CREATE TABLE IF NOT EXISTS mystarstrike (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)