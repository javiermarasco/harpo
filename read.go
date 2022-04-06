package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func ReadAzSecret(path string, secret_name secret_struct, creds *auth) (value string, error string) {
	secretNameHashed := CreateHash(path + "+" + secret_name.Name)
	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
	uri := fmt.Sprint(base_uri, "/secrets/", secretNameHashed, "?api-version=7.2")
	value = ""
	error = ""
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var response secret_struct

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		value = response.Value

	} else if resp.StatusCode == http.StatusNotFound {
		error = "Secret not found, check the path or the secret name."
	}
	return value, error
}

func ReadAWSSecret(path string, secretname string) (string, error) {
	region := os.Getenv("AWS_REGION")
	inputForHash := path + "+" + secretname
	secretNameHashed := CreateHash(inputForHash)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	conn := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		o.Region = region
	})
	mysecretvalueinput := secretsmanager.GetSecretValueInput{
		SecretId: &secretNameHashed,
	}
	mysecret, err := conn.GetSecretValue(context.TODO(), &mysecretvalueinput)

	if err != nil {
		return "", err
	}
	return aws.ToString(mysecret.SecretString), nil
}

func ReadGCPSecret(path string, secretname string) (string, error) {
	// Try to get the parent id for GCP
	gcp_parent := os.Getenv("GCP_parent")
	if gcp_parent == "" {
		fmt.Println(" Environment variable GCP_parent needs to be defined with format 'projects/parentid'.")
		os.Exit(1)
	}
	//parent := "projects/842557969287"
	inputForHash := path + "+" + secretname
	secretNameHashed := CreateHash(inputForHash)
	// Create the client.
	ctx := context.Background()
	client, errclient := secretmanager.NewClient(ctx)
	if errclient != nil {
		return "", errclient
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: gcp_parent + "/secrets/" + strings.ToLower(secretNameHashed) + "/versions/latest",
	}

	secret, errget := client.AccessSecretVersion(ctx, req)
	if errget != nil {
		return "", errget
	}

	return string(secret.Payload.Data), nil
}
