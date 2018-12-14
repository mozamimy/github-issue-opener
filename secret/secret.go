package secret

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func FetchSecretValue(secretName string, secretKey string) (string, error) {
	retVal := ""

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

	var secret map[string]string
	err = json.Unmarshal(([]byte)(*getSecretValueResp.SecretString), &secret)
	if err != nil {
		return retVal, err
	}
	retVal = secret[secretKey]

	return retVal, nil
}
