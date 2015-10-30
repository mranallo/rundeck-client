package cmd

import (
	"fmt"
	"github.com/lusis/go-rundeck/src/rundeck.v12"
	"strings"
)

func RunJob(projectid string, jobname string, options string) {
	var jobID string

	client := rundeck.NewClientFromEnv()

	jobByName, err := client.FindJobByName(jobname, projectid)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	jobID = jobByName.ID

	arguments := parseArguments(options)

	o := rundeck.RunOptions{LogLevel: "INFO", RunAs: "", Arguments: arguments}
	data, err := client.RunJob(jobID, o)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	var executionID string
	for _, d := range data.Executions {
		executionID = d.ID
	}
	GetExecutionstate(executionID, projectid)
	fmt.Printf("\nTo see the log from this execution, run 'rundeck-client execution output %s'\n\n", executionID)
}

func parseArguments(options string) string {
	//input is quoted, comma-separated key/pairs: option1=option1,option2=option2
	//output should be : '-option1 option1 -option2 option2'
	var arguments string
	if options != "" {
		j := strings.Split(options, ",")
		for _, p := range j {
			s := strings.Split(p, "=")
			k, v := s[0], s[1]
			k = "-" + k
			arguments = arguments + " " + k + " " + v
		}
	}
	return strings.Trim(arguments, " ")
}
