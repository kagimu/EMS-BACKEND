package requests

type PreAdmissionInterviewTestReportRequest struct {
	AdmissionStatus                 string `json:"admission_status"`
	ApplicationNumber               string `json:"application_number"`
	ClassNumber                     string `json:"class_number"`
	AttentionNeeded                 string `json:"attention_needed"`
	ReasonForChangingAdmissionClass string `json:"reason_for_changing_admission_class"`
	ReasonForNotAdmittingStudent    string `json:"reason_for_not_admitting_student"`
}
