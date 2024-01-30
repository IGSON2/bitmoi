package api

import (
	"bitmoi/backend/utilities"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ImdIntervalRequest struct {
	Mode         string `json:"mode" validate:"required,oneof=competition practice" query:"mode"`
	Identifier   string `json:"identifier" validate:"required" query:"identifier"`
	ReqInterval  string `json:"reqinterval" validate:"required,oneof=5m 15m 1h 4h 1d" query:"reqinterval"`
	MinTimestamp int64  `json:"min_timestamp" validate:"required,number" query:"min_timestamp"`
	MaxTimestamp int64  `json:"max_timestamp" validate:"required,number" query:"max_timestamp"`
}

func (s *Server) getImdInterval(c *fiber.Ctx) error {
	req := new(ImdIntervalRequest)

	err := c.QueryParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}

	if errs := utilities.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)

	err = json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	cdd, err := s.selectInterChart(info, req.ReqInterval, req.MinTimestamp, req.MaxTimestamp, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot select intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
	}

	if !info.IsPracticeMode() {
		cdd.encodeChart(info.PriceFactor, info.VolumeFactor, info.TimeFactor)
	}

	return c.Status(fiber.StatusOK).JSON(OnePairChart{OneChart: cdd})
}
