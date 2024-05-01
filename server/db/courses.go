package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"encoding/json"
	//"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Course struct {
	id     string `json:"id"`
	name   string `json:"name"`
	year string `json:"year"`
}

func AddCourse(id,name,year string){
	statement, _ := db_obj.Prepare("INSERT INTO courses (id, name, year) VALUES (?, ?, ?)")
	statement.Exec(id,name,year)
}
func IsCourse(userid,courseid string) bool{
	statement, _ := db_obj.Prepare("select id,name,year from enroll join courses on enroll.course_id=courses.id where user_id=? and course_id=?")
	rows,_ :=statement.Query(userid,courseid)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}
func GetCourses(userid string) string{
	statement, _ := db_obj.Prepare("select id,name,year from enroll join courses on enroll.course_id=courses.id where user_id=?")
	rows,_ :=statement.Query(userid)
	var courses []Course
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	for rows.Next() {
		var tempC Course
		rows.Scan(&tempC.id, &tempC.name, &tempC.year)
		//fmt.Printf("%s\n",tempC.id)
		courses=append(courses,tempC)
	}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "ID", "Name", "Year"})
	for i := 0; i < len(courses); i++ {
		t.AppendRow([]interface{}{i+1,courses[i].id, courses[i].name, courses[i].year})
		t.AppendSeparator()
	}
	return t.Render()
}


func IsAdmin(userid string) bool{
	statement, _ := db_obj.Prepare("select * from enroll join courses on enroll.course_id=courses.id join selection on selection.user_id=enroll.user_id where enroll.user_id=? and isadmin ='true'")
	rows,_ :=statement.Query(userid)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}
