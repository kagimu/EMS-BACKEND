package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	enrollmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var enrollmentRepository = enrollmentRepositories.NewEnrolmentRepository()

func CreateBillPaymentPlan(payload requests.BillPaymentPlanRequest, currentUserNumber string) error {
	startDate, installments, totalAmount, err := parsePlanInput(payload.StartDate, payload.Installments)
	if err != nil {
		return err
	}

	invoiceByNumber, err := invoiceRepositories.FindInvoiceByNumber(payload.PromissoryNoteNumber)
	if err != nil {
		return errors.New("failed to fetch promissory note: " + err.Error())
	}
	if invoiceByNumber.Number == "" {
		return errors.New("invalid promissory note number")
	}

	student, err := studentRepositories.FindStudentByNumber(payload.StudentNumber)
	if err != nil {
		return errors.New("failed to fetch student: " + err.Error())
	}
	if student.Number == "" {
		return errors.New("invalid student number")
	}

	billable, err := invoiceRepositories.GetBillableByInvoiceNumber(invoiceByNumber.Number)
	if err != nil {
		return errors.New("failed to fetch invoice billable: " + err.Error())
	}

	var studyPeriodNumber string
	var billableNumber *string
	if billable.Number == "" {
		enrollments, err := enrollmentRepository.AllEnrolments(student.Number, "", "", "", constants.StatusActive)
		if err != nil || len(enrollments) == 0 {
			return errors.New("invalid invoice billable: unable to determine study period for payment plan")
		}
		studyPeriodNumber = enrollments[0].StudyPeriodNumber
	} else {
		studyPeriodNumber = billable.StudyPeriodNumber
		billableNumber = stringPointer(billable.Number)
	}

	alreadyExists, err := invoiceRepositories.FindActiveBillPaymentPlanByStudentAndPromissoryNote(payload.StudentNumber, payload.PromissoryNoteNumber)
	if err != nil {
		return errors.New("checking active payment plan failed: " + err.Error())
	}
	if alreadyExists.Number != "" {
		return errors.New("an active payment plan already exists for this student and promissory note")
	}

	plan := models.BillPaymentPlan{
		Number:               uuid.NewV4().String(),
		PromissoryNoteNumber: stringPointer(payload.PromissoryNoteNumber),
		BillableNumber:       billableNumber,
		StudentNumber:        payload.StudentNumber,
		StudyPeriodNumber:    studyPeriodNumber,
		TotalAmount:          totalAmount,
		StartDate:            startDate,
		Notes:                payload.Notes,
		Status:               constants.StatusActive,
		UserNumber:           currentUserNumber,
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := invoiceRepositories.SaveBillPaymentPlan(plan, tx); err != nil {
			return errors.New("failed to create bill payment plan: " + err.Error())
		}

		for i := range installments {
			installments[i].BillPaymentPlanNumber = plan.Number
			if err := invoiceRepositories.SaveBillPaymentPlanInstallment(installments[i], tx); err != nil {
				return errors.New("failed to create installment: " + err.Error())
			}
		}

		return nil
	})
}

func GetBillPaymentPlans(planNumber, promissoryNoteNumber, studentNumber, status string) ([]responses.BillPaymentPlanResponse, error) {
	plans, err := invoiceRepositories.FindBillPaymentPlans(planNumber, promissoryNoteNumber, studentNumber, status)
	if err != nil {
		return nil, errors.New("failed to fetch bill payment plans: " + err.Error())
	}

	response := make([]responses.BillPaymentPlanResponse, 0)
	for _, plan := range plans {
		response = append(response, buildBillPaymentPlanResponse(plan))
	}

	return response, nil
}

