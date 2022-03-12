package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func WriteSecret(keyvault_name string, path string, secret secret_struct, creds *auth) {
	base_uri := fmt.Sprint("https://", keyvault_name, ".vault.azure.net")
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
