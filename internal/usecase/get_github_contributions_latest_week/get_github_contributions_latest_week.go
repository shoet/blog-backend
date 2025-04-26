package get_github_contributions_latest_week

import (
	"context"
	"fmt"
	"time"

	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure/adapter"
)

type GitHubV4APIAdapter interface {
	GetContributions(ctx context.Context, username string, fromDateUTC time.Time, toDateUTC time.Time) (adapter.GitHubContributionWeeks, error)
}

type Usecase struct {
	githubAPIv4 GitHubV4APIAdapter
	clock       clocker.Clocker
}

func NewUsecase(
	githubAPIv4 GitHubV4APIAdapter,
	clock clocker.Clocker,
) *Usecase {
	return &Usecase{
		githubAPIv4: githubAPIv4,
		clock:       clock,
	}
}

func (u *Usecase) Run(
	ctx context.Context, username string, numOfLatestWeeks int,
) (adapter.GitHubContributionWeeks, error) {
	// todate: 次の土曜日
	// fromdate: 次の土曜日から直近numOfLatestWeeks週間前の翌日開始

	c := u.clock.Now()
	nextSat := (int(time.Saturday) - int(c.Weekday()) + 7) % 7
	toDateTime := c.AddDate(0, 0, nextSat)
	fromDateTime := toDateTime.AddDate(0, 0, -7*numOfLatestWeeks+1)

	contributions, err := u.githubAPIv4.GetContributions(ctx, username, fromDateTime, toDateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %v", err)
	}

	return contributions, nil
}
