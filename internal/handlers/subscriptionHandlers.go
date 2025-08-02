package handlers

import (
	"log"
	"net/http"
	subscriptionService "rest_service/internal/subscriptionService"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHadler struct {
	service subscriptionService.SubscriptionService
}

func NewSubscriptionHadler(s subscriptionService.SubscriptionService) *SubscriptionHadler {
	return &SubscriptionHadler{service: s}
}

// ListSubscriptions godoc
// @Summary      Получить список подписок
// @Description  Возвращает список всех подписок
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}  subscriptionService.Subscription
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionHadler) ListSubscriptions(c *gin.Context) {
	log.Println("[ListSubscriptions] Вход в хендлер")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Println("Неверные параметры запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный номер страницы"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || page < 1 || limit > 100 {
		log.Println("Неверное количество элементов")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный количество элементов"})
		return
	}

	paginatedResponse, err := h.service.ListSubscriptions(page, limit)
	if err != nil {
		log.Printf("[ListSubscriptions] Ошибка получения подписок: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить список подписок"})
		return
	}

	log.Printf("[ListSubscriptions] Подписок получено: %d\n", len(paginatedResponse.Data))
	c.JSON(http.StatusOK, paginatedResponse)
}

// CreateSubscription godoc
// @Summary      Создать новую подписку
// @Description  Создает подписку с переданными параметрами
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      subscriptionService.RequestBody  true  "Данные подписки"
// @Success      200           {object}  subscriptionService.Subscription
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionHadler) CreateSubscription(c *gin.Context) {
	log.Println("[CreateSubscription] Вход в хендлер")

	var req subscriptionService.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateSubscription] Ошибка привязки JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	sub, err := h.service.CreateSubscriptions(req)
	if err != nil {
		log.Printf("[CreateSubscription] Ошибка создания подписки: %v\n", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[CreateSubscription] Подписка создана: %+v\n", sub)
	c.JSON(http.StatusOK, sub)
}

// GetSubscriptionByID godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает подписку по уникальному идентификатору
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки"
// @Success      200  {object}  subscriptionService.Subscription
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHadler) GetSubscriptionByID(c *gin.Context) {
	idstr := c.Param("id")
	log.Printf("[GetSubscriptionByID] Поиск подписки по ID: %s\n", idstr)

	sub, err := h.service.GetSubscriptionByID(idstr)
	if err != nil {
		log.Printf("[GetSubscriptionByID] Ошибка: %v\n", err)
		c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[GetSubscriptionByID] Подписка найдена: %+v\n", sub)
	c.JSON(http.StatusOK, sub)
}

// UpdateSubscriptionByID godoc
// @Summary      Обновить подписку по ID
// @Description  Обновляет данные подписки с указанным ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      string                           true  "ID подписки"
// @Param        subscription  body      subscriptionService.RequestBody  true  "Обновленные данные подписки"
// @Success      200           {object}  subscriptionService.Subscription
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHadler) UpdateSubscriptionByID(c *gin.Context) {
	log.Println("[UpdateSubscriptionByID] Вход в хендлер")

	var req subscriptionService.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateSubscriptionByID] Ошибка привязки JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	idstr := c.Param("id")
	updatedSub, err := h.service.UpdateSubcriptionByID(req, idstr)
	if err != nil {
		log.Printf("[UpdateSubscriptionByID] Ошибка обновления подписки ID=%s: %v\n", idstr, err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[UpdateSubscriptionByID] Подписка обновлена: %+v\n", updatedSub)
	c.JSON(http.StatusOK, updatedSub)
}

// DeleteSubcriptionByID godoc
// @Summary      Удалить подписку по ID
// @Description  Удаляет подписку с указанным ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "ID подписки"
// @Success      204  {string}  string  "No Content"
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHadler) DeleteSubcriptionByID(c *gin.Context) {
	idstr := c.Param("id")
	log.Printf("[DeleteSubcriptionByID] Удаление подписки ID=%s\n", idstr)

	err := h.service.DeleteSubcriptionByID(idstr)
	if err != nil {
		log.Printf("[DeleteSubcriptionByID] Ошибка удаления: %v\n", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[DeleteSubcriptionByID] Подписка удалена ID=%s\n", idstr)
	c.JSON(http.StatusNoContent, "")
}

// GetAmountOfsubscriptions godoc
// @Summary      Получить сумму подписок по фильтрам
// @Description  Возвращает общую сумму подписок за указанный период с учетом фильтров
// @Tags         subscriptions
// @Produce      json
// @Param        start_date    query     string  false  "Дата начала (YYYY-MM-DD)"
// @Param        end_date      query     string  false  "Дата окончания (YYYY-MM-DD)"
// @Param        user_id       query     string  false  "ID пользователя"
// @Param        name_service  query     string  false  "Название сервиса"
// @Success      200           {object} map[string]int
// @Failure      500           {object}  map[string]string
// @Router       /subscriptions/amountSubscriptions [get]
func (h *SubscriptionHadler) GetAmountOfsubscriptions(c *gin.Context) {
	log.Println("[GetAmountOfsubscriptions] Вход в хендлер")

	params := subscriptionService.RequestParametersСalculatingSum{
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("name_service"),
	}

	log.Printf("[GetAmountOfsubscriptions] Параметры: %+v\n", params)

	total, err := h.service.GetAmountOfsubscriptions(params)
	if err != nil {
		log.Printf("[GetAmountOfsubscriptions] Ошибка вычисления суммы: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[GetAmountOfsubscriptions] Сумма: %v\n", total)
	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
