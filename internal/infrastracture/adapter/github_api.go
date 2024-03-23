package adapter

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type GitHubV4APIClient struct {
	githubPersonalAccessToken string
}

func NewGitHubV4APIClient(
	githubPersonalAccessToken string,
) *GitHubV4APIClient {
	return &GitHubV4APIClient{
		githubPersonalAccessToken: githubPersonalAccessToken,
	}
}

type GitHubContributionWeeks []struct {
	ContributionDays []struct {
		Date              string `json:"date"`
		Color             string `json:"color"`
		ContributionCount int    `json:"contributionCount"`
	} `json:"contributionDays"`
}

func (g *GitHubV4APIClient) GetContributions(ctx context.Context, username string, fromDateUTC time.Time, toDateUTC time.Time) (GitHubContributionWeeks, error) {
	apiUrl, err := url.Parse("https://api.github.com/graphql")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	var query struct {
		User struct {
			ContributionCollection struct {
				ContributionCalendar struct {
					Weeks GitHubContributionWeeks `json:"weeks"`
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.githubPersonalAccessToken, TokenType: "Bearer"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(apiUrl.String(), httpClient)

	type DateTime struct{ time.Time }
	variables := map[string]interface{}{
		"login": graphql.String(username),
		"from":  DateTime{fromDateUTC},
		"to":    DateTime{toDateUTC},
	}
	if err := client.Query(context.Background(), &query, variables); err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return query.User.ContributionCollection.ContributionCalendar.Weeks, nil
}
