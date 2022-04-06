package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func DeleteAzSecret(path string, secret secret_struct, creds *auth) {
	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
	inputForHash := path + "+" + secret.Name
	secretNameHashed := CreateHash(inputForHash)

	uri := fmt.Sprint(base_uri, "/secrets/", secretNameHashed, "?api-version=7.2")

	// DELETE
	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully deleted secret from Azure Keyvault")
	} else {
		fmt.Println("Delete failed with error: ", resp)
	}
}

func DeleteAwsSecret(path string, secretname string) {
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

	deletesecret := secretsmanager.DeleteSecretInput{
		SecretId: aws.String(secretNameHashed),
	}
	_, err = conn.DeleteSecret(context.TODO(), &deletesecret)
	if err != nil {
		fmt.Println("Error deleting secret.")
	}
}

func DeleteGCPSecret(path string, secretname string) {

	// Try to get the parent id for GCP
	gcp_parent := os.Getenv("GCP_parent")
	if gcp_parent == "" {
		fmt.Println(" Environment variable GCP_parent needs to be defined with format 'projects/parentid'.")
		os.Exit(1)
	}

	//Parent := "projects/842557969287"
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer c.Close()
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: gcp_parent + "/secrets/" + secretname,
	}

	err = c.DeleteSecret(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
	}
}
