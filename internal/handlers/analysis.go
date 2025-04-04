package handlers

import (
	"avan/internal/services"
	"avan/views/analysis"
	"avan/views/layout"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AnalysisHandler struct {
	AnalysisService services.AnalysisService
}

func (w AnalysisHandler) Index(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate *time.Time
	if startDateStr != "" {
		t, _ := time.Parse(time.DateOnly, startDateStr)
		startDate = &t
	} else {
		t := time.Now().Add(1000000 * time.Hour)
		startDate = &t
	}

	var endDate *time.Time
	if endDateStr != "" {
		t, _ := time.Parse(time.DateOnly, endDateStr)
		endDate = &t
	} else {
		t := time.Now().Add(-1000000 * time.Hour)
		endDate = &t
	}

	nodes, healthy, unhealthy := w.AnalysisService.GetNodeAnalysis()
	databaseType := w.AnalysisService.AnalysisRepository.Base.GetDBDialector()
	consumerCount := w.AnalysisService.AnalysisRepository.Base.Config.KafkaConsumerCount

	ownerID := GetUserIDFromContext(c)
	workspaceID := GetWorkspaceIDFromContext(c)
	messageSummary, err := w.AnalysisService.GetMessageSummary(workspaceID)
	providersStats, err := w.AnalysisService.GetProvidersStats(workspaceID, ownerID)
	timeDistribution, err := w.AnalysisService.GetTimeDistribution(workspaceID, ownerID, startDate, endDate)

	c.JSON(200, timeDistribution)
	return
	data := analysis.AnalysisData{
		TotalMessages:   messageSummary.TotalMessages,
		SuccessMessages: messageSummary.SuccessMessages,
		FailedMessages:  messageSummary.FailedMessages,
		SuccessRate:     messageSummary.SuccessRate,
		ErrorRate:       messageSummary.ErrorRate,
		AvgSendTime:     messageSummary.AvgSendTime,

		DatabaseType: databaseType,

		ConsumerCount:    consumerCount,
		ProvidersStats:   providersStats,
		TimeDistribution: timeDistribution,
		NodesCount:       nodes,
		HealthyNodes:     healthy,
		UnhealthyNodes:   unhealthy,
	}

	defaultContent := analysis.Index(data)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func NewAnalysisHandler(analysisService services.AnalysisService) AnalysisHandler {
	return AnalysisHandler{AnalysisService: analysisService}
}
