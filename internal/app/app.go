package app

import (
	"ykjam/new_tx/config"
	"ykjam/new_tx/pkg/logger"
	"ykjam/new_tx/pkg/postgres"
)

func Run(cfg *config.Config) {

	logger.SetLevel(cfg.Log.Level)
	log := logger.Get()

	pg, err := postgres.New(cfg.Postgres.Url)
	if err != nil {
		log.Fatal("app - Run - postgres.New: %s", err)
	}
	defer pg.Close()

	//rp := repo.New(pg)
	//userUseCase := use_case.NewUser(rp)
	//depUseCase := use_case.NewDepartment(rp)
	//posUseCase := use_case.NewPosition(rp)

}
