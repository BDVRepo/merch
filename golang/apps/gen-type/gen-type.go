package main

import (
	"bdv-avito-merch/libs/3_infrastructure/db_manager"
	"bdv-avito-merch/libs/4_common/env_vars"
	"bdv-avito-merch/libs/4_common/smart_context"
	"os"

	"gorm.io/gen"
)

type Querier interface {
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	env_vars.LoadEnvVars()
	os.Setenv("LOG_LEVEL", "info")
	logger := smart_context.NewSmartContext()

	dbm, err := db_manager.NewDbManager(logger)
	if err != nil {
		logger.Fatalf("NewDbManager failed: %v", err)
	}

	logger = logger.WithDB(dbm.GetGORM())

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./libs/2_generated_models/model",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(logger.GetDB())

	g.GenerateAllTable()
	g.ApplyBasic()

	g.Execute()
}
