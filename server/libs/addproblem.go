package libs

import (
	"encoding/json"
	"net/http"
	"os"
	"io"
	"strconv"
	//"strings"
	//"fmt"
	//"log"
)
var input_addproblem =[]Help {
	Help{"id","string",false,[]string{}},
	Help{"title","string",false,[]string{}},
	Help{"desc","file",false,[]string{}},
	Help{"active","string",false,[]string{}},
	Help{"courseid","string",false,[]string{}},
	Help{"checker","string",false,[]string{}},
}

func AddProblem(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_addproblem)
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
	// to do verify admin
	if isAdmin(claims["user"].(string)) ==false{
		data, _ := json.Marshal(Output{"", "Not Authorised User"})
		w.Write(data)
		return
	}
	if !isRoot(claims["user"].(string)) && getSelectedCourse(claims["user"].(string))!=r.FormValue ("courseid"){
		data, _ := json.Marshal(Output{"", "Not Authorised User. Select the course"})
		w.Write(data)
		return
	}
	
	prd_desc, multipartFileHeader, err := r.FormFile("desc")
	out_file:=r.FormValue ("id")+"_"+multipartFileHeader.Filename
	dst, err := os.Create(getProblemPath(out_file))
	if err != nil {
		data, _ := json.Marshal(Output{"", "Output save error!!"})
		w.Write(data)
		return
	}
	io.Copy(dst,prd_desc)
	//dst.WriteString(prd_desc)

	addProblem(r.FormValue ("id"),r.FormValue ("title"),out_file,r.FormValue ("active"),r.FormValue ("courseid"),r.FormValue ("checker"))
	data, _ := json.Marshal(Output{"Course Added Successfully", ""})
	w.Write(data)
}



var input_addtc =[]Help {
	Help{"input","file",false,[]string{}},
	Help{"output","file",false,[]string{}},
	Help{"sample","string",false,[]string{}},
	Help{"problem_id","string",false,[]string{}},
}

func AddTC(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_addproblem)
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
	// to do verify admin
	if isAdmin(claims["user"].(string)) ==false{
		data, _ := json.Marshal(Output{"", "Not Authorised User"})
		w.Write(data)
		return
	}
	if !isRoot(claims["user"].(string)) && getSelectedCourse(claims["user"].(string))!=r.FormValue ("courseid"){
		data, _ := json.Marshal(Output{"", "Not Authorised User. Select the course"})
		w.Write(data)
		return
	}
	tcid:= addTC(r.FormValue ("sample"),r.FormValue ("problem_id"))
	prd_input, multipartFileHeader, err := r.FormFile("input")
	inp_file:=strconv.FormatInt(tcid,10)+"_"+multipartFileHeader.Filename
	dst, err := os.Create(getProblemPath(inp_file))
	if err != nil {
		data, _ := json.Marshal(Output{"", "Input save error!!"})
		w.Write(data)
		return
	}
	io.Copy(dst,prd_input)

	prd_output, multipartFileHeader, err := r.FormFile("output")
	out_file:=strconv.FormatInt(tcid,10)+"_"+multipartFileHeader.Filename
	dst, err = os.Create(getProblemPath(out_file))
	if err != nil {
		data, _ := json.Marshal(Output{"", "Output save error!!"})
		w.Write(data)
		return
	}
	io.Copy(dst,prd_output)
	//dst.WriteString(prd_desc)
	updatefiles(inp_file,out_file,tcid)
	
	data, _ := json.Marshal(Output{"Course Added Successfully", ""})
	w.Write(data)
}