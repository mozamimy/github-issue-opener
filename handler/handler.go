package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/mozamimy/lambda-github-issue-opener/parameter"
	"github.com/mozamimy/lambda-github-issue-opener/snsevent"
)

func HandleRequest(snsEvent snsevent.SNSEvent) {
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

	req := &github.IssueRequest{
		Title: github.String(parameter.IssueSubject),
		Body:  github.String(parameter.IssueBody),
	}

	ctx := context.Background()
	owner := parameter.RepositoryOwner
	repo := parameter.Repository
	_, _, err = client.Issues.Create(ctx, owner, repo, req)
	if err != nil {
		log.Fatal(err)
	}
}
