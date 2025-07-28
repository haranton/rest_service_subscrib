package handlers

import (
	"net/http"
	subscriptionService "rest_service/internal/subscriptionService"
	"time"

	"github.com/gin-gonic/gin"
)

type SubscriptionHadler struct {
	service subscriptionService.SubscriptionService
}

func NewSubscriptionHadler(s subscriptionService.SubscriptionService) *SubscriptionHadler {
	return &SubscriptionHadler{service: s}
}

func (h *SubscriptionHadler) ListSubscriptions(c *gin.Context) {

	subs, err := h.service.ListSubscriptions()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить список подписок"})
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHadler) CreateSubscription(c *gin.Context) {

	var req subscriptionService.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	sub, err := h.service.CreateSubscriptions(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHadler) GetSubscriptionByID(c *gin.Context) {

	idstr := c.Param("id")

	sub, err := h.service.GetSubscriptionByID(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)

}

func (h *SubscriptionHadler) UpdateSubscriptionByID(c *gin.Context) {

	var req subscriptionService.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	idstr := c.Param("id")
	updatedSub, err := h.service.UpdateSubcriptionByID(req, idstr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedSub)
}

func (h *SubscriptionHadler) DeleteSubcriptionByID(c *gin.Context) {

	idstr := c.Param("id")

	err := h.service.DeleteSubcriptionByID(idstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}

	c.JSON(http.StatusNoContent, "")
}

func (h *SubscriptionHadler) GetAmountOfsubscriptions(c *gin.Context) {

	params := subscriptionService.RequestParametersСalculatingSum{
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("name_service"),
	}

	total, err := h.service.GetAmountOfsubscriptions(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": total,
	})

}

func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())
	return yearDiff*12 + monthDiff + 1 // +1, чтобы включить начальный месяц
}
