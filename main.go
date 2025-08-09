package main

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func main() {

	// if os.Args[1] == "TODO" {
	// 	fmt.Println("This will be where prod function goes")
	// } else if os.Args[1] == "test" {
	// 	fmt.Println("This will be where test function goes")
	// } else {
	// 	readValuesDev()
	// }
	if os.Args[1] == "dev" {
		readValuesDev()
	}

}

func readValuesDev() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := ssm.NewFromConfig(cfg)

	output, err := client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:      aws.String("/Dev"),
		Recursive: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("/home/abi/testfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, param := range output.Parameters {
		if _, err := file.Write([]byte("define('" + (path.Base(aws.ToString(param.Name)) + "' ,'" + aws.ToString(param.Value)) + "');\n")); err != nil {
			log.Fatal(err)
		}
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
