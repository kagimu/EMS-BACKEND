package services

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrAttachmentInvalid  = errors.New("unsupported or invalid attachment")
	ErrAttachmentTooLarge = errors.New("attachment exceeds the 10 MB limit")
	ErrAttachmentNotFound = errors.New("attachment not found")
)

const maxSupportAttachmentBytes = 10 * 1024 * 1024

var attachmentTypes = map[string][]string{
	".png":  {"image/png"},
	".jpg":  {"image/jpeg"},
	".jpeg": {"image/jpeg"},
	".pdf":  {"application/pdf"},
	".txt":  {"text/plain"},
	".csv":  {"text/plain", "text/csv", "application/csv"},
}

func attachmentDir() string {
	dir := strings.TrimSpace(os.Getenv("EMS_SUPPORT_ATTACHMENT_DIR"))
	if dir == "" {
		dir = "./storage/support-attachments"
	}
	return dir
}

func toAttachmentResponse(attachment supportModels.Attachment) supportResponses.AttachmentResponse {
	return supportResponses.AttachmentResponse{
		Number: attachment.Number, TicketNumber: attachment.TicketNumber,
		OriginalName: attachment.OriginalName, ContentType: attachment.ContentType,
		SizeBytes: attachment.SizeBytes, Visibility: attachment.Visibility, CreatedAt: attachment.CreatedAt,
	}
}

func SaveSupportAttachment(ticket supportModels.Ticket, uploader, visibility string, header *multipart.FileHeader) (*supportResponses.AttachmentResponse, error) {
	visibility = strings.ToUpper(strings.TrimSpace(visibility))
	if visibility != supportModels.VisibilityCustomer && visibility != supportModels.VisibilityInternal {
		return nil, ErrAttachmentInvalid
	}
	if header == nil || header.Size <= 0 || header.Size > maxSupportAttachmentBytes {
		if header != nil && header.Size > maxSupportAttachmentBytes {
			return nil, ErrAttachmentTooLarge
		}
		return nil, ErrAttachmentInvalid
	}
	original := filepath.Base(strings.TrimSpace(header.Filename))
	extension := strings.ToLower(filepath.Ext(original))
	allowed, ok := attachmentTypes[extension]
	if !ok {
		return nil, ErrAttachmentInvalid
	}

	source, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()
	sniff := make([]byte, 512)
	read, _ := io.ReadFull(source, sniff)
	contentType := strings.TrimSpace(strings.Split(http.DetectContentType(sniff[:read]), ";")[0])
	validType := false
	for _, candidate := range allowed {
		if candidate == contentType {
			validType = true
			break
		}
	}
	if !validType {
		return nil, ErrAttachmentInvalid
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	dir := attachmentDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	storedName := uuid.NewV4().String() + extension
	path := filepath.Join(dir, storedName)
	target, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	written, copyErr := io.Copy(target, io.LimitReader(source, maxSupportAttachmentBytes+1))
	closeErr := target.Close()
	if copyErr != nil || closeErr != nil || written > maxSupportAttachmentBytes {
		os.Remove(path)
		if written > maxSupportAttachmentBytes {
			return nil, ErrAttachmentTooLarge
		}
		return nil, ErrAttachmentInvalid
	}

	attachment := supportModels.Attachment{
		Number: uuid.NewV4().String(), TicketNumber: ticket.Number, UploadedByNumber: uploader,
		OriginalName: original, StoredName: storedName, ContentType: contentType,
		SizeBytes: written, Visibility: visibility,
	}
	if err := config.DB.Create(&attachment).Error; err != nil {
		os.Remove(path)
		return nil, err
	}
	response := toAttachmentResponse(attachment)
	return &response, nil
}

func ListSupportAttachments(ticketNumber string, includeInternal bool) ([]supportResponses.AttachmentResponse, error) {
	query := config.DB.Where("ticket_number = ?", ticketNumber)
	if !includeInternal {
		query = query.Where("visibility = ?", supportModels.VisibilityCustomer)
	}
	var attachments []supportModels.Attachment
	if err := query.Order("created_at ASC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	result := make([]supportResponses.AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		result = append(result, toAttachmentResponse(attachment))
	}
	return result, nil
}

// AttachmentPath authorizes a download and returns the internal path. The path
// is never exposed in an API response.
func AttachmentPath(number, organisationFilter string, includeInternal bool) (supportModels.Attachment, string, error) {
	var attachment supportModels.Attachment
	if err := config.DB.Where("number = ?", number).First(&attachment).Error; err != nil {
		return attachment, "", ErrAttachmentNotFound
	}
	if !includeInternal && attachment.Visibility != supportModels.VisibilityCustomer {
		return attachment, "", ErrAttachmentNotFound
	}
	var ticket supportModels.Ticket
	if err := config.DB.Where("number = ?", attachment.TicketNumber).First(&ticket).Error; err != nil {
		return attachment, "", ErrAttachmentNotFound
	}
	if organisationFilter != "" && ticket.OrganisationNumber != organisationFilter {
		return attachment, "", ErrAttachmentNotFound
	}
	path := filepath.Join(attachmentDir(), filepath.Base(attachment.StoredName))
	if _, err := os.Stat(path); err != nil {
		return attachment, "", ErrAttachmentNotFound
	}
	return attachment, path, nil
}
