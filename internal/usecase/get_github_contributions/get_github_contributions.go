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
	ctx context.Context, username string, fromDateUTC time.Time, toDateUTC time.Time,
) (adapter.GitHubContributionWeeks, error) {
	contributions, err := u.githubAPIv4.GetContributions(ctx, username, fromDateUTC, toDateUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %v", err)
	}
	return contributions, nil
}
