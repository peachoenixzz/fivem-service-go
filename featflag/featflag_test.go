//go:build unit

package featflag

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/stretchr/testify/assert"
)

func TestListFeatureFlag(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.Config
		want string
	}{
		{"list all feature flag with default value",
			config.Config{},
			`{"isLimitMaxBalanceOnCreate":false}`},
		{"list all feature flag with isLimitMaxBalanceOnCreate = true",
			config.Config{FeatureFlag: config.FeatureFlag{IsLimitMaxBalanceOnCreate: true}},
			`{"isLimitMaxBalanceOnCreate":true}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := New(tc.cfg)

			// Assertions
			if assert.NoError(t, h.List(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.JSONEq(t, tc.want, rec.Body.String())
			}
		})
	}
}
