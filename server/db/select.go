package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"fmt"
)


func SelectCourse(userid,courseid string) {
	statement, _ := db_obj.Prepare("update selection set course_id=?  where user_id=?")
	statement.Exec(courseid,userid)
}
func SelectProblem(userid,probid string) {
	statement, _ := db_obj.Prepare("update selection set problem_id=?  where user_id=?")
	statement.Exec(probid,userid)
}
func GetSelectedCourse(userid string) string{
	statement, _ := db_obj.Prepare("select course_id from selection  where user_id=?")
	rows,_ :=statement.Query(userid)
	var tempC string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tempC)
		return tempC
	}
	return "Not Selected"
}
func GetSelectedProblem(userid string) string{
	statement, _ := db_obj.Prepare("select problem_id from selection  where user_id=?")
	rows,_ :=statement.Query(userid)
	var tempC string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tempC)
		return tempC
	}
	return "Not Selected"
}/*
func SelectProblem(userid,probid string) {
	statement, _ := db_obj.Prepare("update selection set problem_id=?  where user_id=?")
	statement.Exec(probid,userid)
}*/
func ClearSelection(userid string) {
	statement, _ := db_obj.Prepare("delete from selection WHERE user_id=?")
	statement.Exec(userid)
}

func CreateSelection(userid string) {
	statement, _ := db_obj.Prepare("insert into selection (user_id,course_id,problem_id) VALUES (?, 'Not Selected','Not Selected')")
	statement.Exec(userid)
}

