package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Constantes para os caminhos, facilitando a manutenção
const (
	migrationsPath = "./store/migrations"
	configPath     = "./store/migrations/tern.conf"
)

var rootCmd = &cobra.Command{
	Use:   "dbtool",
	Short: "Uma ferramenta de CLI para gerenciar migrações de banco de dados com tern.",
	Long: `dbtool é uma aplicação de linha de comando que serve como um wrapper
para o 'tern', facilitando a execução de migrações, rollbacks e verificação de status.`,
}

// Execute adiciona todos os comandos filhos ao comando raiz e define os sinalizadores apropriadamente.
// É chamado por main.main(). Só precisa acontecer uma vez para o rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Erros devem ser enviados para a saída de erro padrão (stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize é o local ideal para carregar configurações.
	// Será executado antes da função Run de qualquer comando.
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(createMigration)
	rootCmd.AddCommand(generate)
}

// initConfig carrega as variáveis de ambiente do arquivo .env.
// O godotenv.Load não retorna erro se o arquivo .env não existir.
func initConfig() {
	godotenv.Load()
}

func executeTernCommand(action string, args ...string) {
	cmdArgs := []string{
		action,
		"--migrations",
		migrationsPath,
		"--config",
		configPath,
	}
	cmdArgs = append(cmdArgs, args...)

	// Usar strings.Join para uma saída mais limpa, parecida com um comando real.
	fmt.Printf("Executando: tern %s\n", strings.Join(cmdArgs, " "))

	cmd := exec.Command("tern", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao executar o comando tern:", err)
		fmt.Fprintln(os.Stderr, "Saída do comando:", string(output))
		os.Exit(1)
	}

	fmt.Println("Comando executado com sucesso!")
	fmt.Println(string(output))
}
