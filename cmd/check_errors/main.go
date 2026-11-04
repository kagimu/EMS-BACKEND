package main

import (
	"fmt"
	"log"
	
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	var err error
	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	
	var items []models.QBSyncQueue
	err = config.DB.Where("status = ?", "FAILED").Order("created_at desc").Limit(5).Find(&items).Error
	if err != nil {
		log.Fatal(err)
	}
	
	for _, item := range items {
		fmt.Printf("Action: %s\nError: %s\nPayload: %s\n---\n", item.ActionType, item.ErrorMessage, item.QbXmlPayload)
	}
}
