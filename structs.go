package main

type secret_struct struct {
	Name  string            `json:"name"`
	Value string            `json:"value"`
	Tags  map[string]string `json:"tags"`
}

type secret_list struct {
	Value []struct {
		ID         string `json:"id"`
		Attributes struct {
			Enabled         bool   `json:"enabled"`
			Created         int    `json:"created"`
			Updated         int    `json:"updated"`
			RecoveryLevel   string `json:"recoveryLevel"`
			RecoverableDays int    `json:"recoverableDays"`
		} `json:"attributes"`
		Tags map[string]string `json:"tags"`
	} `json:"value"`
	NextLink string `json:"nextLink"`
}

type auth struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Resource     string `json:"resource"`
	TenantID     string `json:"tenant_id"`
	Token        string `json:"access_token"`
	KeyVault     string `json:"-"`
}
