package get_github_contributions

import (
	"context"
	"fmt"
	"time"

	"github.com/shoet/blog/internal/infrastracture/adapter"
)

type GitHubV4APIAdapter interface {
	GetContributions(ctx context.Context, username string, fromDateUTC time.Time, toDateUTC time.Time) (adapter.GitHubContributionWeeks, error)
}

type Usecase struct {
	githubAPIv4 GitHubV4APIAdapter
}

func NewUsecase(
	githubAPIv4 GitHubV4APIAdapter,
) *Usecase {
	return &Usecase{
		githubAPIv4: githubAPIv4,
	}
}

func (u *Usecase) Run(
	ctx context.Context, username string, fromDateUTCStr string, toDateUTCStr string,
) (adapter.GitHubContributionWeeks, error) {

	fromTime, err := time.Parse(time.RFC3339, fromDateUTCStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	toTime, err := time.Parse(time.RFC3339, toDateUTCStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	contributions, err := u.githubAPIv4.GetContributions(ctx, username, fromTime, toTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %v", err)
	}

	return contributions, nil
}
