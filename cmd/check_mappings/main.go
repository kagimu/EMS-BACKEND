package main

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
)

func main() {
	
	var mappings []models.QBCustomerMapping
	config.DB.Find(&mappings)
	
	fmt.Printf("Total mappings: %d\n", len(mappings))
	for _, m := range mappings {
		fmt.Printf("- User: %s, Profile: %s, Name: %s\n", m.UserNumber, m.QBProfileNumber, m.QbFullName)
	}
}
