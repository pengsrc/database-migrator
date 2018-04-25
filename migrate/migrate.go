package migrate

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Import Driver for MySQL
	"github.com/rubenv/sql-migrate"
	"gopkg.in/gorp.v1"

	"github.com/pengsrc/go-shared/convert"
)

// SQLConfig is config for SQL connection.
type SQLConfig struct {
	Dialect string

	Address  string // Format like "host:port"
	Database string

	User     string
	Password string
}

// SQLConnection represents a database connection used for migration.
type SQLConnection struct {
	db      *sql.DB
	dialect string

	source migrate.MigrationSource
}

// CloseMySQLConnection closes the database connection.
func (c *SQLConnection) CloseMySQLConnection() error {
	return c.db.Close()
}

// LocalMigrationSource migrates from a local files.
type LocalMigrationSource struct {
	Dir string
}

// FindMigrations implements migrate.MigrationSource{} interface.
func (l LocalMigrationSource) FindMigrations() (m []*migrate.Migration, err error) {
	m = make([]*migrate.Migration, 0)

	sqlExtension := ".sql"

	infos, err := ioutil.ReadDir(l.Dir)
	if err != nil {
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if strings.HasSuffix(info.Name(), sqlExtension) {
			file, err := os.Open(filepath.Join(l.Dir, info.Name()))
			if err != nil {
				return m, err
			}

			migration, err := migrate.ParseMigration(
				info.Name()[:len(info.Name())-len(sqlExtension)], file,
			)
			if err != nil {
				file.Close()
				return nil, err
			}

			file.Close()
			m = append(m, migration)
		}
	}

	// Make sure migrations are sorted
	sort.Sort(ByID(m))
	return
}

// ByID sorts migrations by ID.
type ByID []*migrate.Migration

func (b ByID) Len() int           { return len(b) }
func (b ByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByID) Less(i, j int) bool { return b[i].Less(b[j]) }

// NewMySQLConnection creates a new database connection for MySQL.
func NewMySQLConnection(config *SQLConfig, migrationsDir string) (connection *SQLConnection, err error) {
	// Lowercase the SQL dialect to avoid case issue.
	config.Dialect = strings.ToLower(config.Dialect)

	db, err := sql.Open(config.Dialect, fmt.Sprintf(
		`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true`,
		config.User, config.Password, config.Address, config.Database,
	))
	if err != nil {
		return nil, err
	}

	connection = &SQLConnection{
		db:      db,
		dialect: config.Dialect,
		source: &LocalMigrationSource{
			Dir: migrationsDir,
		},
	}
	return
}

// Status prints the database schema status.
func Status(c *SQLConnection) (err error) {
	records, err := migrate.GetMigrationRecords(c.db, c.dialect)
	if err != nil {
		return
	}

	fmt.Println("Applied At               Migration")
	fmt.Println(strings.Repeat("=", 80))

	for _, migration := range records {
		appliedAt := convert.TimeToString(migration.AppliedAt, convert.ISO8601)
		fmt.Printf("%s  -  %s\n", appliedAt, migration.Id)
	}

	newRecords, _, err := migrate.PlanMigration(c.db, c.dialect, c.source, migrate.Up, 0)
	if err != nil {
		return
	}

	for _, migration := range newRecords {
		fmt.Printf("Pending               -  %s\n", migration.Id)
	}

	return
}

// Up executes a single database migration. It returns the last migration name.
func Up(c *SQLConnection) (record string, err error) {
	count, err := migrate.ExecMax(c.db, c.dialect, c.source, migrate.Up, 1)
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	records, err := migrate.GetMigrationRecords(c.db, c.dialect)
	if err != nil {
		return
	}

	record = records[len(records)-1].Id
	return
}

// Down reverts a single database migration. It returns the reverted migration name.
func Down(c *SQLConnection) (record string, err error) {
	records, err := migrate.GetMigrationRecords(c.db, c.dialect)
	if err != nil {
		return
	}

	count, err := migrate.ExecMax(c.db, c.dialect, c.source, migrate.Down, 1)
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	record = records[len(records)-1].Id
	return
}

// Sync migrates the database to latest schema.
func Sync(c *SQLConnection) (done []string, err error) {
	count, err := migrate.ExecMax(c.db, c.dialect, c.source, migrate.Up, 0)
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	records, err := migrate.GetMigrationRecords(c.db, c.dialect)
	if err != nil {
		return
	}

	for i := len(records) - count; i < len(records); i++ {
		done = append(done, records[i].Id)
	}
	return
}

func init() {
	migrate.SetTable("schema_migrations")
	migrate.MigrationDialects["mysql"] = gorp.MySQLDialect{
		Engine: "InnoDB", Encoding: "utf8mb4",
	}
}
