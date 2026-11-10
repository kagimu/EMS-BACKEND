package services

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	commercialRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/requests"
	commercialResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/responses"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrProfileNotFound     = errors.New("customer profile not found")
	ErrProfileExists       = errors.New("a commercial profile already exists for this organisation")
	ErrProfileOrganisation = errors.New("organisation not found")
	ErrProfileValidation   = errors.New("invalid customer profile")
	ErrProfileManager      = errors.New("relationship manager user not found")
)

var profileEmailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func toProfileResponse(profile commercialModels.CustomerProfile) commercialResponses.CustomerProfileResponse {
	return commercialResponses.CustomerProfileResponse{
		Number: profile.Number, OrganisationNumber: profile.OrganisationNumber, CustomerSince: profile.CustomerSince,
		CommercialStatus: profile.CommercialStatus, CustomerCategory: profile.CustomerCategory,
		RegisteredBusinessName: profile.RegisteredBusinessName, TaxIdentification: profile.TaxIdentification,
		BillingContactName: profile.BillingContactName, BillingEmail: profile.BillingEmail,
		BillingPhone: profile.BillingPhone, PhysicalAddress: profile.PhysicalAddress, Country: profile.Country,
		RelationshipManagerNumber: profile.RelationshipManagerNumber, InternalNotes: profile.InternalNotes,
		CreatedByNumber: profile.CreatedByNumber, UpdatedByNumber: profile.UpdatedByNumber,
		CreatedAt: profile.CreatedAt, UpdatedAt: profile.UpdatedAt,
	}
}

func encodeProfileState(profile commercialModels.CustomerProfile) string {
	encoded, _ := json.Marshal(toProfileResponse(profile))
	return string(encoded)
}

func recordCommercialEvent(tx *gorm.DB, profile commercialModels.CustomerProfile, actor, action, before, after, reason string) error {
	return tx.Create(&commercialModels.CommercialEvent{
		Number: uuid.NewV4().String(), EntityType: "CUSTOMER_PROFILE", EntityNumber: profile.Number,
		OrganisationNumber: profile.OrganisationNumber, Action: action, BeforeState: before,
		AfterState: after, Reason: strings.TrimSpace(reason), ActorUserNumber: actor,
	}).Error
}

func validateProfileInput(organisationNumber string, request commercialRequests.CustomerProfileRequest) (string, string, error) {
	var entityCount int64
	if err := config.DB.Model(&eosModel.Entity{}).Where("number = ?", organisationNumber).Count(&entityCount).Error; err != nil {
		return "", "", err
	}
	if entityCount == 0 {
		return "", "", ErrProfileOrganisation
	}
	status, err := commercialModels.NormaliseCommercialStatus(request.CommercialStatus)
	if err != nil {
		return "", "", err
	}
	category, err := commercialModels.NormaliseCustomerCategory(request.CustomerCategory)
	if err != nil {
		return "", "", err
	}
	email := strings.ToLower(strings.TrimSpace(request.BillingEmail))
	if email == "" || !profileEmailPattern.MatchString(email) {
		return "", "", ErrProfileValidation
	}
	if manager := strings.TrimSpace(request.RelationshipManagerNumber); manager != "" {
		var count int64
		if err := config.DB.Model(&authModel.User{}).Where("number = ?", manager).Count(&count).Error; err != nil {
			return "", "", err
		}
		if count == 0 {
			return "", "", ErrProfileManager
		}
	}
	return status, category, nil
}

func applyProfileInput(profile *commercialModels.CustomerProfile, request commercialRequests.CustomerProfileRequest, status, category string) {
	profile.CommercialStatus = status
	profile.CustomerCategory = category
	profile.RegisteredBusinessName = strings.TrimSpace(request.RegisteredBusinessName)
	profile.TaxIdentification = strings.TrimSpace(request.TaxIdentification)
	profile.BillingContactName = strings.TrimSpace(request.BillingContactName)
	profile.BillingEmail = strings.ToLower(strings.TrimSpace(request.BillingEmail))
	profile.BillingPhone = strings.TrimSpace(request.BillingPhone)
	profile.PhysicalAddress = strings.TrimSpace(request.PhysicalAddress)
	profile.Country = strings.TrimSpace(request.Country)
	profile.RelationshipManagerNumber = strings.TrimSpace(request.RelationshipManagerNumber)
	profile.InternalNotes = strings.TrimSpace(request.InternalNotes)
	if request.CustomerSince != nil {
		profile.CustomerSince = request.CustomerSince
	}
}

func GetProfile(organisationNumber string) (*commercialResponses.CustomerProfileResponse, error) {
	var profile commercialModels.CustomerProfile
	if err := config.DB.Where("organisation_number = ?", organisationNumber).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	response := toProfileResponse(profile)
	return &response, nil
}

func CreateProfile(actorNumber, organisationNumber string, request commercialRequests.CustomerProfileRequest) (*commercialResponses.CustomerProfileResponse, error) {
	status, category, err := validateProfileInput(organisationNumber, request)
	if err != nil {
		return nil, err
	}
	var existing int64
	config.DB.Model(&commercialModels.CustomerProfile{}).Where("organisation_number = ?", organisationNumber).Count(&existing)
	if existing > 0 {
		return nil, ErrProfileExists
	}
	profile := commercialModels.CustomerProfile{
		Number: uuid.NewV4().String(), OrganisationNumber: organisationNumber,
		CustomerSince: request.CustomerSince, CreatedByNumber: actorNumber, UpdatedByNumber: actorNumber,
	}
	if profile.CustomerSince == nil {
		now := time.Now()
		profile.CustomerSince = &now
	}
	applyProfileInput(&profile, request, status, category)

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}
		after := encodeProfileState(profile)
		return recordCommercialEvent(tx, profile, actorNumber, "CUSTOMER_PROFILE_CREATED", "", after, request.Reason)
	})
	if err != nil {
		return nil, err
	}
	response := toProfileResponse(profile)
	return &response, nil
}

func UpdateProfile(actorNumber, organisationNumber string, request commercialRequests.CustomerProfileRequest) (*commercialResponses.CustomerProfileResponse, error) {
	status, category, err := validateProfileInput(organisationNumber, request)
	if err != nil {
		return nil, err
	}
	var profile commercialModels.CustomerProfile
	if err := config.DB.Where("organisation_number = ?", organisationNumber).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	beforeProfile := profile
	before := encodeProfileState(beforeProfile)
	statusChanged := profile.CommercialStatus != status
	appliedSince := profile.CustomerSince
	applyProfileInput(&profile, request, status, category)
	if request.CustomerSince == nil {
		profile.CustomerSince = appliedSince
	}
	profile.UpdatedByNumber = actorNumber

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&profile).Error; err != nil {
			return err
		}
		after := encodeProfileState(profile)
		action := "CUSTOMER_PROFILE_UPDATED"
		if statusChanged {
			action = "CUSTOMER_STATUS_CHANGED"
		}
		return recordCommercialEvent(tx, profile, actorNumber, action, before, after, request.Reason)
	})
	if err != nil {
		return nil, err
	}
	response := toProfileResponse(profile)
	return &response, nil
}
