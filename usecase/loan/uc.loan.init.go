package loan

import (
	"amartha/config"
	"context"
)

type Usecase struct {
	ctx              context.Context
	cfg              *config.Config
	repoLoan         RepoLoanInterface
	repoUpload       RepoUploadInterface
	repoNotification RepoNotificationInterface
}

func New(
	ctx context.Context,
	cfg *config.Config,
	repoLoan RepoLoanInterface,
	repoUpload RepoUploadInterface,
	repoNotification RepoNotificationInterface,
) (*Usecase, error) {
	return &Usecase{
		ctx:              ctx,
		cfg:              cfg,
		repoLoan:         repoLoan,
		repoUpload:       repoUpload,
		repoNotification: repoNotification,
	}, nil
}
