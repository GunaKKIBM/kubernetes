package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	docgen "k8s.io/kubernetes/test/e2e_node/docgenerator"
)

// Parse file
// Title Level
// Code block
// Paragraph
// Insert space, leve

var filteredData = make(map[docgen.TestCaseType][]interface{})

func Parsefile() error {

	file, err := os.ReadFile(os.Getenv("JSON_LOG"))
	if err != nil {
		return  err
	}
	err = json.Unmarshal(file, &filteredData)
	if err != nil {
		return  err
	}
	return nil	
}

 

func main() {
	err := Parsefile()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(filteredData)

	var builder strings.Builder

	for testcasetype, tests := range filteredData {
		builder.WriteString(docgen.TitleLevel(string(testcasetype), 1))

		for _, test := range tests {
			bytes, err := json.Marshal(test)
			if err != nil {
				fmt.Println("Error marshalling", err)
			}
			var podOP = new(docgen.PodConditionsTestOutput)
			err = json.Unmarshal(bytes, podOP)
			if err != nil {
				fmt.Println("Error unmarshalling", err)
			}
			podOP.Generate(&builder)
		}
	}

	err = os.WriteFile(os.Getenv("DOC_PATH"), []byte(builder.String()), 0644)
        	if err != nil {
                fmt.Printf("Error writing file: %v\n", err)
			}
}
