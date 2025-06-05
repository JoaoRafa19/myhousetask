package migrations

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	dbDriver      = "mysql"
	dbUser        = "root"
	dbPass        = "root"
	dbHost        = "localhost"
	dbPort        = "3307"
	dbName        = "testegrpc"
	migrationsDir = "file://db/migrations"
)

func Run() (*sql.DB, error) {
	// Conecta ao MySQL sem um banco de dados específico para poder criar o nosso.
	initDb, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/")
	if err != nil {
		log.Fatalf("Falha ao conectar no MySQL: %v", err)
	}
	_, err = initDb.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatalf("Falha ao criar o banco de dados: %v", err)
	}
	initDb.Close()

	// Agora, conecta ao banco de dados que sabemos que existe.
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		return nil, err
	}

	// Executa as migrações
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not create mysql driver instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsDir,
		dbName,
		driver,
	)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("failed to get migration version: %v", err)
	}

	if dirty {
		log.Println("database is dirty, forcing to version", version-1)
		if err := m.Force(int(version - 1)); err != nil {
			log.Fatalf("failed to force migration version: %v", err)
		}
	}

	errm := m.Up(); 
	if errm != nil && errm != migrate.ErrNoChange {
		log.Fatalf("could not run migrations: %v", err)
	}

	if errm != nil && errm == migrate.ErrNoChange {
		log.Println("No changes to the database")
	}

	log.Println("Migrations ran successfully")
	return db, nil
}
