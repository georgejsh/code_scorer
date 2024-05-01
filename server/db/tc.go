package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"io"
	//"fmt"
	"strconv"
)

type TC struct {
	id string `json:"id"`
	input   string `json:"input"`
	output string `json:"output"`
	checker string `json:"checker"`
}
func AddTC(sample,problem_id string) int64{
	statement, _ := db_obj.Prepare("INSERT INTO tc (sample, problem_id) VALUES (?, ?)")
	res, _ :=statement.Exec(sample,problem_id)
	id, _ := res.LastInsertId()
    return id
}
func Updatefiles(input,output string,id int64){
	statement, _ := db_obj.Prepare("update tc set input=? and output =? where id=?")
	statement.Exec(input,output,id)
}

func GetTC(id string) TC{
	statement, _ := db_obj.Prepare("select id,input,output from tc where id=?")
	rows,_ :=statement.Query(id)
	defer rows.Close()
	for rows.Next() {
		var tempC TC
		rows.Scan(&tempC.id,&tempC.input,&tempC.output)
		return tempC
	}
	return TC{}
}

func GetSampleTC(userid,probid string) (string){
	statement, _ := db_obj.Prepare("select id,input, output from tc where problem_id=? and sample=\"true\" order by random()")
	rows,_ :=statement.Query(probid)
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	tcid:=AddTCFile(userid,probid,"true")
	os.Mkdir(getFolderNameForTC(tcid),os.FileMode(0666))

	newFileInputPath := GetInputFileTC(tcid)
    newFileInput, _ := os.Create(newFileInputPath)

    defer newFileInput.Close()

	//newFileOutputPath := "./files/"+string(tcid)+"_sample_output.txt"
    //newFileOutput, _ := os.Create(newFileOutputPath)
    //defer newFileOutput.Close()
	var inp []string
	var tc_ids = ""
	for rows.Next() {
		var tempC TC
		rows.Scan(&tempC.id,&tempC.input,&tempC.output)
		inp=append(inp,tempC.input)

		//outfile, _ := os.Open(tempC.output)
		//fmt.Printf("%s\n",tempC.id)
		//defer outfile.Close()
		//io.Copy(newFileOutput, outfile)

		tc_ids+=tempC.id+","
	}
	newFileInput.WriteString(strconv.Itoa(len(inp)) + "\n")
	for i := 0; i < len(inp); i++{
		inpfile, _ := os.Open(GetProblemPath(inp[i]))
		defer inpfile.Close()
		io.Copy(newFileInput, inpfile)
		newFileInput.WriteString("\n")
	}
	UpdateTCFile(tc_ids,strconv.FormatInt(tcid,10))
	return newFileInputPath
}

func GetFullTC(userid,probid string) (string){
	statement, _ := db_obj.Prepare("select id,input, output from tc where problem_id=? and sample=\"false\" order by random()")
	rows,_ :=statement.Query(probid)
	defer rows.Close()
	//fmt.Printf("Reached here\n")
	
	tcid:=AddTCFile(userid,probid,"false")
	os.Mkdir(getFolderNameForTC(tcid),os.FileMode(0666))
	newFileInputPath :=  GetInputFileTC(tcid)
    newFileInput, _ := os.Create(newFileInputPath)
    defer newFileInput.Close()

	//newFileOutputPath := "./files/"+string(tcid)+"_output.txt"
    //newFileOutput, _ := os.Create(newFileOutputPath)
    //defer newFileOutput.Close()
	
	var inp []string
	var tc_ids = ""
	for rows.Next() {
		var tempC TC
		rows.Scan(&tempC.id,&tempC.input,&tempC.output)
		inp=append(inp,tempC.input)

		//outfile, _ := os.Open(tempC.output)
		//fmt.Printf("%s\n",tempC.id)
		//defer outfile.Close()
		//io.Copy(newFileOutput, outfile)

		tc_ids+=tempC.id+","
	}
	newFileInput.WriteString(strconv.Itoa(len(inp)) + "\n")
	for i := 0; i < len(inp); i++{
		inpfile, _ := os.Open(GetProblemPath(inp[i]))
		defer inpfile.Close()
		io.Copy(newFileInput, inpfile)
		newFileInput.WriteString("\n")
	}
	UpdateTCFile(tc_ids,strconv.FormatInt(tcid,10))
	return newFileInputPath
}