package main

import (
	"github.com/spf13/cobra"
)

var createMigration = &cobra.Command{
	Use:     "create <nome_da_migration>",
	Aliases: []string{"c"},
	Short:   "Cria um novo arquivo de migração",
	Long:    `Cria um novo arquivo de migração SQL na pasta de migrações.`,
	Args:    cobra.ExactArgs(1), // Garante que exatamente um argumento seja passado
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		executeTernCommand("new", migrationName)
	},
}
