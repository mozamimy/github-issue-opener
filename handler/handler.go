package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v39/github"
	"github.com/mozamimy/github-issue-opener/parameter"
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
	for _, issue := range parameter.Issues {
		req := &github.IssueRequest{
			Title:  github.String(issue.Subject),
			Body:   github.String(issue.Body),
			Labels: &parameter.IssueLabels,
		}
		_, _, err = client.Issues.Create(ctx, owner, repo, req)
		if err != nil {
			log.Fatal(err)
		}
	}
}
