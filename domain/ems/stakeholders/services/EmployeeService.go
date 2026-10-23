package services

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/api/responses"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	stakeholderRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/repositories"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrEmployeeNameRequired   = errors.New("first name and second name are required")
	ErrEmployeeGenderRequired = errors.New("gender is required")
	ErrEmployeeGenderInvalid  = errors.New("invalid gender")
	ErrEmployeeEmailRequired  = errors.New("email is required")
	ErrEmployeeEmailInvalid   = errors.New("invalid email")
	ErrEmployeePhoneInvalid   = errors.New("invalid phone number")
	ErrEmployeeEmailExists    = errors.New("an account already exists for this email")
	ErrEmployeePhoneExists    = errors.New("an account already exists for this phone number")
	ErrEmployeeTitleInvalid   = errors.New("title must be 100 characters or fewer")
)

var employeeEmailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// CreateEmployee creates a staff identity (with a server-generated user number),
// its contacts, organisation membership and the initial stakeholder record in a
// single transaction. No system role or permission is granted.
//
// The activation email is best-effort: a mail failure never rolls back the
// employee record, and InvitationSent reports the outcome.
func CreateEmployee(actorNumber, organisationNumber string, request requests.CreateEmployeeRequest) (*responses.EmployeeResponse, error) {
	firstName := strings.TrimSpace(request.FirstName)
	secondName := strings.TrimSpace(request.SecondName)
	lastName := strings.TrimSpace(request.LastName)
	gender := strings.TrimSpace(request.Gender)
	email := strings.ToLower(strings.TrimSpace(request.Email))
	phone := strings.TrimSpace(request.Phone)
	title := strings.TrimSpace(request.Title)

	if firstName == "" || secondName == "" {
		return nil, ErrEmployeeNameRequired
	}
	if gender == "" {
		return nil, ErrEmployeeGenderRequired
	}
	if utf8.RuneCountInString(gender) > 8 {
		return nil, ErrEmployeeGenderInvalid
	}
	if email == "" {
		return nil, ErrEmployeeEmailRequired
	}
	if !employeeEmailPattern.MatchString(email) || utf8.RuneCountInString(email) > 40 {
		return nil, ErrEmployeeEmailInvalid
	}
	if utf8.RuneCountInString(title) > 100 {
		return nil, ErrEmployeeTitleInvalid
	}
	if firstNameLength := utf8.RuneCountInString(firstName); firstNameLength > 60 {
		return nil, ErrEmployeeNameRequired
	}
	if secondNameLength := utf8.RuneCountInString(secondName); secondNameLength > 60 {
		return nil, ErrEmployeeNameRequired
	}
	if lastNameLength := utf8.RuneCountInString(lastName); lastNameLength > 60 {
		return nil, ErrEmployeeNameRequired
	}
	if phone != "" {
		phone = commons.CheckValidPhone(phone)
		if phone == "" {
			return nil, ErrEmployeePhoneInvalid
		}
	}
	stakeholderType, typeErr := stakeholderModel.NormalizeStakeholderType(request.StakeholderType)
	if typeErr != nil {
		return nil, ErrInvalidStakeholderType
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}

	userNumber := uuid.NewV4().String()
	invitationToken := uuid.NewV4().String()

	user := authModel.User{
		Number:     userNumber,
		FirstName:  firstName,
		SecondName: secondName,
		LastName:   lastName,
		Gender:     gender,
		Type:       constants.TypeStaff,
		Status:     "Active",
	}

	stakeholder := stakeholderModel.Stakeholder{
		Number:             uuid.NewV4().String(),
		OrganisationNumber: organisationNumber,
		UserNumber:         userNumber,
		StakeholderType:    stakeholderType,
		Title:              title,
		ValidFrom:          request.ValidFrom,
		ValidTo:            request.ValidTo,
		Status:             constants.StatusActive,
		AddedByNumber:      actorNumber,
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		var contact authModel.Contact
		if err := tx.Where("LOWER(contact) = ? AND LOWER(contact_type) = ?", email, "email").First(&contact).Error; err == nil {
			return ErrEmployeeEmailExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if phone != "" {
			var phoneContact authModel.Contact
			if err := tx.Where("LOWER(contact) = ? AND LOWER(contact_type) = ?", strings.ToLower(phone), "phone").First(&phoneContact).Error; err == nil {
				return ErrEmployeePhoneExists
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&authModel.Contact{
			Contact:     email,
			ContactType: "Email",
			UserNumber:  userNumber,
		}).Error; err != nil {
			return err
		}

		if phone != "" {
			if err := tx.Create(&authModel.Contact{
				Contact:     phone,
				ContactType: "Phone",
				UserNumber:  userNumber,
			}).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&eosModel.EntityUser{
			EntityNumber:  organisationNumber,
			UserNumber:    userNumber,
			AddedByNumber: actorNumber,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&authModel.PasswordReset{
			Number:     uuid.NewV4().String(),
			UserNumber: userNumber,
			Token:      invitationToken,
			Status:     "Inactive",
		}).Error; err != nil {
			return err
		}

		existing, err := stakeholderRepositories.FindActiveByOrgUserType(organisationNumber, userNumber, stakeholderType, tx)
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrDuplicateActiveStakeholder
		}

		return stakeholderRepositories.Create(&stakeholder, tx)
	}); err != nil {
		return nil, err
	}

	invitationSent := emailSender.Info{
		Email:       email,
		Token:       invitationToken,
		MailType:    "WholesomeEd: User Account Confirmation",
		UserName:    secondName,
		CurrentYear: time.Now().Year(),
		Type:        "registration",
	}.SendEmail() == nil

	response := responses.EmployeeResponse{
		UserNumber:     userNumber,
		InvitationSent: invitationSent,
	}
	if stakeholderResponse, err := responseFor(stakeholder); err == nil {
		response.Stakeholder = *stakeholderResponse
	} else {
		response.Stakeholder = toResponse(stakeholder, false, strings.TrimSpace(strings.Join([]string{firstName, secondName, lastName}, " ")))
	}

	return &response, nil
}
