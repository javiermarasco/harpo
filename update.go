package main

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func UpdateAWSSecret(path string, secretname string, newvalue string) error {
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

	input := secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secretNameHashed),
		SecretString: aws.String(newvalue),
	}
	_, errupdate := conn.UpdateSecret(context.TODO(), &input)
	if errupdate != nil {
		return err
	}
	return nil

}
func UpdateGCPSecret(path string, secretname string, newvalue string) error {
	// Try to get the parent id for GCP
	gcp_parent := os.Getenv("GCP_parent")
	if gcp_parent == "" {
		fmt.Println(" Environment variable GCP_parent needs to be defined with format 'projects/parentid'.")
		os.Exit(1)
	}

	//parent := "projects/842557969287"
	payload := []byte(newvalue)
	inputForHash := path + "+" + secretname
	secretNameHashed := CreateHash(inputForHash)

	ctx := context.Background()
	client, errclient := secretmanager.NewClient(ctx)
	if errclient != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", errclient)
	}
	defer client.Close()

	reqver := &secretmanagerpb.AddSecretVersionRequest{
		Parent: gcp_parent + "/secrets/" + secretNameHashed,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}
	resultver, errver := client.AddSecretVersion(ctx, reqver)
	if errver != nil {
		return fmt.Errorf("failed to add secret version: %v", errver)
	}

	fmt.Println("Secret updated: ", resultver.Name)
	return nil

}
