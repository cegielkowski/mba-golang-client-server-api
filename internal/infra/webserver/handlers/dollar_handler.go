package handlers

import (
	"context"
	"encoding/json"
	"github.com/cegielkowski/mba-golang-client-server-api/internal/entity"
	"github.com/cegielkowski/mba-golang-client-server-api/internal/infra/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type DollarHandler struct {
	DollarDB database.DollarInterface
}

func NewDollarHandler(db database.DollarInterface) *DollarHandler {
	return &DollarHandler{
		DollarDB: db,
	}
}

func (h *DollarHandler) GetDollar(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqChan := make(chan *http.Request)
	errChan := make(chan error)
	go func() {
		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			errChan <- err
		}
		reqChan <- req
	}()

	var req *http.Request
	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Request timeout")
		return c.SendStatus(fiber.StatusRequestTimeout)
	case <-errChan:
		return c.SendStatus(fiber.StatusBadRequest)
	case req = <-reqChan:
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var result entity.DollarResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	go func() {
		err = h.DollarDB.Create(ctx, &entity.Dollar{
			ID:        uuid.New().String(),
			Value:     result.Currency.Value,
			CreatedAt: time.Now(),
		})
		errChan <- err
	}()

	select {
	case <-time.After(10 * time.Millisecond):
		log.Println("Database timeout")
		return c.SendStatus(fiber.StatusRequestTimeout)
	case err = <-errChan:
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"value": result.Currency.Value})
}