func UpdateBillPaymentPlan(planNumber string, payload requests.UpdateBillPaymentPlanRequest, currentUserNumber string) error {
	plan, err := invoiceRepositories.FindBillPaymentPlanByNumber(planNumber)
	if err != nil {
		return errors.New("failed to fetch bill payment plan: " + err.Error())
	}
	if plan.Number == "" {
		return errors.New("invalid bill payment plan number")
	}

	startDate, installments, totalAmount, err := parsePlanInput(payload.StartDate, payload.Installments)
	if err != nil {
		return err
	}

	invoiceByNumber, err := invoiceRepositories.FindInvoiceByNumber(payload.PromissoryNoteNumber)
	if err != nil {
		return errors.New("failed to fetch promissory note: " + err.Error())
	}
	if invoiceByNumber.Number == "" {
		return errors.New("invalid promissory note number")
	}

	student, err := studentRepositories.FindStudentByNumber(payload.StudentNumber)
	if err != nil {
		return errors.New("failed to fetch student: " + err.Error())
	}
	if student.Number == "" {
		return errors.New("invalid student number")
	}

	billable, err := invoiceRepositories.GetBillableByInvoiceNumber(invoiceByNumber.Number)
	if err != nil {
		return errors.New("failed to fetch invoice billable: " + err.Error())
	}

	var studyPeriodNumber string
	var billableNumber *string
	if billable.Number == "" {
		enrollments, err := enrollmentRepository.AllEnrolments(student.Number, "", "", "", constants.StatusActive)
		if err != nil || len(enrollments) == 0 {
			return errors.New("invalid invoice billable: unable to determine study period for payment plan")
		}
		studyPeriodNumber = enrollments[0].StudyPeriodNumber
	} else {
		studyPeriodNumber = billable.StudyPeriodNumber
		billableNumber = stringPointer(billable.Number)
	}

	plan.PromissoryNoteNumber = stringPointer(payload.PromissoryNoteNumber)
	plan.BillableNumber = billableNumber
	plan.StudentNumber = payload.StudentNumber
	plan.StudyPeriodNumber = studyPeriodNumber
	plan.StartDate = startDate
	plan.Notes = payload.Notes
	plan.TotalAmount = totalAmount
	plan.UserNumber = currentUserNumber
	if payload.Status != "" {
		plan.Status = payload.Status
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := invoiceRepositories.SaveBillPaymentPlan(plan, tx); err != nil {
			return errors.New("failed to update bill payment plan: " + err.Error())
		}

		if err := invoiceRepositories.DeleteBillPaymentPlanInstallmentsByPlanNumber(plan.Number, tx); err != nil {
			return errors.New("failed to clear old installments: " + err.Error())
		}

		for i := range installments {
			installments[i].BillPaymentPlanNumber = plan.Number
			if err := invoiceRepositories.SaveBillPaymentPlanInstallment(installments[i], tx); err != nil {
				return errors.New("failed to save installments: " + err.Error())
			}
		}

		return nil
	})
}

func PayBillPaymentPlanInstallment(payload requests.PayInstallmentRequest, currentUserNumber string) error {
	installment, err := invoiceRepositories.FindBillPaymentPlanInstallmentByNumber(payload.InstallmentNumber)
	if err != nil {
		return errors.New("failed to fetch installment: " + err.Error())
	}
	if installment.Number == "" {
		return errors.New("invalid installment number")
	}

	plan, err := invoiceRepositories.FindBillPaymentPlanByNumber(installment.BillPaymentPlanNumber)
	if err != nil {
		return errors.New("failed to fetch payment plan: " + err.Error())
	}
	if plan.Number == "" {
		return errors.New("invalid payment plan number")
	}

	installment.PaidAmount += payload.PaidAmount
	if installment.PaidAmount >= installment.Amount {
		installment.Status = "PAID"
	} else if installment.PaidAmount > 0 {
		installment.Status = "PARTIALLY_PAID"
	}

	paidOn, err := time.Parse(time.DateOnly, payload.PaidOn)
	if err != nil {
		return errors.New("invalid payment date format")
	}

	installmentPayment := models.BillPaymentPlanInstallmentPayment{
		Number:            uuid.NewV4().String(),
		InstallmentNumber: installment.Number,
		PaidAmount:        payload.PaidAmount,
		PaidOn:            paidOn,
		UserNumber:        currentUserNumber,
	}

	plan.PaidAmount += payload.PaidAmount
	// Check if all installments are paid to mark the entire plan as PAID
	allPaid := true
	for _, inst := range plan.Installments {
		// We need to use the updated installment if it's the one we just modified
		currentPaidAmount := inst.PaidAmount
		if inst.Number == installment.Number {
			currentPaidAmount = installment.PaidAmount
		}
		if currentPaidAmount < inst.Amount {
			allPaid = false
			break
		}
	}

	if allPaid {
		plan.Status = "PAID"
	} else if plan.PaidAmount > 0 {
		plan.Status = "PARTIALLY_PAID"
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := invoiceRepositories.SaveBillPaymentPlanInstallmentPayment(installmentPayment, tx); err != nil {
			return errors.New("failed to save installment payment record: " + err.Error())
		}
		if err := invoiceRepositories.SaveBillPaymentPlanInstallment(installment, tx); err != nil {
			return errors.New("failed to save installment: " + err.Error())
		}
		if err := invoiceRepositories.SaveBillPaymentPlan(plan, tx); err != nil {
			return errors.New("failed to update payment plan: " + err.Error())
		}
		return nil
	})
}

