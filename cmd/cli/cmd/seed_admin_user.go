package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/infrastracture/services"
	"github.com/spf13/cobra"
)

var seedAdminUserCmd = &cobra.Command{
	Use:   "add-admin",
	Short: "Add admin user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("failed to create config: %v", err)
		}
		db, err := infrastracture.NewDBMySQL(ctx, cfg)
		if err != nil {
			fmt.Printf("failed to create db: %v", err)
			os.Exit(1)
		}
		c := clocker.RealClocker{}
		userRepo, err := repository.NewUserRepository(&c)
		if err != nil {
			fmt.Printf("failed to create user repository: %v", err)
			os.Exit(1)
		}
		adminService, err := services.NewAdminService(db, userRepo)
		if err != nil {
			fmt.Printf("failed to create admin service: %v", err)
			os.Exit(1)
		}
		adminService.SeedAdminUser(ctx, cfg)
	},
}

func init() {
	rootCmd.AddCommand(seedAdminUserCmd)
}

type authService interface {
	SeedAdminUser(cfg *config.Config) (*models.User, error)
}

type AuthController struct {
	service authService
}

func (ac *AuthController) addAdminUser() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	admin, err := ac.service.SeedAdminUser(cfg)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	_ = admin
	return nil
}
