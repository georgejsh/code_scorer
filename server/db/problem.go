package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jedib0t/go-pretty/v6/table"
	//"fmt"
)

type Problem struct {
	id     string `json:"id"`
	title  string `json:"title"`
	desc   string `json:"desc"`
	active string `json:"active"`
}



func AddProblem(id,title,desc,active,courseid,checker string) {
	statement, _ := db_obj.Prepare("INSERT INTO problem (id,title, desc, active, course_id,checker) VALUES (?,?, ?, ?, ?)")
	statement.Exec(id,title,desc,active,courseid,checker)
}

func GetActiveProblems(courseid string) string{
	statement, _ := db_obj.Prepare("select id,title,desc,active from problem where course_id=? and active<(select CURRENT_TIMESTAMP) ")
	rows,_ :=statement.Query(courseid)
	var problems []Problem
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	for rows.Next() {
		var tempC Problem
		rows.Scan(&tempC.id,&tempC.title , &tempC.desc, &tempC.active)
		//fmt.Printf("%s\n",tempC.id)
		problems=append(problems,tempC)
	}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "ID", "Title"})
	for i := 0; i < len(problems); i++ {
		t.AppendRow([]interface{}{i+1,problems[i].id, problems[i].title})
		t.AppendSeparator()
	}
	return t.Render()
}

func GetProblems(courseid string) string{
	statement, _ := db_obj.Prepare("select id,title,desc,active from problem where course_id=?")
	rows,_ :=statement.Query(courseid)
	var problems []Problem
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	for rows.Next() {
		var tempC Problem
		rows.Scan(&tempC.id,&tempC.title , &tempC.desc, &tempC.active)
		//fmt.Printf("%s\n",tempC.id)
		problems=append(problems,tempC)
	}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "ID", "Title", "Description", "Active Time"})
	for i := 0; i < len(problems); i++ {
		t.AppendRow([]interface{}{i+1,problems[i].id, problems[i].title, problems[i].desc, problems[i].active})
		t.AppendSeparator()
	}
	return t.Render()
}
func GetProblem(probid string) string {
	statement, _ := db_obj.Prepare("select desc from problem where active<(select CURRENT_TIMESTAMP) and id=?")
	rows,_ :=statement.Query(probid)
	var loc string
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	for rows.Next() {
		rows.Scan(&loc)
		return loc
	}
	return "NA"
}

func GetChecker(probid string) string {
	statement, _ := db_obj.Prepare("select checker from problem where id=?")
	rows,_ :=statement.Query(probid)
	var loc string
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	for rows.Next() {
		rows.Scan(&loc)
		return loc
	}
	return "NA"
}

func IsProblem(courseid,userid,problemid string) bool{
	statement, _ := db_obj.Prepare("select id from enroll join problem on enroll.course_id=problem.course_id where enroll.course_id=? and user_id=? and problem.id=? ")
	rows,_ :=statement.Query(courseid,userid,problemid)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}




