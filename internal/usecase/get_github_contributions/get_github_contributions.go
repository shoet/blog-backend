package get_github_contributions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	query := `
	query($login:String!, $from:DateTime!, $to:DateTime!) {
	  user(login: $login) {
		contributionsCollection(from: $from, to: $to) {
			contributionCalendar {
				weeks {
					contributionDays {
						color
						contributionCount
						date
					}
				}
			}
		}
	  }
	}
	`

	type variables struct {
		Login string `json:"login"`
		From  string `json:"from"`
		To    string `json:"to"`
	}

	requestBody := struct {
		Query     string            `json:"query"`
		Variables map[string]string `json:"variables"`
	}{
		Query: query,
		Variables: map[string]string{
			"login": "shoet",
			"from":  "2024-03-01T15:00:00Z",
			"to":    "2024-03-22T14:59:59Z",
		},
	}

	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	request, err := http.NewRequest("POST", apiUrl.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	var githubResponse GitHubContributionResponse
	if err := json.NewDecoder(response.Body).Decode(&githubResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return githubResponse.Data.User.ContributionsCollection.ContributionCalendar.Weeks, nil
}
