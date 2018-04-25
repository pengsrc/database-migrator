package constants

// MigrationExample is the example migration file.
const MigrationExample = "" +
	"-- +migrate Up\n" +
	"-- SQL in section 'Up' is executed when this migration is applied.\n" +
	"CREATE TABLE `examples` (\n" +
	"  `id` bigint(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n" +
	"  PRIMARY KEY (`id`)\n" +
	")\n" +
	"  ENGINE = InnoDB\n" +
	"  AUTO_INCREMENT = 10000\n" +
	"  DEFAULT CHARSET = utf8;\n" +
	"\n" +
	"-- +migrate Down\n" +
	"-- SQL section 'Down' is executed when this migration is rolled back.\n" +
	"DROP TABLE `examples`;\n"
