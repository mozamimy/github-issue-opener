package parameter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestNew(t *testing.T) {
	snsEventRaw, err := ioutil.ReadFile("../test/fixtures/event_filled_with_attrs.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	var snsEvent events.SNSEvent
	err = json.Unmarshal(snsEventRaw, &snsEvent)
	if err != nil {
		t.Fatalf("%v", err)
	}

	os.Setenv("GITHUB_API_BASE_URL", "https://api.example.com")
	os.Setenv("GITHUB_API_UPLOADS_BASE_URL", "https://uploads.example.com")
	os.Setenv("GITHUB_APP_KEY", "dummy")
	// TODO: We should mock API call for AWS Secret Manager and replace with dummy secret key name
	os.Setenv("GITHUB_APP_PRIVATE_KEY_SECRET_NAME", "lambda/github-issue-opener")
	os.Setenv("GITHUB_BASE_URL", "https://example.com")
	os.Setenv("GITHUB_INSTALLATION_ID", "1")
	os.Setenv("GITHUB_INTEGRATION_ID", "2")
	os.Setenv("ISSUE_BODY_TEMPLATE", "{{.Message}}")
	os.Setenv("ISSUE_SUBJECT_TEMPLATE", "{{.Subject}}")
	os.Setenv("REPOSITORY", "syndrome")
	os.Setenv("REPOSITORY_OWNER", "rabbit")
	os.Setenv("ISSUE_LABELS", "duplicate,question")

	result, err := New(snsEvent)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// TODO: We should mock API call for AWS Secret Manager and replace with dummy secret key name
	result.GitHubAppPrivateKey = "dummy"
	issues := []Issue{
		Issue{
			Subject: "example subject: foo",
			Body:    "example message: bar",
		},
	}
	expect := Parameter{
		GitHubAppKey:            "dummy",
		GitHubBaseURL:           "https://example.com",
		GitHubIntegrationID:     2,
		GitHubInstallationID:    1,
		GitHubAppPrivateKey:     "dummy",
		GitHubAPIBaseURL:        "https://api.example.com",
		GitHubAPIUploadsBaseURL: "https://uploads.example.com",
		RepositoryOwner:         "fox",
		Repository:              "ailment",
		IssueLabels:             []string{"bug", "wontfix"},
		Issues:                  issues,
	}

	if reflect.DeepEqual(result, expect) {
		// ok
	} else {
		t.Fatalf("Do not same values, expect: %+v, got: %+v", expect, result)
	}

	snsEventRaw, err = ioutil.ReadFile("../test/fixtures/event.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	var snsEvent2 events.SNSEvent
	err = json.Unmarshal(snsEventRaw, &snsEvent2)
	if err != nil {
		t.Fatalf("%v", err)
	}

	result, err = New(snsEvent2)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// TODO: We should mock API call for AWS Secret Manager and replace with dummy secret key name
	result.GitHubAppPrivateKey = "dummy"
	issues = []Issue{
		Issue{
			Subject: "example subject",
			Body:    "example message",
		},
	}
	expect = Parameter{
		GitHubAppKey:            "dummy",
		GitHubBaseURL:           "https://example.com",
		GitHubIntegrationID:     2,
		GitHubInstallationID:    1,
		GitHubAppPrivateKey:     "dummy",
		GitHubAPIBaseURL:        "https://api.example.com",
		GitHubAPIUploadsBaseURL: "https://uploads.example.com",
		RepositoryOwner:         "rabbit",
		Repository:              "syndrome",
		IssueLabels:             []string{"duplicate", "question"},
		Issues:                  issues,
	}
	if reflect.DeepEqual(result, expect) {
		// ok
	} else {
		t.Fatalf("Do not same values, expect: %+v, got: %+v", expect, result)
	}

	// Test a case which has no ISSUE_LABALS environment variable
	os.Setenv("ISSUE_LABELS", "")

	result, err = New(snsEvent2)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// TODO: We should mock API call for AWS Secret Manager and replace with dummy secret key name
	result.GitHubAppPrivateKey = "dummy"
	expect = Parameter{
		GitHubAppKey:            "dummy",
		GitHubBaseURL:           "https://example.com",
		GitHubIntegrationID:     2,
		GitHubInstallationID:    1,
		GitHubAppPrivateKey:     "dummy",
		GitHubAPIBaseURL:        "https://api.example.com",
		GitHubAPIUploadsBaseURL: "https://uploads.example.com",
		RepositoryOwner:         "rabbit",
		Repository:              "syndrome",
		IssueLabels:             []string{},
		Issues:                  issues,
	}
	if reflect.DeepEqual(result, expect) {
		// ok
	} else {
		t.Fatalf("Do not same values, expect: %+v, got: %+v", expect, result)
	}

	os.Setenv("PARSE_JSON_MODE", "1")
	os.Setenv("ISSUE_BODY_TEMPLATE", "{{.ParsedMessage.rabbit}}")
	os.Setenv("ISSUE_SUBJECT_TEMPLATE", "{{.SNSEntity.Subject}}")

	snsEventRaw, err = ioutil.ReadFile("../test/fixtures/event_with_json_message.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	var snsEvent3 events.SNSEvent
	err = json.Unmarshal(snsEventRaw, &snsEvent3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	result, err = New(snsEvent3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	issues = []Issue{
		Issue{
			Subject: "example subject",
			Body:    "syndrome",
		},
	}
	// TODO: We should mock API call for AWS Secret Manager and replace with dummy secret key name
	result.GitHubAppPrivateKey = "dummy"
	expect = Parameter{
		GitHubAppKey:            "dummy",
		GitHubBaseURL:           "https://example.com",
		GitHubIntegrationID:     2,
		GitHubInstallationID:    1,
		GitHubAppPrivateKey:     "dummy",
		GitHubAPIBaseURL:        "https://api.example.com",
		GitHubAPIUploadsBaseURL: "https://uploads.example.com",
		RepositoryOwner:         "rabbit",
		Repository:              "syndrome",
		IssueLabels:             []string{},
		Issues:                  issues,
	}
	if reflect.DeepEqual(result, expect) {
		// ok
	} else {
		t.Fatalf("Do not same values, expect: %+v, got: %+v", expect, result)
	}
}
