package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"pv-service/entities/dao"
)

type DBConnection interface {
	ConnectToDB() error
	GetAllData() *[]dao.PVData
	GetNonDailyDataBetweenStartAndEndTime(startTime uint32, endTime uint32) (*[]dao.PVData, error)
	GetDailyDataBetweenStartAndEndTime(startTime uint32, endTime uint32) (*[]dao.PVData, error)
}

type dbConnection struct {
	postgres *gorm.DB
}

func GetDBConnection() DBConnection {
	dbConn := new(dbConnection)
	err := dbConn.ConnectToDB()
	if err != nil {
		return nil
	}
	return dbConn
}

func (c *dbConnection) ConnectToDB() error {
	db, err := gorm.Open(postgres.Open("host=192.168.2.115 user=raskob password=raskob dbname=postgres port=5432"))
	if err != nil {
		return err
	}
	c.postgres = db
	return nil
}

func (c dbConnection) GetAllData() *[]dao.PVData {
	var (
		data []dao.PVData
	)
	err := c.postgres.Find(&data).Error
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return &data
}

func (c *dbConnection) GetNonDailyDataBetweenStartAndEndTime(startTime uint32, endTime uint32) (*[]dao.PVData, error) {
	var (
		data []dao.PVData
	)
	err := c.postgres.
		Order("time ASC").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Where("total_e IS NULL").
		Find(&data).Error
	return &data, err
}

func (c *dbConnection) GetDailyDataBetweenStartAndEndTime(startTime uint32, endTime uint32) (*[]dao.PVData, error) {
	var (
		data []dao.PVData
	)
	err := c.postgres.
		Order("time ASC").
		Where("time BETWEEN ? AND ?", startTime, endTime).
		Where("total_e IS NOT NULL").
		Find(&data).Error
	return &data, err
}
