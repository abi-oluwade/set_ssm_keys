package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func main() {

	directory, err := os.ReadDir(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range directory {
		current_file := (os.Args[2] + "/" + file.Name())
		fmt.Println("Reading file ==> " + current_file)
		readValues(current_file)
	}

}

func readValues(file string) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client := ssm.NewFromConfig(cfg)

	var ssm_env string
	switch os.Args[1] {
	case "dev":
		ssm_env = "/Dev"
	case "test":
		ssm_env = "/Test"
	case "prod":
		ssm_env = "/Prod"
	default:
		ssm_env = "/Dev"
	}

	if os.Args[1] != "dev" && os.Args[1] != "test" && os.Args[1] != "prod" {
		log.Print("Cannot find SSM Store environment provided, must be either 'dev','test' or 'prod'")
		log.Fatal()
	}

	param_output, err := client.GetParametersByPath(context.Background(), &ssm.GetParametersByPathInput{
		Path:           aws.String(ssm_env),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	file_input, err := os.ReadFile(file)
	if err != nil || !strings.Contains(file, "/") {
		log.Print("Cannot open or read to file path provided. Moving onto the next...")
	}

	lines := strings.Split(string(file_input), "\n")

	for i, line := range lines {
		for _, param := range param_output.Parameters {

			if strings.Contains(line, path.Base(aws.ToString(param.Name))) {
				fmt.Println("Replacing Placeholder in " + line + " with ==> " + aws.ToString(param.Value))

				lines[i] = strings.ReplaceAll(line, "@"+path.Base(aws.ToString(param.Name))+"@", aws.ToString(param.Value))

				updated_file := strings.Join(lines, "\n")
				err = os.WriteFile(file, []byte(updated_file), 0644)
				if err != nil {
					log.Fatal(err)
				}

			}

		}
	}
}
