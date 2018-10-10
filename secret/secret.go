package secret

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secret struct {
	GitHubAppPrivateKey string `json:"GitHubAppPrivateKey"`
}

func FetchSecret(secretName string) (*Secret, error) {
	retVal := new(Secret)

	smSession, err := session.NewSession()
	if err != nil {
		return retVal, err
	}

	svc := secretsmanager.New(smSession)

	getSecretValueResp, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return retVal, err
	}

	err = json.Unmarshal(([]byte)(*getSecretValueResp.SecretString), retVal)
	if err != nil {
		return retVal, err
	}

	return retVal, nil
}
