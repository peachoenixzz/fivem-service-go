package playerquest

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ResponseValidateItem struct {
	LimitType      string `json:"limit_type"`
	Name           string `json:"item_name"`
	Category       string `json:"category"`
	MaxLimit       int64  `json:"max_limit"`
	Point          int64  `json:"point"`
	RemainQuantity int64  `json:"remaining_quantity"`
	ExpireDateItem int    `json:"expire_date_item"`
}

func (h Handler) CreateQuestPlayer(c echo.Context) error {

	return c.JSON(http.StatusOK, res)
}
