package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/spf13/cobra"
)

// rollbackCmd representa o comando rollback
var rollbackCmd = &cobra.Command{
	Use:     "rollback [opções]",
	Aliases: []string{"r"},
	Short:   "Reverte a última migração (tern rollback)",
	Long: `Reverte a migração mais recente aplicada ao banco de dados.
Você pode passar argumentos adicionais que serão repassados para o tern.
Exemplo: dbtool rollback`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		)

		conn, err := pgx.Connect(ctx, connString)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not connect to database: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = conn.Close(ctx) }()

		m, err := migrate.NewMigrator(ctx, conn, "schema_version")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Não foi possível criar o migrador: %v\n", err)
			os.Exit(1)
		}

		migrationsFS := os.DirFS(migrationsPath)
		if err := m.LoadMigrations(migrationsFS); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Não foi possível carregar os arquivos de migração: %v\n", err)
			os.Exit(1)
		}

		version, err := m.GetCurrentVersion(ctx)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Não foi possível obter a versão atual: %v\n", err)
			os.Exit(1)
		}

		if version == 0 {
			fmt.Println("Nenhuma migração para reverter.")
			return
		}

		if err := m.MigrateTo(ctx, version-1); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Erro ao reverter a migração: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Rolled back to version %d\n", version-1)
	},
}
