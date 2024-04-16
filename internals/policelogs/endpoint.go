package policelogs

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	Status  int
	Message interface{}
}

type Response struct {
	ArrestPlayerName      string        `json:"arrest_player_name"`
	PolicePlayerName      string        `json:"police_player_name"`
	ArrestSteamPlayerName string        `json:"arrest_steam_player_name"`
	PoliceSteamPlayerName string        `json:"police_steam_player_name"`
	ArrestJobPlayer       string        `json:"arrest_job_player"`
	ArrestSexPlayer       string        `json:"arrest_sex_player"`
	Case                  []interface{} `json:"case"`
	CaseQuantity          []interface{} `json:"case_quantity"`
	CaseCustom            []interface{} `json:"case_custom"`
	TimeCustom            []interface{} `json:"time_custom"`
	FineCustom            []interface{} `json:"fine_custom"`
	AllMiliSec            int64         `json:"all_milisec"`
	AllMinute             int64         `json:"all_mins"`
	AllFine               int64         `json:"all_fine"`
	PoliceDecreaseTime    int64         `json:"police_decrease_time"`
}

type Request struct {
	ArrestPlayerName      string        `json:"arrest_player_name"`
	PolicePlayerName      string        `json:"police_player_name"`
	ArrestSteamPlayerName string        `json:"arrest_steam_player_name"`
	PoliceSteamPlayerName string        `json:"police_steam_player_name"`
	ArrestJobPlayer       string        `json:"arrest_job_player"`
	ArrestSexPlayer       string        `json:"arrest_sex_player"`
	Case                  []interface{} `json:"case"`
	CaseQuantity          []interface{} `json:"case_quantity"`
	CaseCustom            []interface{} `json:"case_custom"`
	TimeCustom            []interface{} `json:"time_custom"`
	FineCustom            []interface{} `json:"fine_custom"`
	AllMiliSec            int64         `json:"all_milisec"`
	AllMinute             int64         `json:"all_mins"`
	AllFine               int64         `json:"all_fine"`
	PoliceDecreaseTime    int64         `json:"police_decrease_time"`
}

type CaseDetail struct {
	Imprison interface{} `json:"imprison"`
	Fine     interface{} `json:"fine"`
	Count    interface{} `json:"count"`
}

