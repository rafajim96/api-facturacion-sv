package drivers

import (
	"fmt"
	"time"
	errPackage "github.com/MarlonG1/api-facturacion-sv/config/error"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"gorm.io/gorm"
)

type DriverConfig interface {
	GetDSN() gorm.Dialector
	GetDriverName() string
	GetHost() string
	GetStringConnection() string
}

type DbConnection struct {
	Db     *gorm.DB
	Config *gorm.Config
	Driver DriverConfig
	Err    error
}

func NewDatabaseConnection(driver DriverConfig) *DbConnection {
	return &DbConnection{
		Driver: driver,
		Config: &gorm.Config{},
	}
}

func (d *DbConnection) Open() error {
	const maxRetries = 10
	const retryDelay = 3 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		d.Db, d.Err = gorm.Open(d.Driver.GetDSN(), d.Config)
		if d.Err == nil {
			logs.Info("Database connection has been set successfully", map[string]interface{}{
				"Database type:": d.Driver.GetDriverName(),
				"Database host":  d.Driver.GetHost(),
			})
			return nil
		}

		if attempt < maxRetries {
			logs.Warn(fmt.Sprintf("Failed to connect to DB [attempt %d/%d]", attempt, maxRetries), map[string]interface{}{
				"Database type:":      d.Driver.GetDriverName(),
				"Database connection": d.Driver.GetStringConnection(),
				"Database error":      d.Err.Error(),
			})
			time.Sleep(retryDelay)
			continue
		}
	}

	logs.Error("Exceeded max retries to connect to database", map[string]interface{}{
		"Database type:":      d.Driver.GetDriverName(),
		"Database connection": d.Driver.GetStringConnection(),
		"Last error":          d.Err.Error(),
	})

	return d.Err
}

func (d *DbConnection) Close() error {
	dbInstance, err := d.Db.DB()
	if err != nil {
		logs.Error(errPackage.ErrFailedToGetDBInstance.Error(), map[string]interface{}{
			"Database type:":      d.Driver.GetDriverName(),
			"Database connection": d.Driver.GetStringConnection(),
			"Database error":      err.Error(),
		})
		return err
	}

	if err := dbInstance.Close(); err != nil {
		logs.Error(errPackage.ErrFailedToCloseDbConnection.Error(), map[string]interface{}{
			"Database type:":      d.Driver.GetDriverName(),
			"Database connection": d.Driver.GetStringConnection(),
			"Database error":      err.Error(),
		})
		return err
	}

	logs.Info(fmt.Sprintf("Database connection close successfully"), map[string]interface{}{
		"Database type:": d.Driver.GetDriverName(),
	})

	return nil
}
