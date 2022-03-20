package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func WriteAzSecret(path string, secret secret_struct, creds *auth) {
	base_uri := fmt.Sprint("https://", creds.KeyVault, ".vault.azure.net")
	inputForHash := path + "+" + secret.Name
	secretName := CreateHash(inputForHash)

	uri := fmt.Sprint(base_uri, "/secrets/", secretName, "?api-version=7.2")

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
		panic(err)
	}

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
	} else {
		fmt.Println("Get failed with error: ", resp)
	}

}

func WriteAWSSecret(path string, secretname string, secretvalue string) {
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
		fmt.Println("Error writing secret.")
	}
}
