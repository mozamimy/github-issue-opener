package parameter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mozamimy/github-issue-opener/secret"
)

type Issue struct {
	Subject string
	Body    string
}

type Parameter struct {
	GitHubAppKey            string
	GitHubBaseURL           string
	GitHubIntegrationID     int
	GitHubInstallationID    int
	GitHubAppPrivateKey     string
	GitHubAPIBaseURL        string
	GitHubAPIUploadsBaseURL string
	RepositoryOwner         string
	Repository              string
	IssueLabels             []string
	Issues                  []Issue
}

func fetchEnvVarAsInt(key string) (int, error) {
	retVarStr := os.Getenv(key)
	if retVarStr == "" {
		return 0, fmt.Errorf("err %s variable is required", key)
	}

	retVar, err := strconv.Atoi(retVarStr)
	if err != nil {
		return 0, err
	}

	return retVar, nil
}

func renderTemplate(templateStr string, snsEntity events.SNSEntity) (string, error) {
	template := template.Must(template.New("template").Parse(templateStr))
	var templateBuffer bytes.Buffer
	err := template.Execute(&templateBuffer, snsEntity)
	if err != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

func splitAndTrim(str string, sep string) []string {
	ret := []string{}
	for _, splitStr := range strings.Split(str, sep) {
		trimmedStr := strings.TrimSpace(splitStr)
		if trimmedStr != "" {
			ret = append(ret, trimmedStr)
		}
	}
	return ret
}

func handleEnvVars(parameter *Parameter) (string, string, error) {
	var err error

	parameter.GitHubAppKey = os.Getenv("GITHUB_APP_KEY")
	if parameter.GitHubAppKey == "" {
		return "", "", fmt.Errorf("GITHUB_APP_KEY variable is required")
	}

	parameter.GitHubBaseURL = os.Getenv("GITHUB_BASE_URL")
	if parameter.GitHubBaseURL == "" {
		return "", "", fmt.Errorf("GITHUB_BASE_URL variable is required")
	}

	parameter.GitHubIntegrationID, err = fetchEnvVarAsInt("GITHUB_INTEGRATION_ID")
	if err != nil {
		return "", "", err
	}

	parameter.GitHubInstallationID, err = fetchEnvVarAsInt("GITHUB_INSTALLATION_ID")
	if err != nil {
		return "", "", err
	}

	gitHubAppPrivateKeySecretName := os.Getenv("GITHUB_APP_PRIVATE_KEY_SECRET_NAME")
	if gitHubAppPrivateKeySecretName == "" {
		return "", "", fmt.Errorf("GITHUB_APP_PRIVATE_KEY_SECRET_NAME variable is required")
	}
	secrets, err := secret.FetchSecret(gitHubAppPrivateKeySecretName)
	if err != nil {
		return "", "", err
	}
	parameter.GitHubAppPrivateKey = secrets.GitHubAppPrivateKey

	parameter.GitHubAPIBaseURL = os.Getenv("GITHUB_API_BASE_URL")
	if parameter.GitHubAPIBaseURL == "" {
		return "", "", fmt.Errorf("GITHUB_API_BASE_URL variable is required")
	}

	parameter.GitHubAPIUploadsBaseURL = os.Getenv("GITHUB_API_UPLOADS_BASE_URL")
	if parameter.GitHubAPIUploadsBaseURL == "" {
		return "", "", fmt.Errorf("GITHUB_API_UPLOADS_BASE_URL variable is required")
	}

	parameter.RepositoryOwner = os.Getenv("REPOSITORY_OWNER")
	if parameter.RepositoryOwner == "" {
		return "", "", fmt.Errorf("REPOSITORY_OWNER variable is required")
	}
	parameter.Repository = os.Getenv("REPOSITORY")
	if parameter.Repository == "" {
		return "", "", fmt.Errorf("REPOSITORY variable is required")
	}

	parameter.IssueLabels = splitAndTrim(os.Getenv("ISSUE_LABELS"), ",")

	issueSubjectTemplateStr := os.Getenv("ISSUE_SUBJECT_TEMPLATE")
	if issueSubjectTemplateStr == "" {
		return "", "", fmt.Errorf("ISSUE_SUBJECT_TEMPLATE variable is required")
	}
	issueBodyTemplateStr := os.Getenv("ISSUE_BODY_TEMPLATE")
	if issueBodyTemplateStr == "" {
		return "", "", fmt.Errorf("ISSUE_BODY_TEMPLATE variable is required")
	}

	return issueSubjectTemplateStr, issueBodyTemplateStr, nil
}

func handleMessageAttributes(snsEntity *events.SNSEntity, issueSubjectTemplateStr *string, issueBodyTemplateStr *string, parameter *Parameter) error {
	var err error
	for key, attr := range snsEntity.MessageAttributes {
		switch attrMap := attr.(type) {
		case map[string]interface{}:
			switch attrValue := attrMap["Value"].(type) {
			case string:
				switch attrMap["Type"] {
				case "String.Array":
					if key == "Labels" {
						var labels []string
						err := json.Unmarshal([]byte(attrValue), &labels)
						if err != nil {
							return fmt.Errorf("failed to unmarshal message attribute value: %v", attrValue)
						}
						parameter.IssueLabels = labels
					} else {
						return fmt.Errorf("unknown message attribute key: %v", key)
					}
				case "String":
					switch key {
					case "IssueSubjectTemplate":
						issueSubjectTemplateStr = &attrValue
					case "IssueBodyTemplate":
						issueBodyTemplateStr = &attrValue
					case "RepositoryOwner":
						parameter.RepositoryOwner = attrValue
					case "Repository":
						parameter.Repository = attrValue
					default:
						return fmt.Errorf("unknown message attribute key: %v", key)
					}
				case "Number":
					parameter.GitHubInstallationID, err = strconv.Atoi(attrValue)
					if err != nil {
						return fmt.Errorf("failed to convert from string to integer (key: %v, value: %v)", key, attrValue)
					}
				}
			default:
				return fmt.Errorf("unknown type is found in message attribute [%v]: %v", reflect.TypeOf(attrValue), attrValue)
			}
		default:
			return fmt.Errorf("unknown type is found in message attribute [%v]: %v", reflect.TypeOf(attrMap), attrMap)
		}
	}

	issueSubject, err := renderTemplate(*issueSubjectTemplateStr, *snsEntity)
	if err != nil {
		return err
	}
	issueBody, err := renderTemplate(*issueBodyTemplateStr, *snsEntity)
	if err != nil {
		return err
	}
	parameter.Issues = append(parameter.Issues, Issue{Subject: issueSubject, Body: issueBody})
	return nil
}

func New(snsEvent events.SNSEvent) (Parameter, error) {
	parameter := Parameter{}
	issueSubjectTemplateStr, issueBodyTemplateStr, err := handleEnvVars(&parameter)
	if err != nil {
		return Parameter{}, err
	}

	for _, record := range snsEvent.Records {
		snsEntity := record.SNS
		err = handleMessageAttributes(&snsEntity, &issueSubjectTemplateStr, &issueBodyTemplateStr, &parameter)
		if err != nil {
			return Parameter{}, err
		}
	}
	return parameter, nil
}
