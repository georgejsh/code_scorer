package libs

import (
	"encoding/json"
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"regexp"
	"fmt"
)

var input_submit =[]Help {
	Help{"name","enum",false,[]string{"sample","full"}},
	Help{"output","file",false,[]string{}},
	Help{"code","file",false,[]string{}},
}
func Submit(w http.ResponseWriter, r *http.Request) {
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_submit)
		data, _ := json.Marshal(Output{string(output), ""})
		
		w.Write(data)
		return
	}
	claims,err := getProblemID(r)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	fmt.Printf("problem id: %d \n",claims["problem"].(string))
	var tcid int64
	if r.FormValue("name") == "sample"{
		tcid=isActiveSampleTC(claims["user"].(string),claims["problem"].(string))
		
	}else{
		tcid=isActiveFullTC(claims["user"].(string),claims["problem"].(string))
		
	}
	if tcid==-1{
		data, _ := json.Marshal(Output{"", "Submission time exceeded!!"})
		w.Write(data)
		return
	}
	fmt.Printf("tcid: %d \n",tcid)
	output,_,_ := r.FormFile("output")
    if err != nil {
        // Handle the error (e.g., log it or return an error response)
        data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
    }
	defer output.Close()

	source,_,_ := r.FormFile("code")
    if err != nil {
        // Handle the error (e.g., log it or return an error response)
        data, _ := json.Marshal(Output{"", err.Error()+"here"})
		w.Write(data)
		return
    }
	defer source.Close()
	
	submitid:= createSubmit(tcid)
	fmt.Printf("submitid: %d \n",submitid)
	out_file:=getOutputFileTC(tcid,submitid)
	code_file:=getSourceFileTC(tcid,submitid)
	

	dst, err := os.Create(out_file)
	if err != nil {
		data, _ := json.Marshal(Output{"", "Output save error!!"})
		w.Write(data)
		return
	}
	//dst.WriteString(output)
	
	io.Copy(dst,output)
	dst.Close()
	dst, err = os.Create(code_file)
	if err != nil {
		data, _ := json.Marshal(Output{"", "Output save error!!"})
		w.Write(data)
		return
	}
	io.Copy(dst,source)
	//dst.WriteString(source)
	dst.Close()
	

	result:=getTCList(tcid)


	content, err := ioutil.ReadFile(out_file)    
	t := regexp.MustCompile(`Case [0-9]*:`) 
	outpertc := t.Split(string(content), -1) 
	if len(result)!=len(outpertc)-1{
		data, _ := json.Marshal(Output{"", "Presentation Error(Number of testcases in output file wrong)!! "})
		w.Write(data)
		return
	}
	score:=0
	for i:=0;i<len(outpertc)-1;i++ {
		tc_tmp_file:=getOutputFileTCCase(tcid,submitid,int64(i))
		os.WriteFile(tc_tmp_file, []byte(outpertc[i+1]), 0666)
		score+=checker(getProblemPath(result[i][0]),getProblemPath(result[i][1]),tc_tmp_file,getChecker(claims["problem"].(string)))
	}

	updateScore(submitid,int64(score))
	
	if score==len(outpertc)-1{
		data, _ := json.Marshal(Output{"All Testcases Passed!!", ""})
		w.Write(data)
		return
	}else{
		data, _ := json.Marshal(Output{"Wrong Answer!!", ""})
		w.Write(data)
		return
	}

	
	//var desc string
	
	
	return
}
