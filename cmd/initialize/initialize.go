package initialize

import (
	"amartha/config"
	"amartha/internal/pkg/db"
	repoLoan "amartha/repository/loan"
	repoNotification "amartha/repository/notification"
	repoUpload "amartha/repository/upload"
	ucLoan "amartha/usecase/loan"
	"context"
	"log"
)

type App struct {
	Loan *ucLoan.Usecase
}

func Initialize(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Println("Initialize App")

	log.Println("Initialize DB")
	db, err := db.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialize Loan Repository")
	repoLoan, err := repoLoan.New(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialize Upload Repository")
	repoUpload, err := repoUpload.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialize Notification Repository")
	repoNotification, err := repoNotification.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialize Usecase")
	uc, err := ucLoan.New(ctx, cfg, repoLoan, repoUpload, repoNotification)
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		Loan: uc,
	}, nil

}
