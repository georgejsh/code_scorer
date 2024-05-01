package libs
import (
	"server/db"
	"os"
	//"bytes"
	//"bufio"
	//"io"
)

func checkPassword(userid,pass string) (bool){
	return db.CheckPassword(userid,pass)
}

func getPassword(userid string) (string,error){
	return db.GetPassword(userid)
}

func resetPassowrd(userid,pass string){
	db.ResetPassowrd(userid,pass)
}

func addUser(userid,pass string){
	db.AddUser(userid,pass)
}

func isAdmin(userid string) bool{
	if userid=="root" {
		return true
	}
	return db.IsAdmin(userid)
}

func isRoot(userid string) bool{
	if userid=="root" {
		return true
	}
	return false
}
func getCourses(courseid string) string{
	return db.GetCourses(courseid)
}

func getProblems(userid,courseid string) string{
	if isAdmin(userid){
		return db.GetProblems(courseid)
	}else {
		return db.GetActiveProblems(courseid)
	}
}
func clearSelection(userid string) {
	db.ClearSelection(userid)
}
func createSelection(userid string) {
	db.CreateSelection(userid)
}
func selectCourse(userid,courseid string) bool {
	if ! isCourse(userid,courseid) {
		return false
	}
	db.SelectCourse(userid,courseid)
	return true
}

func isCourse(userid,courseid string) bool {
	return db.IsCourse(userid,courseid)
}
func getSelectedCourse(userid string) string{
	return db.GetSelectedCourse(userid)
}

func getSelectedProblem(courseid string) string{
	return db.GetSelectedProblem(courseid)
}

func selectProblem(userid,courseid,probid string) bool {
	if ! isProblem(courseid,userid,probid) {return false}
	db.SelectProblem(userid,probid)
	return true
}
func isProblem(courseid,userid,probid string) bool {
	return db.IsProblem(courseid,userid,probid)
}

func getProblem(probid string) string{
	return db.GetProblem(probid)
}

func getSampleTC(userid,probid string) (string){
	prev:=isActiveSampleTC(userid,probid)
	if  prev!=-1{
		return  db.GetInputFileTC(prev)	
	}
	return db.GetSampleTC(userid,probid)
}
func isActiveSampleTC(userid,probid string)int64 {
	return db.IsActiveSampleTC(userid,probid)
}

func isActiveFullTC(userid,probid string)int64 {
	return db.IsActiveFullTC(userid,probid)
}
func getFullTC(userid,probid string) (string){
	prev:=isActiveFullTC(userid,probid)
	if  prev!=-1{
		return  db.GetInputFileTC(prev)
	}
	return db.GetFullTC(userid,probid)
}
func getProblemPath(prob_id string) string{
	return  db.GetProblemPath(prob_id)
}
/*
func getFolderNameForTC(tc_list_id int64){
	return  db.GetFolderNameForTC(tc_list_id)
}*/

func getOutputFileTC(tc_list_id,submitid int64) string {
	return  db.GetOutputFileTC(tc_list_id,submitid)
}
func getSourceFileTC(tc_list_id,submitid int64) string {
	return  db.GetSourceFileTC(tc_list_id,submitid)
}
/*
func getFullOutputFileTC(tc_list_id,submitid int64) string{
	return  db.GetFullOutputFileTC(tc_list_id,submitid)
}*/

func getOutputFileTCCase(tc_list_id,submitid,caseid int64) string{
	return db.GetOutputFileTCCase(tc_list_id,submitid,caseid)
}
/*
func getSampleTCOutput(userid,probid string) (string){
	perv:=IsActiveSampleTC(userid,probid)
	if  prev!=-1{
		return  "./files/"+prev"_sample_output.txt"
	}
	return "NA"
}
func getFullTCOutput(userid,probid string) (string){
	perv:=IsActiveFullTC(userid,probid)
	if  prev!=-1{
		return  "./files/"+prev"_output.txt"
	}
	return "NA"
}*/

func getTCList(tclistid int64) [][2]string {
	return db.GetTCList(tclistid)
}

func createSubmit(tclistid int64)int64{
	submitid:= db.CreateSubmit(tclistid)
	os.Mkdir(db.GetFolderNameForSubmit(tclistid,submitid),os.FileMode(0666))
	return submitid
}

func updateScore(submitid,score int64){
	db.UpdateScore(submitid,score)
}

func getChecker(probid string) string{
	return db.GetChecker(probid)
}


//add

func addCourse(id,name,year string){
	db.AddCourse(id,name,year)
}

func enroll(userid,courseid,isadmin string) {
	if isadmin == "true"{
		db.EnrollAdmin(userid,courseid)
	}else {
		db.Enroll(userid,courseid)
	}
}

func addProblem(id,title,desc,active,courseid,checker string) {
	db.AddProblem(id,title,desc,active,courseid,checker)
}

func addTC(sample,problem_id string) int64{
	return db.AddTC(sample,problem_id )
}

func updatefiles(input,output string,id int64){
	db.Updatefiles(input,output,id)
}