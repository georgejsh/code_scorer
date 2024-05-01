package db
import (
	//"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	//"github.com/jedib0t/go-pretty/v6/table"
	//"fmt"
	"strconv"
)
var base_folder="./db/files/"
func getFolderNameForTC(tc_list_id int64) string{
	return  base_folder+strconv.FormatInt(tc_list_id,10)
}

func GetInputFileTC(tc_list_id int64) string {
	return  base_folder+strconv.FormatInt(tc_list_id,10)+"/input.txt"
}

func GetProblemPath(prob_id string)  string{
	return  base_folder+prob_id
}
func GetFolderNameForSubmit(tc_list_id,submitid int64) string{
	return  base_folder+strconv.FormatInt(tc_list_id,10)+"/"+strconv.FormatInt(submitid,10)
}
func GetOutputFileTC(tc_list_id,submitid int64) string {
	return  base_folder+strconv.FormatInt(tc_list_id,10)+"/"+strconv.FormatInt(submitid,10)+"/output.txt"
}

func GetSourceFileTC(tc_list_id,submitid int64) string {
	return  base_folder+strconv.FormatInt(tc_list_id,10)+"/"+strconv.FormatInt(submitid,10)+"/code.cpp"
}



func GetOutputFileTCCase(tc_list_id,submitid,caseid int64) string{
	return  base_folder+strconv.FormatInt(tc_list_id,10)+"/"+strconv.FormatInt(submitid,10)+"/"+strconv.FormatInt(caseid,10)
}
