package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/network"
	"avan/internal/repositories"
	"time"
)

type AnalysisService struct {
	AnalysisRepository repositories.AnalysisRepository
	GossipService      network.GossipService

	MessageService MessageService
}

func (a AnalysisService) GetNodeAnalysis() (int, int, int) {
	return a.GossipService.GetNodeHealth()
}

func (a AnalysisService) GetMessageSummary(workspaceID int64) (sqlc.MessageSummaryMysqlRow, error) {
	return a.MessageService.GetMessageSummary(workspaceID)
}

func (a AnalysisService) GetProvidersStats(workspaceID, ownerID int64) (map[string]sqlc.GetProviderStatsMysqlRow, error) {
	return a.MessageService.GetProvidersStats(workspaceID, ownerID)
}

func (a AnalysisService) GetTimeDistribution(workspaceID, ownerID int64, startDate, endDate *time.Time) (map[string]sqlc.GetTimeDistributionMysqlRow, error) {
	return a.MessageService.GetTimeDistribution(workspaceID, ownerID, startDate, endDate)
}

func NewAnalysisService(analysisRepository repositories.AnalysisRepository, gossipService network.GossipService, messageService MessageService) AnalysisService {
	return AnalysisService{AnalysisRepository: analysisRepository, GossipService: gossipService, MessageService: messageService}
}
