package get_github_contributions

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Usecase struct {
}

func NewUsecase() *Usecase {
	return &Usecase{}
}

type GitHubContributionResponse struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks GitHubContributionWeeks `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

type GitHubContributionWeeks []struct {
	ContributionDays []GitHubContribution `json:"contributionDays"`
}

type GitHubContribution struct {
	Date              string `json:"date"`
	Color             string `json:"color"`
	ContributionCount int    `json:"contributionCount"`
}

func (u *Usecase) Run(
	ctx context.Context, githubToken string, username string, fromDateUTC string, toDateUTC string,
) (GitHubContributionWeeks, error) {
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
		&oauth2.Token{AccessToken: githubToken, TokenType: "Bearer"},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(apiUrl.String(), httpClient)

	fromTime, err := time.Parse(time.RFC3339, fromDateUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	toTime, err := time.Parse(time.RFC3339, toDateUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	type DateTime struct{ time.Time }
	variables := map[string]interface{}{
		"login": graphql.String(username),
		"from":  DateTime{fromTime},
		"to":    DateTime{toTime},
	}
	if err := client.Query(context.Background(), &query, variables); err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	return query.User.ContributionCollection.ContributionCalendar.Weeks, nil
}
