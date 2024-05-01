package libs

import (
	"encoding/json"
	"net/http"
)


func GetHelp(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	type task struct {
		Command string `json:"command"`
		Url     string `json:"url"`
	}
	var result []task
	for k, _ := range commands {
		if ! commands_access[k] {
			result = append(result, task{k, "http://localhost:3333/" + k+"/"})
		}
	}
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	output, _ := json.Marshal(result)

	data, _ := json.Marshal(Output{string(output), ""})
	w.Write(data)
}
