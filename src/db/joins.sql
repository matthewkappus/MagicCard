SELECT DISTINCT stu415.perm_id, stu415.student_name, starstrike.teacher, starstrike.title, starstrike.comment, starstrike.teacher  FROM stu415 join starstrike WHERE stu415.perm_id = starstrike.perm_id;

SELECT DISTINCT  stu415.student_name, starstrike.teacher, starstrike.title, starstrike.comment, starstrike.teacher  FROM stu415 LEFT JOIN starstrike WHERE stu415.perm_id = starstrike.perm_id;

SELECT * from starstrike left join stu415 on stu415.perm_id = starstrike.perm_id;