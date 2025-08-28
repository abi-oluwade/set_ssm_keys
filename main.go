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

	readValues()

}

func readValues() {
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

	output, err := client.GetParametersByPath(context.Background(), &ssm.GetParametersByPathInput{
		Path:           aws.String(ssm_env),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	// TODO - Add functionality to loop through files in directory

	input, err := os.ReadFile(os.Args[2])
	if err != nil || !strings.Contains(os.Args[2], "/") {
		log.Print("Cannot open or read to file path provided.")
		log.Fatal()
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		for _, param := range output.Parameters {

			if strings.Contains(line, path.Base(aws.ToString(param.Name))) {
				fmt.Println("Replacing Placeholder in -->" + line)

				if strings.Contains(os.Args[2], "test") {

					lines[i] = strings.ReplaceAll(line, path.Base(aws.ToString(param.Name)+"_TOKEN"), aws.ToString(param.Value))
				} else {
					lines[i] = strings.ReplaceAll(line, "@"+path.Base(aws.ToString(param.Name))+"@", aws.ToString(param.Value))
				}

				output2 := strings.Join(lines, "\n")
				err = os.WriteFile(os.Args[2], []byte(output2), 0644)
				if err != nil {
					log.Fatal(err)
				}

			}

		}
	}
}
