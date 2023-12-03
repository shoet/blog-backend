package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/services"
	"github.com/shoet/blog/store"
	"github.com/spf13/cobra"
)

var seedAdminUserCmd = &cobra.Command{
	Use:   "add-admin",
	Short: "Add admin user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("failed to create config: %v", err)
		}
		db, err := store.NewDBMySQL(ctx, cfg)
		if err != nil {
			fmt.Printf("failed to create db: %v", err)
			os.Exit(1)
		}
		c := clocker.RealClocker{}
		userRepo, err := store.NewUserRepository(&c)
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
	SeedAdminUser(ctx context.Context, cfg *config.Config) (*models.User, error)
}

type AuthController struct {
	service authService
}

func (ac *AuthController) addAdminUser() error {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	admin, err := ac.service.SeedAdminUser(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	_ = admin
	return nil
}
