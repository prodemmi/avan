package handlers

import (
	"avan/internal/services"
	"avan/views/analysis"
	"avan/views/layout"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
		start, err = parseFilterDate(startDate)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid start date")
		}
		end, err = parseFilterDate(endDate)
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

func parseFilterDate(value string) (time.Time, error) {
	value = strings.NewReplacer("۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4", "۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9", "٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4", "٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9").Replace(strings.TrimSpace(value))
	if date, err := time.ParseInLocation(time.DateOnly, value, time.Local); err == nil {
		return date, nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == '/' || r == '-' || r == '.' })
	if len(parts) != 3 {
		return time.Time{}, errors.New("invalid date")
	}
	jy, err := strconv.Atoi(parts[0])
	if err != nil || jy < 1200 || jy > 1600 {
		return time.Time{}, errors.New("invalid Jalali year")
	}
	jm, err := strconv.Atoi(parts[1])
	if err != nil || jm < 1 || jm > 12 {
		return time.Time{}, errors.New("invalid Jalali month")
	}
	jd, err := strconv.Atoi(parts[2])
	if err != nil || jd < 1 || jd > 31 || (jm > 6 && jm < 12 && jd > 30) || (jm == 12 && jd > 30) {
		return time.Time{}, errors.New("invalid Jalali day")
	}
	gy, gm, gd := jalaliToGregorian(jy, jm, jd)
	if jm == 12 && jd == 30 {
		nextY, nextM, nextD := jalaliToGregorian(jy+1, 1, 1)
		if gy == nextY && gm == nextM && gd == nextD {
			return time.Time{}, errors.New("invalid Jalali leap day")
		}
	}
	return time.Date(gy, time.Month(gm), gd, 0, 0, 0, 0, time.Local), nil
}

func jalaliToGregorian(jy, jm, jd int) (gy, gm, gd int) {
	jy += 1595
	days := -355668 + 365*jy + (jy/33)*8 + ((jy%33 + 3) / 4) + jd
	if jm < 7 {
		days += (jm - 1) * 31
	} else {
		days += (jm-7)*30 + 186
	}
	gy = 400 * (days / 146097)
	days %= 146097
	if days > 36524 {
		gy += 100 * ((days - 1) / 36524)
		days = (days - 1) % 36524
		if days >= 365 {
			days++
		}
	}
	gy += 4 * (days / 1461)
	days %= 1461
	if days > 365 {
		gy += (days - 1) / 365
		days = (days - 1) % 365
	}
	gd = days + 1
	leap := gy%4 == 0 && gy%100 != 0 || gy%400 == 0
	monthDays := [...]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if leap {
		monthDays[1] = 29
	}
	for gm = 1; gm <= 12 && gd > monthDays[gm-1]; gm++ {
		gd -= monthDays[gm-1]
	}
	return gy, gm, gd
}

func NewAnalysisHandler(analysisService services.AnalysisService) AnalysisHandler {
	return AnalysisHandler{AnalysisService: analysisService}
}
