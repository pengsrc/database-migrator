package cmds

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pengsrc/database-migrator/constants"
	"github.com/pengsrc/database-migrator/migrate"
	"github.com/pengsrc/go-shared/check"
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once.
func Execute() {
	check.ErrorForExit(constants.Name, migratorCMD.Execute())
}

func init() {
	migratorCMD.SilenceErrors = true
	migratorCMD.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	migratorCMD.AddCommand(statusCMD)
	migratorCMD.AddCommand(syncCMD)
	migratorCMD.AddCommand(upCMD)
	migratorCMD.AddCommand(downCMD)
	migratorCMD.AddCommand(newCMD)

	addDatabaseFlags(statusCMD)
	addDatabaseFlags(syncCMD)
	addDatabaseFlags(upCMD)
	addDatabaseFlags(downCMD)
	addNewFlags(newCMD)
}

// migratorCMD represents the migrate sub-command.
var migratorCMD = &cobra.Command{
	Use:   constants.Name,
	Short: "Database schema migration tool",
	Long:  "Database schema migration tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// statusCMD gets the database schema status.
var statusCMD = &cobra.Command{
	Use:   "status <service>",
	Short: "Show database schema status",
	Long:  "Show database schema status",
	Run: func(cmd *cobra.Command, args []string) {
		connection := loadMySQLConnection()
		defer connection.CloseMySQLConnection()

		check.ErrorForExit(constants.Name, migrate.Status(connection))
	},
}

// syncCMD runs one migration.
var syncCMD = &cobra.Command{
	Use:   "sync",
	Short: "Sync database schema",
	Long:  "Sync database schema",
	Run: func(cmd *cobra.Command, args []string) {
		connection := loadMySQLConnection()
		defer connection.CloseMySQLConnection()

		migrations, err := migrate.Sync(connection)
		check.ErrorForExit(constants.Name, err)

		for i, m := range migrations {
			migrations[i] = fmt.Sprintf("- %s\n", m)
		}

		if len(migrations) > 0 {
			fmt.Printf("Migrated to the latest database schema.\n%s", strings.Join(migrations, ""))
		} else {
			fmt.Println("Already has the latest database schema.")
		}
	},
}

// upCMD runs one migration.
var upCMD = &cobra.Command{
	Use:   "up",
	Short: "Run one migration",
	Long:  "Run one migration",
	Run: func(cmd *cobra.Command, args []string) {
		connection := loadMySQLConnection()
		defer connection.CloseMySQLConnection()

		migration, err := migrate.Up(connection)
		check.ErrorForExit(constants.Name, err)

		if migration != "" {
			fmt.Printf("Run migration complete, %s.\n", migration)
		} else {
			fmt.Println("No migration exceuted.")
		}
	},
}

// downCMD reverts one migration.
var downCMD = &cobra.Command{
	Use:   "down",
	Short: "Revert one migration",
	Long:  "Revert one migration",
	Run: func(cmd *cobra.Command, args []string) {
		connection := loadMySQLConnection()
		defer connection.CloseMySQLConnection()

		migration, err := migrate.Down(connection)
		check.ErrorForExit(constants.Name, err)

		if migration != "" {
			fmt.Printf("Revert migration complete, %s.\n", migration)
		} else {
			fmt.Println("No migration exceuted.")
		}
	},
}

// newCMD create an new migration.
var newCMD = &cobra.Command{
	Use:   "new",
	Short: "Create an new migration",
	Long:  "Create an new migration",
	Run: func(cmd *cobra.Command, args []string) {
		if dbMigrationsDir == "" {
			fmt.Println("Please specify migration directory.")
			return
		}
		if dbMigrationName == "" {
			fmt.Println("Please specify migration name.")
			return
		}

		check.ErrorForExit(constants.Name, check.Dir(dbMigrationsDir))
		target := filepath.Join(dbMigrationsDir, fmt.Sprintf(
			"%s_%s.sql",
			time.Now().UTC().Format("20060102150405"), dbMigrationName,
		))

		check.ErrorForExit(
			constants.Name,
			ioutil.WriteFile(target, []byte(constants.MigrationExample), 0644),
		)
		fmt.Printf("New migration created at \"%s\".\n", target)
	},
}

// Flags for database migrations commands.
var (
	dbFlagDialect  string
	dbFlagHost     string
	dbFlagPort     int
	dbFlagDatabase string
	dbFlagUser     string
	dbFlagPassword string

	dbMigrationsDir string
	dbMigrationName string

	dbFlagHelp bool
)

// addDatabaseFlags add flags for command which need database connection.
func addDatabaseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&dbFlagDialect, "dialect", "", "MySQL", "Specify database dialect",
	)
	cmd.Flags().StringVarP(
		&dbFlagHost, "host", "h", "127.0.0.1", "Specify database host",
	)
	cmd.Flags().IntVarP(
		&dbFlagPort, "port", "p", 3306, "Specify database port",
	)
	cmd.Flags().StringVarP(
		&dbFlagDatabase, "database", "d", "", "Specify database name",
	)
	cmd.Flags().StringVarP(
		&dbFlagUser, "user", "u", "root", "Specify user",
	)
	cmd.Flags().StringVarP(
		&dbFlagPassword, "password", "P", "", "Specify password",
	)

	cmd.Flags().StringVarP(
		&dbMigrationsDir, "migrations", "m", "", "Specify migration directory",
	)

	// Manually add help flag to avoid conflict.
	cmd.Flags().BoolVarP(
		&dbFlagHelp, "help", "", false, "Show help",
	)

	cmd.MarkFlagRequired("dialect")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("database")
	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("migrations")
}

// addNewFlags add flags for new command.
func addNewFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&dbMigrationsDir, "migrations", "m", "", "Specify migration directory",
	)
	cmd.Flags().StringVarP(
		&dbMigrationName, "name", "n", "", "Specify migration name",
	)

	// Manually add help flag to avoid conflict.
	cmd.Flags().BoolVarP(
		&dbFlagHelp, "help", "", false, "Show help",
	)

	cmd.MarkFlagRequired("migrations")
	cmd.MarkFlagRequired("name")
}

func loadMySQLConnection() (connection *migrate.SQLConnection) {
	connection, err := migrate.NewMySQLConnection(&migrate.SQLConfig{
		Dialect:  dbFlagDialect,
		Address:  fmt.Sprintf("%s:%d", dbFlagHost, dbFlagPort),
		Database: dbFlagDatabase,
		User:     dbFlagUser,
		Password: dbFlagPassword,
	}, dbMigrationsDir)
	check.ErrorForExit(constants.Name, err)

	return
}
