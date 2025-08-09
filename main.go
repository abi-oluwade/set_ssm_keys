package main

import (
	"context"
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

	file, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil || !strings.Contains(os.Args[2], "/") {
		log.Print("Cannot open or read to file path provided.")
		log.Fatal()
	}

	for _, param := range output.Parameters {
		if _, err := file.Write([]byte("define('" + (path.Base(aws.ToString(param.Name)) + "' ,'" + aws.ToString(param.Value)) + "');\n")); err != nil {
			log.Fatal(err)
		}

	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
