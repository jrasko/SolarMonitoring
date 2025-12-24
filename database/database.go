package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	postgres *gorm.DB
}

func GetDBConnection(host, user, password string) (Connection, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432", host, user, password)

	db, err := gorm.Open(
		postgres.Open(dsn),
	)
	return Connection{postgres: db}, err
}

func (c *Connection) GetNonDailyData(ctx context.Context, startTime uint32, endTime uint32) ([]MinuteData, error) {
	var data dbMinuteDataSet

	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Where("total_e IS NULL AND dc1_u IS NOT NULL").
		Find(&data).Error
	return data.toExternal(), err
}

func (c *Connection) GetDailyData(ctx context.Context, startTime uint32, endTime uint32) ([]DailyData, error) {
	var data []DailyData

	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("total_e IS NOT NULL").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Find(&data).Error
	return data, err
}

func (c *Connection) GetZappiData(ctx context.Context, begin time.Time, end time.Time) ([]ZappiData, error) {
	var data []ZappiData

	err := c.postgres.
		WithContext(ctx).
		Order("plugged_in ASC").
		Where("plugged_in BETWEEN ? AND ? OR unplugged BETWEEN ? AND ?", begin, end, begin, end).
		Find(&data).Error
	return data, err
}
