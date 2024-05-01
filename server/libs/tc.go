package libs

import (
	"encoding/json"
	"net/http"
	"os"
	"io"
	//"fmt"
)

var input_tc =[]Help {
	Help{"name","enum",false,[]string{"sample","full"}},
}
func GetTC(w http.ResponseWriter, r *http.Request) {
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_tc)
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
	var desc string
	var filename string
	if r.FormValue("name") == "sample"{
		desc=getSampleTC(claims["user"].(string),claims["problem"].(string));
		w.Header().Set("Content-Type", "text/plain")
		filename="sample.txt"

	}else{
		desc=getFullTC(claims["user"].(string),claims["problem"].(string));
		w.Header().Set("Content-Type", "text/plain")
		filename="full.txt"
	}
	file, err := os.Open(desc)
    if err != nil {
        data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
    }
    defer file.Close()

    // Set appropriate headers
    w.Header().Set("Content-Disposition", "attachment; filename="+filename)

    // Pipe the file contents to the response writer
    _, err = io.Copy(w, file)
    if err != nil {
        data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
    }
	return
}
