package libs
import (
	"strings"
	"io/ioutil"
	"fmt"
)
func checker(input_file,output_file,submit_file,checker_type string) int{
	//input, _ := ioutil.ReadFile(input_file)
	output, _ := ioutil.ReadFile(output_file)
	submit, _ := ioutil.ReadFile(submit_file)
	fmt.Printf("files %s %s \n",output_file, submit_file)
	if checker_type=="EQUAL"{
		fmt.Printf("%s %s \n",string(output), string(submit))
		if  strings.Compare(string(output), string(submit)) == 0 {
			return 1
		}else {
			return 0
		}
	}
	return 0
}