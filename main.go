package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func main() {

	readValues()

}

func readValues() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := ssm.NewFromConfig(cfg)

	var ppm_env string
	switch os.Args[1] {
	case "dev":
		ppm_env = "/Dev"
	case "test":
		ppm_env = "/Test"
	case "prod":
		ppm_env = "/Prod"
	default:
		fmt.Println("/Dev")
	}

	output, err := client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:      aws.String(ppm_env),
		Recursive: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
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
