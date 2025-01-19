package main

import (
	"cashflow/api"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Println("Did not load an .env file")
	}

	err = repository.Init(models.Supabase)
	if err != nil {
		panic("Could not initiate the app : " + err.Error())
	}
	defer repository.Close()

	repo := repository.Get()
	bs := domain.NewBalanceService(repo)
	ts := domain.NewTransactionService(repo, &bs)
	ps := domain.NewPlanService(repo)
	as := domain.NewAuthService(repo)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
		AuthService:        &as,
		PlanService:        &ps,
	}

	api.Run(app)
}
