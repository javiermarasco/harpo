package main

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func WriteAzSecret(path string, secret secret_struct, creds *auth) error {
	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
	inputForHash := path + "+" + secret.Name
	secretNameHashed := CreateHash(inputForHash)

	uri := fmt.Sprint(base_uri, "/secrets/", secretNameHashed, "?api-version=7.2")

	var newtags = make(map[string]string)
	newtags = PathToTags(path)

	newtags["SecretName"] = secret.Name

	type Secret struct {
		Tags  map[string]string `json:"tags,omitempty"`
		Value string            `json:"value,omitempty"`
	}

	secreto := Secret{
		Value: secret.Value,
		Tags:  newtags,
	}

	jsonData, _ := json.Marshal(secreto)
	// PUT
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprint("bearer ", creds.Token))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
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
	} else {
		fmt.Println("Get failed with error: ", resp)
		return err
	}
	return nil
}

func WriteAWSSecret(path string, secretname string, secretvalue string) error {
	region := os.Getenv("AWS_REGION")
	inputForHash := path + "+" + secretname
	secretNameHashed := CreateHash(inputForHash)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	conn := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		o.Region = region
	})

	tags := []types.Tag{}
	secretnametag := types.Tag{Key: aws.String("SecretName"), Value: aws.String(secretname)}
	pathtag := types.Tag{Key: aws.String("Path"), Value: aws.String(path)}
	tags = append(tags, secretnametag)
	tags = append(tags, pathtag)

	input := secretsmanager.CreateSecretInput{
		Name:         aws.String(secretNameHashed),
		SecretString: aws.String(secretvalue),
		Tags:         tags,
	}

	_, err = conn.CreateSecret(context.TODO(), &input)
	if err != nil {
		return err
	}
	return nil
}

func WriteGCPSecret(path string, secretname string, secretvalue string) error {
	// Try to get the parent id for GCP
	gcp_parent := os.Getenv("GCP_parent")
	if gcp_parent == "" {
		fmt.Println(" Environment variable GCP_parent needs to be defined with format 'projects/parentid'.")
		os.Exit(1)
	}
	//parent := "projects/842557969287"
	payload := []byte(secretvalue)
	inputForHash := path + "+" + secretname
	secretNameHashed := CreateHash(inputForHash)

	// Create the client.
	ctx := context.Background()
	client, errclient := secretmanager.NewClient(ctx)
	if errclient != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", errclient)
	}
	defer client.Close()

	gcplabels := map[string]string{
		"secretname": strings.ToLower(secretname),
		"path":       PathToGCP(path),
	}

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   gcp_parent,
		SecretId: secretNameHashed,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			Labels: gcplabels,
		},
	}

	// Call the API.
	_, errcreate := client.CreateSecret(ctx, req)
	if errcreate != nil {
		return fmt.Errorf("failed to create secret: %v", errcreate)
	}

	// Build the request.
	reqver := &secretmanagerpb.AddSecretVersionRequest{
		Parent: gcp_parent + "/secrets/" + secretNameHashed,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	resultver, errver := client.AddSecretVersion(ctx, reqver)
	if errver != nil {
		return fmt.Errorf("failed to add secret version: %v", errver)
	}

	fmt.Println("Secret created: ", resultver.Name)
	return nil
}
