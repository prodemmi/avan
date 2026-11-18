package handlers

import (
	"avan/internal/services"
	"avan/views/analysis"
	"avan/views/layout"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AnalysisHandler struct {
	AnalysisService services.AnalysisService
}

func (w AnalysisHandler) Index(c *gin.Context) {
	rangeName := c.DefaultQuery("range", "last_7_days")
	start, end, err := analysisDateRange(rangeName, c.Query("start_date"), c.Query("end_date"), time.Now())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nodes, healthy, unhealthy := w.AnalysisService.GetNodeAnalysis()
	consumerCount := w.AnalysisService.AnalysisRepository.Base.Config.RedisConsumerCount

	ownerID := GetUserIDFromContext(c)
	workspaceID := GetWorkspaceIDFromContext(c)
	messageSummary, err := w.AnalysisService.GetMessageSummary(workspaceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	providersStats, err := w.AnalysisService.GetProvidersStats(workspaceID, ownerID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	timeDistribution, err := w.AnalysisService.GetTimeDistribution(workspaceID, ownerID, &start, &end)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := analysis.AnalysisData{
		TotalMessages:   messageSummary.TotalMessages,
		SuccessMessages: messageSummary.SuccessMessages,
		FailedMessages:  messageSummary.FailedMessages,
		SuccessRate:     messageSummary.SuccessRate,
		ErrorRate:       messageSummary.ErrorRate,
		AvgSendTime:     messageSummary.AvgSendTime,

		ConsumerCount:    consumerCount,
		ProvidersStats:   providersStats,
		TimeDistribution: timeDistribution,
		Range:            rangeName,
		StartDate:        start.Format(time.DateOnly),
		EndDate:          end.Format(time.DateOnly),
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

func analysisDateRange(rangeName, startDate, endDate string, now time.Time) (time.Time, time.Time, error) {
	var start, end time.Time
	switch rangeName {
	case "last_hour":
		start, end = now.Add(-time.Hour), now
	case "last_2_hours":
		start, end = now.Add(-2*time.Hour), now
	case "last_24_hours":
		start, end = now.Add(-24*time.Hour), now
	case "last_7_days":
		start, end = now.AddDate(0, 0, -7), now
	case "last_30_days":
		start, end = now.AddDate(0, 0, -30), now
	case "custom":
		var err error
		start, err = time.ParseInLocation(time.DateOnly, startDate, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid start date")
		}
		end, err = time.ParseInLocation(time.DateOnly, endDate, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid end date")
		}
		end = end.AddDate(0, 0, 1).Add(-time.Nanosecond)
	default:
		return time.Time{}, time.Time{}, errors.New("invalid date range")
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, errors.New("end date must be on or after start date")
	}
	return start, end, nil
}

func NewAnalysisHandler(analysisService services.AnalysisService) AnalysisHandler {
	return AnalysisHandler{AnalysisService: analysisService}
}
