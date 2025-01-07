package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"
	version "github.com/optimism-java/interopbackend/migration"
	v0 "github.com/optimism-java/interopbackend/migration/version/v0"
	"log"

	"gorm.io/gorm"
)

var migrationOptions = gormigrate.Options{}

func InitMigrate(db *gorm.DB) {
	m := gormigrate.New(db, &migrationOptions, []*gormigrate.Migration{})

	m.InitSchema(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(v0.ModelSchemaList...)
		if err != nil {
			return err
		}
		return nil
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Init Migration successfully")
}

func upgradeLatestMigrate(db *gorm.DB) {
	m := gormigrate.New(db, &migrationOptions, version.ModelSchemaList)
	if len(version.ModelSchemaList) == 0 {
		return
	}

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Upgrade Latest Migration successfully")
}

func Migrate(db *gorm.DB) {
	InitMigrate(db)

	upgradeLatestMigrate(db)
}
