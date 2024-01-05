package database

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	postgres *gorm.DB
}

func GetDBConnection() (DBConnection, error) {
	var dbConn DBConnection
	if err := dbConn.ConnectToDB(); err != nil {
		return DBConnection{}, err
	}
	return dbConn, nil
}

func (c *DBConnection) ConnectToDB() error {
	db, err := gorm.Open(postgres.Open("host=solar_db user=raskob password=raskob dbname=postgres port=5432"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{SlowThreshold: time.Second}),
	})
	if err != nil {
		return err
	}
	c.postgres = db
	return nil
}

func (c *DBConnection) GetNonDailyData(ctx context.Context, startTime uint32, endTime uint32) ([]MinuteData, error) {
	var data dbMinuteDataSet

	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Where("total_e IS NULL AND dc1_u IS NOT NULL").
		Find(&data).Error
	return data.toExternal(), err
}

func (c *DBConnection) GetDailyData(ctx context.Context, startTime uint32, endTime uint32) ([]DailyData, error) {
	var data []DailyData

	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("total_e IS NOT NULL").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Find(&data).Error
	return data, err
}

func (c *DBConnection) GetZappiData(ctx context.Context, begin time.Time, end time.Time) ([]ZappiData, error) {
	var data []ZappiData

	err := c.postgres.
		WithContext(ctx).
		Order("plugged_in ASC").
		Where("plugged_in BETWEEN ? AND ? OR unplugged BETWEEN ? AND ?", begin, end, begin, end).
		Find(&data).Error
	return data, err
}
