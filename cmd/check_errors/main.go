package main

import (
	"fmt"
	"log"
	
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
)

func main() {
	config.ConnectDB() 
	
	var items []models.QBSyncQueue
	err := config.DB.Where("status = ?", "FAILED").Order("created_at desc").Limit(5).Find(&items).Error
	if err != nil {
		log.Fatal(err)
	}
	
	for _, item := range items {
		fmt.Printf("Action: %s\nError: %s\nPayload: %s\n---\n", item.ActionType, item.ErrorMessage, item.QbXmlPayload)
	}
}
