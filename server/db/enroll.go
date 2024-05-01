package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"fmt"
)


func Enroll(userid,courseid string) {
	statement, _ := db_obj.Prepare("INSERT INTO enroll (user_id, course_id, isadmin) VALUES (?, ?, ?)")
	statement.Exec(userid,courseid,"false")
}
func EnrollAdmin(userid,courseid string){
	statement, _ := db_obj.Prepare("INSERT INTO enroll (user_id, course_id, isadmin) VALUES (?, ?, ?)")
	statement.Exec(userid,courseid,"true")
}


