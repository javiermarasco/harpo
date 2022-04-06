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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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

func GetAwsSecrets(path string) ([]types.SecretListEntry, error) {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	conn := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		o.Region = region
	})
	values := []string{path}

	myfiltervalues := types.Filter{
		Key:    "tag-value",
		Values: values,
	}

	myfilters := []types.Filter{myfiltervalues}

	mylistsecretinput := secretsmanager.ListSecretsInput{
		Filters: myfilters,
	}

	result, err := conn.ListSecrets(context.TODO(), &mylistsecretinput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result.SecretList, nil
}

func GetGcpSecrets(path string) {

	// Try to get the parent id for GCP
	gcp_parent := os.Getenv("GCP_parent")
	if gcp_parent == "" {
		fmt.Println(" Environment variable GCP_parent needs to be defined with format 'projects/parentid'.")
		os.Exit(1)
	}

	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer c.Close()

	req := &secretmanagerpb.ListSecretsRequest{
		//Parent: "projects/842557969287",
		Parent: gcp_parent,
	}
	it := c.ListSecrets(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("An error ocurred while retrieving the secret from GCP.")
		}

		pathslice := strings.Split(path, "/")
		respsecretpathslice := strings.Split(strings.Replace(resp.Labels["path"], "_", "/", -1), "/")

		pathmath := false
		if len(pathslice) <= len(respsecretpathslice) {
			for i, _ := range pathslice {
				if pathslice[i] == respsecretpathslice[i] {
					pathmath = true
				} else {
					pathmath = false
				}
			}
		}
		if pathmath == true {
			fmt.Println("The path for the secret is: ", strings.Replace(resp.Labels["path"], "_", "/", -1)+"/"+resp.Labels["secretname"])
		}
	}
}
