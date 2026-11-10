// Command owner-bootstrap creates the normal EMS user account that will act as
// the software provider/developer. It never creates an organisation, licence or
// customer profile and it never prints the password.
//
// Run from the backend directory with the normal .env/database configuration:
//
//	go run ./cmd/owner-bootstrap \
//	  -email owner@vendor.example \
//	  -password 'choose-a-strong-password' \
//	  -first-name Owner -second-name Developer -last-name Account -gender Male
//
// Then add the printed user number to EMS_SOFTWARE_OWNER_USER_NUMBERS and
// restart the backend. The first customer registration is not used as the
// software-provider identity.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/hassPassTokenGen"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found; using process environment variables")
	}

	email := flag.String("email", os.Getenv("EMS_OWNER_BOOTSTRAP_EMAIL"), "software provider email")
	password := flag.String("password", os.Getenv("EMS_OWNER_BOOTSTRAP_PASSWORD"), "software provider password")
	firstName := flag.String("first-name", "Software", "first name")
	secondName := flag.String("second-name", "Provider", "second name")
	lastName := flag.String("last-name", "Account", "last name")
	gender := flag.String("gender", "Other", "gender")
	phone := flag.String("phone", "", "optional phone number")
	flag.Parse()

	if err := run(*email, *password, *firstName, *secondName, *lastName, *gender, *phone); err != nil {
		log.Fatal(err)
	}
}

func run(email, password, firstName, secondName, lastName, gender, phone string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	firstName = strings.TrimSpace(firstName)
	secondName = strings.TrimSpace(secondName)
	lastName = strings.TrimSpace(lastName)
	gender = strings.TrimSpace(gender)
	phone = strings.TrimSpace(phone)

	if email == "" || password == "" {
		return errors.New("email and password are required (flags or EMS_OWNER_BOOTSTRAP_EMAIL/EMS_OWNER_BOOTSTRAP_PASSWORD)")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}
	if firstName == "" || secondName == "" {
		return errors.New("first name and second name are required")
	}
	if gender == "" {
		return errors.New("gender is required")
	}
	if phone != "" {
		phone = commons.CheckValidPhone(phone)
		if phone == "" {
			return errors.New("invalid phone number")
		}
	}

	db, err := gorm.Open(mysql.Open(config.DbURL(config.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	config.DB = db

	var existingUser authModels.User
	var existingContact authModels.Contact
	if err := db.Where("LOWER(contact) = ? AND LOWER(contact_type) = ?", email, "email").First(&existingContact).Error; err == nil {
		if err := db.Where("number = ?", existingContact.UserNumber).First(&existingUser).Error; err == nil {
			printInstruction(existingUser.Number)
			return nil
		}
	}
	if phone != "" {
		var count int64
		db.Model(&authModels.Contact{}).Where("LOWER(contact) = ? AND LOWER(contact_type) = ?", strings.ToLower(phone), "phone").Count(&count)
		if count > 0 {
			return errors.New("phone number already belongs to another account")
		}
	}

	hashed, err := hassPassTokenGen.Hash(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := authModels.User{
		Number: uuid.NewV4().String(), FirstName: firstName, SecondName: secondName, LastName: lastName,
		Gender: gender, Type: "User", Status: "Active",
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&authModels.Contact{Contact: email, ContactType: "Email", UserNumber: user.Number}).Error; err != nil {
			return err
		}
		if phone != "" {
			if err := tx.Create(&authModels.Contact{Contact: phone, ContactType: "Phone", UserNumber: user.Number}).Error; err != nil {
				return err
			}
		}
		return tx.Create(&authModels.UserPassword{UserNumber: user.Number, UserPassword: hashed}).Error
	})
	if err != nil {
		return fmt.Errorf("create software provider account: %w", err)
	}

	printInstruction(user.Number)
	return nil
}

func printInstruction(userNumber string) {
	fmt.Println("Software provider account is ready.")
	fmt.Println("User number:", userNumber)
	fmt.Println("Add this number to EMS_SOFTWARE_OWNER_USER_NUMBERS and restart the backend.")
	fmt.Println("The password is not displayed. Keep it in your password manager.")
}
