package db

import (
	"fmt"
	"time"

	"github.com/kajikaji0725/Examination/api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

type Config struct {
	Host     string
	Username string
	Password string
	DBname   string
	Port     string
}

func Dsn(config *Config) string {
	return fmt.Sprintf(
		"user=%s password=%s port=%s database=%s host=%s sslmode=disable",
		config.Username,
		config.Password,
		config.Port,
		config.DBname,
		config.Host,
	)
}

func NewController(config *Config) (*Controller, error) {

	db, err := gorm.Open(postgres.Open(Dsn(config)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.LocationDetail{},
	)
	if err != nil {
		return nil, err
	}

	return &Controller{db}, nil
}

func (controller *Controller) FetchDBRequestCount(from string, to string, timeFormat string) ([]model.LocationDetail, error) {

	postalCode := []model.LocationDetail{}
	var err error
	var toTime, fromTime time.Time
	switch {
	case from == "" && to == "":
		err = controller.db.Find(&postalCode).Error
	case from == "":
		toTime, err = time.Parse(timeFormat, to)
		fmt.Println(toTime)
		if err != nil {
			return nil, err
		}
		err = controller.db.Where("date <= ?", toTime).Find(&postalCode).Error
	case to == "":
		fromTime, err = time.Parse(timeFormat, from)
		if err != nil {
			return nil, err
		}
		err = controller.db.Where("date >= ?", fromTime).Find(&postalCode).Error
	default:
		fmt.Println("hoge")
		fromTime, err = time.Parse(timeFormat, from)
		if err != nil {
			return nil, err
		}
		toTime, err = time.Parse(timeFormat, to)
		if err != nil {
			return nil, err
		}
		err = controller.db.Where("date BETWEEN ? AND ?", from, to).Find(&postalCode).Error
	}

	if err != nil {
		return nil, err
	}
	return postalCode, nil
}

func (controller *Controller) SetDBPostalCode(locationDetail *model.LocationDetail) error {
	postalCode := model.LocationDetail{}
	err := controller.db.Model(&postalCode).Create(map[string]interface{}{"postal": locationDetail.Postal, "date": locationDetail.Date}).Error
	if err != nil {
		return err
	}
	return nil
}
