package database

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"pv-service/entities/dao"
	"time"
)

type DBConnection struct {
	postgres *gorm.DB
}

func GetDBConnection() *DBConnection {
	dbConn := new(DBConnection)
	if err := dbConn.ConnectToDB(); err != nil {
		return nil
	}
	return dbConn
}

func (c *DBConnection) ConnectToDB() error {
	db, err := gorm.Open(postgres.Open("host=192.168.2.115 user=raskob password=raskob dbname=postgres port=5432"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{SlowThreshold: time.Second}),
	})
	if err != nil {
		return err
	}
	c.postgres = db
	return nil
}

func (c *DBConnection) GetAllData(ctx context.Context) *[]dao.PVData {
	var (
		data []dao.PVData
	)
	err := c.postgres.WithContext(ctx).Find(&data).Error
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return &data
}

func (c *DBConnection) GetNonDailyDataBetweenStartAndEndTime(ctx context.Context, startTime uint32, endTime uint32) (*[]dao.PVData, error) {
	var (
		data []dao.PVData
	)
	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Where("total_e IS NULL AND dc1_u IS NOT NULL").
		Find(&data).Error
	return &data, err
}

func (c *DBConnection) GetDailyDataBetweenStartAndEndTime(ctx context.Context, startTime uint32, endTime uint32) (*[]dao.PVData, error) {
	var (
		data []dao.PVData
	)
	err := c.postgres.
		WithContext(ctx).
		Order("time ASC").
		Where("total_e IS NOT NULL").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Find(&data).Error
	return &data, err
}

func (c *DBConnection) GetZappiDataBetweenStartAndEnddate(ctx context.Context, begin *time.Time, end *time.Time) (*[]dao.ZappiData, error) {
	data := []dao.ZappiData{}
	err := c.postgres.
		WithContext(ctx).
		Order("plugged_in ASC").
		Where("plugged_in BETWEEN ? AND ? OR unplugged BETWEEN ? AND ?", begin, end, begin, end).
		Find(&data).Error
	return &data, err
}
