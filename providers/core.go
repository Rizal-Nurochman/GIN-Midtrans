package providers

import (
	"log"

	"github.com/Rizal-Nurochman/config"
	"github.com/Rizal-Nurochman/database"

	authHandler "github.com/Rizal-Nurochman/modules/auth/handler"
	authRepo "github.com/Rizal-Nurochman/modules/auth/repository"
	authService "github.com/Rizal-Nurochman/modules/auth/service"

	"github.com/Rizal-Nurochman/modules/user/repository"
	"github.com/Rizal-Nurochman/pkg/constants"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (authService.JWTService, error) {
		return authService.NewJWTService(), nil
	})

	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	err := database.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	jwtService := do.MustInvokeNamed[authService.JWTService](injector, constants.JWTService)

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := authRepo.NewRefreshTokenRepository(db)

	authService := authService.NewAuthService(userRepository, refreshTokenRepository, jwtService, db)

	do.Provide(
		injector, func(i *do.Injector) (authHandler.AuthHandler, error) {
			return authHandler.NewAuthHandler(i, authService), nil
		},
	)
}