func (h Handler) GetPoliceLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	res, err := h.PoliceLog()
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Database Error : ", err.Error())
	}

	logger.Info("get result successfully")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AddPoliceLogEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	var req Request
	err := c.Bind(&req)
	logger.Info("get request event endpoint successfully")
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	playerInfo := user.Claims.(*mw.JwtCustomClaims)
	fmt.Println("JOB", playerInfo.Job)
	fmt.Println("identifier", playerInfo.Identifier)
	fmt.Println("group", playerInfo.Group)

	var resultCase []string
	var resultCaseQuantity []string
	for _, elem := range req.CaseQuantity {
		switch v := elem.(type) {
		case string:
			resultCaseQuantity = append(resultCaseQuantity, v) // already a string, no conversion needed
		case int, float64, bool:
			// use fmt.Sprint to convert non-string types to string
			resultCaseQuantity = append(resultCaseQuantity, fmt.Sprint(v))
		case nil:
			resultCaseQuantity = append(resultCaseQuantity, "") // or however you want to handle nil values
		default:
			fmt.Printf("Unhandled type for: %v\n", v)
		}
	}

	for _, elem := range req.Case {
		switch v := elem.(type) {
		case string:
			resultCase = append(resultCase, v) // already a string, no conversion needed
		case int, float64, bool:
			// use fmt.Sprint to convert non-string types to string
			resultCase = append(resultCase, fmt.Sprint(v))
		case nil:
			resultCase = append(resultCase, "") // or however you want to handle nil values
		default:
			fmt.Printf("Unhandled type for: %v\n", v)
		}
	}

	for _, v := range resultCaseQuantity {
		var casesQuantity map[string]map[string]interface{}
		err := json.Unmarshal([]byte(v), &casesQuantity)
		if err != nil {
			log.Fatalf("Error parsing JSON: %s", err)
		}

		caseNameTranslations := map[string]string{
			"illegal_money_case":        "คดีเงินผิดกฎหมาย",
			"illegal_cement_case":       "คดีปูนผิดกฎหมาย",
			"illegal_ice_case":          "คดีน้ำแข็งผิดกฎหมาย",
			"illegal_drug_skull":        "คดียาเสพติดหัวกะโหลก",
			"illegal_heroin_case":       "คดีฮีโรอีนผิดกฎหมาย",
			"illegal_n_amp_case":        "คดี N-AMP ผิดกฎหมาย",
			"illegal_amp_case":          "คดี AMP ผิดกฎหมาย",
			"illegal_screwdv":           "คดี screwdriver ผิดกฎหมาย",
			"illegal_storerobbery":      "คดีปล้นร้านค้า",
			"housebreaking_case":        "คดีงัดแงะบ้าน",
			"illegal_keycard":           "คดี keycard ผิดกฎหมาย",
			"illegal_cocaine_pack_case": "คดีโคเคนผิดกฎหมาย",
			"illegal_ice_pack_case":     "คดีแพ็คน้ำแข็งผิดกฎหมาย",
			"illegal_cocaine_case":      "คดีโคเคนผิดกฎหมาย",
			"decrease_imprison":         "ลดเวลาจำคุก",
		}

		var builderCaseQuantity strings.Builder
		for caseName, details := range casesQuantity {
			thaiCaseName, ok := caseNameTranslations[caseName]
			if !ok {
				fmt.Printf("No Thai translation for case: %s\n", caseName)
				continue
			}
			builderCaseQuantity.WriteString(fmt.Sprintf("คดี: %s, เวลาติดคุก: %v นาที , ค่าปรับ: %v$ จำนวน: %v ชิ้น\n", thaiCaseName, details["imprison"], details["fine"], details["count"]))
		}

		var builderCase strings.Builder
		for _, vrc := range resultCase {
			var cases map[string]interface{}
			err := json.Unmarshal([]byte(vrc), &cases)
			if err != nil {
				log.Fatalf("Error parsing JSON: %s", err)
			}

			builderCase.WriteString(fmt.Sprintf("คดี: %s, เวลาติดคุก: %v นาที , ค่าปรับ: %v$\n", cases["label_name"], cases["imprison"], cases["fine"]))
		}

		msg := fmt.Sprintf(
			"น้องเทียนเทียนแจ้งเตือน @everyone การจับกุมใหม่เกิดขึ้น! ```\n"+
				"ชื่อผู้เล่นที่ถูกจับ: %v\n"+
				"ชื่อตำรวจ: %v\n"+
				"คดีของผิดกฎหมาย: %v\n"+
				"คดีหลัก : %v\n"+
				"รายละเอียดเพิ่มเติม: %v\n"+
				"เวลา: %v\n"+
				"ค่าปรับ: %v\n"+
				"เวลาทั้งหมด: %d นาที\n"+
				"ค่าปรับทั้งหมด: %d\n"+
				"เวลาที่ลดได้โดยตำรวจ: %d นาที```",
			req.ArrestPlayerName,
			req.PolicePlayerName,
			builderCaseQuantity.String(),
			builderCase.String(),
			req.CaseCustom,
			req.TimeCustom,
			req.FineCustom,
			req.AllMinute,
			req.AllFine,
			(req.PoliceDecreaseTime/1000)/60,
		)

		_, err = h.Discord.ChannelMessageSend("1220316952153554987", msg)
		if err != nil {
			return fmt.Errorf("unable to send message to Discord: %v", err)
		}
	}
	var mes Message
	mes, err = h.InsertMLog(req)
	logger.Info("prepare data to create successfully")
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, mes)
	}
	logger.Info("create successfully")
	return c.JSON(http.StatusCreated, mes)
}

func (h Handler) CaseEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	event := c.Param("event")
	logger.Info("prepare log")
	res, err := h.LogCaseEventAndSteamID(steamID, event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event case and steam id endpoint")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) AllEventAndSteamIDEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	steamID := c.Param("steamid")
	res, err := h.LogAllEventAndSteamID(steamID)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event and steamid endpoint")
	return c.JSON(http.StatusOK, res)
}

func (h Handler) ByEventEndPoint(c echo.Context) error {
	logger := mlog.L(c)
	event := c.Param("event")
	res, err := h.LogCaseEventAll(event)
	if err != nil {
		logger.Error("Database Error : ", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, res)
	}

	logger.Info("get event endpoint successfully")
	return c.JSON(http.StatusOK, res)
}
