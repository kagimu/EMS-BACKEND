package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_whole?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	
	var queues []models.QBSyncQueue
	db.Where("action_type = ?", "CUSTOMER_ADD").Order("id desc").Limit(5).Find(&queues)
	
	for _, q := range queues {
		fmt.Printf("ID: %d, Action: %s, Status: %s, User: %s, TargetProfile: %s, Err: %s\n", q.ID, q.ActionType, q.Status, q.UserNumber, q.QBProfileNumber, q.ErrorMessage)
	}

	var mappings []models.QBCustomerMapping
	db.Order("id desc").Limit(5).Find(&mappings)
	fmt.Println("\nRecent mappings:")
	for _, m := range mappings {
		fmt.Printf("ID: %d, User: %s, Profile: %s, FullName: %s\n", m.ID, m.UserNumber, m.QBProfileNumber, m.QbFullName)
	}
}
