package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/mozamimy/lambda-github-issue-opener/parameter"
)

func HandleRequest(snsEvent events.SNSEvent) {
	parameter, err := parameter.New(snsEvent)
	if err != nil {
		log.Fatal(err)
	}

	itr, err := ghinstallation.New(http.DefaultTransport, parameter.GitHubIntegrationID, parameter.GitHubInstallationID, []byte(parameter.GitHubAppPrivateKey))
	if err != nil {
		log.Fatal(err)
	}
	itr.BaseURL = parameter.GitHubBaseURL

	client, _ := github.NewEnterpriseClient(parameter.GitHubAPIBaseURL, parameter.GitHubAPIUploadsBaseURL, &http.Client{Transport: itr})

	ctx := context.Background()
	owner := parameter.RepositoryOwner
	repo := parameter.Repository
	labels := []string{}
	if parameter.IssueLabel != "" {
		labels = append(labels, parameter.IssueLabel)
	}
	for _, issue := range parameter.Issues {
		req := &github.IssueRequest{
			Title:  github.String(issue.Subject),
			Body:   github.String(issue.Body),
			Labels: &labels,
		}
		_, _, err = client.Issues.Create(ctx, owner, repo, req)
		if err != nil {
			log.Fatal(err)
		}
	}
}