func CancelBillPaymentPlanInstallmentPayment(paymentNumber string) error {
	payment, err := invoiceRepositories.FindBillPaymentPlanInstallmentPaymentByNumber(paymentNumber)
	if err != nil {
		return errors.New("failed to fetch payment: " + err.Error())
	}
	if payment.Number == "" {
		return errors.New("invalid payment number")
	}

	installment, err := invoiceRepositories.FindBillPaymentPlanInstallmentByNumber(payment.InstallmentNumber)
	if err != nil {
		return errors.New("failed to fetch installment: " + err.Error())
	}
	if installment.Number == "" {
		return errors.New("invalid installment number")
	}

	plan, err := invoiceRepositories.FindBillPaymentPlanByNumber(installment.BillPaymentPlanNumber)
	if err != nil {
		return errors.New("failed to fetch payment plan: " + err.Error())
	}
	if plan.Number == "" {
		return errors.New("invalid payment plan number")
	}

	// Revert amounts
	installment.PaidAmount -= payment.PaidAmount
	if installment.PaidAmount <= 0 {
		installment.PaidAmount = 0
		installment.Status = "PENDING"
	} else if installment.PaidAmount < installment.Amount {
		installment.Status = "PARTIALLY_PAID"
	}

	plan.PaidAmount -= payment.PaidAmount
	if plan.PaidAmount <= 0 {
		plan.PaidAmount = 0
	}

	// Re-evaluate plan status
	allPaid := true
	for _, inst := range plan.Installments {
		currentPaidAmount := inst.PaidAmount
		if inst.Number == installment.Number {
			currentPaidAmount = installment.PaidAmount
		}
		if currentPaidAmount < inst.Amount {
			allPaid = false
			break
		}
	}

	if allPaid {
		plan.Status = "PAID"
	} else {
		if plan.PaidAmount > 0 {
			plan.Status = "PARTIALLY_PAID"
		} else {
			plan.Status = "ACTIVE"
		}
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := invoiceRepositories.DeleteBillPaymentPlanInstallmentPaymentByNumber(paymentNumber, tx); err != nil {
			return errors.New("failed to delete payment record: " + err.Error())
		}
		if err := invoiceRepositories.SaveBillPaymentPlanInstallment(installment, tx); err != nil {
			return errors.New("failed to save installment: " + err.Error())
		}
		if err := invoiceRepositories.SaveBillPaymentPlan(plan, tx); err != nil {
			return errors.New("failed to update payment plan: " + err.Error())
		}
		return nil
	})
}

func parsePlanInput(startDate string, installmentRequests []requests.BillPaymentPlanInstallmentRequest) (time.Time, []models.BillPaymentPlanInstallment, float64, error) {
	parsedStartDate, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return time.Time{}, nil, 0, errors.New("invalid start date")
	}

	installments := make([]models.BillPaymentPlanInstallment, 0)
	var totalAmount float64
	var previousDueDate time.Time

	for idx, installmentRequest := range installmentRequests {
		dueDate, err := time.Parse(time.DateOnly, installmentRequest.DueDate)
		if err != nil {
			return time.Time{}, nil, 0, errors.New("invalid installment due date")
		}

		if idx > 0 && dueDate.Before(previousDueDate) {
			return time.Time{}, nil, 0, errors.New("installment due dates must be in ascending order")
		}

		installments = append(installments, models.BillPaymentPlanInstallment{
			Number:     uuid.NewV4().String(),
			DueDate:    dueDate,
			Amount:     installmentRequest.Amount,
			PaidAmount: 0,
			Status:     "PENDING",
			Sequence:   idx + 1,
		})
		totalAmount += installmentRequest.Amount
		previousDueDate = dueDate
	}

	return parsedStartDate, installments, totalAmount, nil
}

func buildBillPaymentPlanResponse(plan models.BillPaymentPlan) responses.BillPaymentPlanResponse {
	installments := make([]responses.BillPaymentPlanInstallmentResponse, 0)
	for _, installment := range plan.Installments {
		payments := make([]responses.BillPaymentPlanInstallmentPaymentResponse, 0)
		for _, pay := range installment.Payments {
			payments = append(payments, responses.BillPaymentPlanInstallmentPaymentResponse{
				Number:     pay.Number,
				PaidAmount: pay.PaidAmount,
				PaidOn:     pay.PaidOn.Format(time.DateOnly),
			})
		}

		installments = append(installments, responses.BillPaymentPlanInstallmentResponse{
			Number:     installment.Number,
			DueDate:    installment.DueDate.Format(time.DateOnly),
			Amount:     installment.Amount,
			PaidAmount: installment.PaidAmount,
			Status:     installment.Status,
			Sequence:   installment.Sequence,
			Payments:   payments,
		})
	}

	promissoryNoteNumber := ""
	if plan.PromissoryNoteNumber != nil {
		promissoryNoteNumber = *plan.PromissoryNoteNumber
	}

	return responses.BillPaymentPlanResponse{
		Number:               plan.Number,
		PromissoryNoteNumber: promissoryNoteNumber,
		StudentNumber:        plan.StudentNumber,
		StudyPeriodNumber:    plan.StudyPeriodNumber,
		TotalAmount:          plan.TotalAmount,
		PaidAmount:           plan.PaidAmount,
		StartDate:            plan.StartDate.Format(time.DateOnly),
		Notes:                plan.Notes,
		Status:               plan.Status,
		Installments:         installments,
	}
}

func stringPointer(value string) *string {
	return &value
}
