package service

import (
	"encoding/json"

	"github.com/roshankumar0036singh/auth-server/internal/dto"
	"github.com/roshankumar0036singh/auth-server/internal/models"
	"github.com/roshankumar0036singh/auth-server/internal/repository"
)

type AuditService struct {
	auditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

// LogEvent creates an audit log entry
func (s *AuditService) LogEvent(userID *string, action, entity, entityID, ip, userAgent string, metadata map[string]interface{}) error {
	metadataJSON := "{}"
	if metadata != nil {
		bytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(bytes)
		}
	}

	log := &models.AuditLog{
		UserID:    userID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		IPAddress: ip,
		UserAgent: userAgent,
		Metadata:  metadataJSON,
	}

	return s.auditRepo.Create(log)
}

// GetUserAuditLogs retrieves the audit logs for a specific user
func (s *AuditService) GetUserAuditLogs(userID string, page, limit int) (*dto.AuditLogsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	logs, err := s.auditRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.auditRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	hasMore := page*limit < int(totalCount)

	return &dto.AuditLogsResponse{
		Logs: logs,
		MetaData: dto.PaginationMetaData{
			TotalCount:  totalCount,
			CurrentPage: page,
			HasMore:     hasMore,
		},
	}, nil
}
