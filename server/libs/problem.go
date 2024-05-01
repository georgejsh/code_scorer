package libs

import (
	"encoding/json"
	"net/http"
	//"io/ioutil"
	"strings"
	"os"
	"io"
	//"fmt"
)

var input_problem =[]Help {
	Help{"problem_id","string",false,[]string{}},
}
func GetProblem(w http.ResponseWriter, r *http.Request) {
	var head string
	
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		w.Header().Set("Content-Type", "application/json")
		output, _ := json.Marshal(input_problem)
		data, _ := json.Marshal(Output{string(output), ""})
		
		w.Write(data)
		return
	}
	err := r.ParseMultipartForm(0)

	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	claims,err:=verifyLogin(r)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	course:=getSelectedCourse(claims["user"].(string))
	if course=="Not Selected"{
		data, _ := json.Marshal(Output{"","No Course Selected yet!!"})
		w.Write(data)
		return
	}
	pass:=selectProblem(claims["user"].(string),course,r.FormValue("problem_id"))
		//fmt.Printf("results %s\n",result)
	if pass != true {
		// Handle error
		data, _ := json.Marshal(Output{"","Wrong problem ID!!"})
		w.Write(data)
		return
	}
	desc:=getProblem(r.FormValue("problem_id"))
	if desc=="NA"{
		data, _ := json.Marshal(Output{"", "Invalid problem ID!!"})
		w.Write(data)
		return
	}
	
	if strings.Contains(desc,".pdf"){
		w.Header().Set("Content-Type", "application/pdf")
	}else if strings.Contains(desc,".jpg"){ 
		w.Header().Set("Content-Type", "image/jpeg")
	}
	file, err := os.Open(getProblemPath(desc))
    if err != nil {
        data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
    }
    defer file.Close()

    // Set appropriate headers
    w.Header().Set("Content-Disposition", "attachment; filename="+desc)

    // Pipe the file contents to the response writer
    _, err = io.Copy(w, file)
    if err != nil {
        data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
    }
	//data, _ := json.Marshal(Output{t.Render(), ""})
	//w.Write(data)
	return
}
