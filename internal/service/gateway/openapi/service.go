package openapi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"downloader_test/internal/service/workers/orchestrator"
)

type HTTPHandler struct {
	publisher Publisher
}

func New(p Publisher) *HTTPHandler {
	return &HTTPHandler{publisher: p}
}

func (s *HTTPHandler) Add(c echo.Context) error {
	var req []string
	if err := c.Bind(&req); err != nil {
		return err
	}
	if s.publisher.Pub(req...) == orchestrator.ErrStopped {
		return c.JSON(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))
	}
	return c.JSON(http.StatusAccepted, http.StatusText(http.StatusAccepted))
}

func (s *HTTPHandler) Configure(g *echo.Group)  {
	g.POST("/loader", s.Add)
}