package parameter

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mozamimy/lambda-github-issue-opener/secret"
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
	IssueLabel              string
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

func New(snsEvent events.SNSEvent) (Parameter, error) {
	parameter := Parameter{}

	parameter.GitHubAppKey = os.Getenv("GITHUB_APP_KEY")
	if parameter.GitHubAppKey == "" {
		return Parameter{}, fmt.Errorf("GITHUB_APP_KEY variable is required")
	}

	parameter.GitHubBaseURL = os.Getenv("GITHUB_BASE_URL")
	if parameter.GitHubBaseURL == "" {
		return Parameter{}, fmt.Errorf("GITHUB_BASE_URL variable is required")
	}

	var err error
	parameter.GitHubIntegrationID, err = fetchEnvVarAsInt("GITHUB_INTEGRATION_ID")
	if err != nil {
		return Parameter{}, err
	}

	parameter.GitHubInstallationID, err = fetchEnvVarAsInt("GITHUB_INSTALLATION_ID")
	if err != nil {
		return Parameter{}, err
	}

	gitHubAppPrivateKeySecretName := os.Getenv("GITHUB_APP_PRIVATE_KEY_SECRET_NAME")
	if gitHubAppPrivateKeySecretName == "" {
		return Parameter{}, fmt.Errorf("GITHUB_APP_PRIVATE_KEY_SECRET_NAME variable is required")
	}
	secrets, err := secret.FetchSecret(gitHubAppPrivateKeySecretName)
	if err != nil {
		return Parameter{}, err
	}
	parameter.GitHubAppPrivateKey = secrets.GitHubAppPrivateKey

	parameter.GitHubAPIBaseURL = os.Getenv("GITHUB_API_BASE_URL")
	if parameter.GitHubAPIBaseURL == "" {
		return Parameter{}, fmt.Errorf("GITHUB_API_BASE_URL variable is required")
	}

	parameter.GitHubAPIUploadsBaseURL = os.Getenv("GITHUB_API_UPLOADS_BASE_URL")
	if parameter.GitHubAPIUploadsBaseURL == "" {
		return Parameter{}, fmt.Errorf("GITHUB_API_UPLOADS_BASE_URL variable is required")
	}

	parameter.RepositoryOwner = os.Getenv("REPOSITORY_OWNER")
	if parameter.RepositoryOwner == "" {
		return Parameter{}, fmt.Errorf("REPOSITORY_OWNER variable is required")
	}
	parameter.Repository = os.Getenv("REPOSITORY")
	if parameter.Repository == "" {
		return Parameter{}, fmt.Errorf("REPOSITORY variable is required")
	}

	parameter.IssueLabel = os.Getenv("ISSUE_LABEL")

	issueSubjectTemplateStr := os.Getenv("ISSUE_SUBJECT_TEMPLATE")
	if issueSubjectTemplateStr == "" {
		return Parameter{}, fmt.Errorf("ISSUE_SUBJECT_TEMPLATE variable is required")
	}
	issueBodyTemplateStr := os.Getenv("ISSUE_BODY_TEMPLATE")
	if issueBodyTemplateStr == "" {
		return Parameter{}, fmt.Errorf("ISSUE_BODY_TEMPLATE variable is required")
	}

	for _, record := range snsEvent.Records {
		snsEntity := record.SNS
		issueSubject, err := renderTemplate(issueSubjectTemplateStr, snsEntity)
		if err != nil {
			return Parameter{}, err
		}
		issueBody, err := renderTemplate(issueBodyTemplateStr, snsEntity)
		if err != nil {
			return Parameter{}, err
		}
		parameter.Issues = append(parameter.Issues, Issue{Subject: issueSubject, Body: issueBody})
	}
	return parameter, nil
}
