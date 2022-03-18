package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetAzSecrets(creds *auth) secret_list {

	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
	uri := fmt.Sprint(base_uri, "/secrets?api-version=7.2")
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var response secret_list
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Secret not found, check the path or the secret name.")
	}
	return response
}

func GetAWSSecret(path string, region string, secretname string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	conn := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		o.Region = region
	})
	input := secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretname),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := conn.GetSecretValue(context.TODO(), &input)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
		fmt.Println(secretString)
	}
}
