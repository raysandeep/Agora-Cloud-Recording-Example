package api

import (
	"fmt"
	"github.com/raysandeep/Agora-Cloud-Recording-Example/utils"
	"math/rand"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/raysandeep/Agora-Cloud-Recording-Example/schemas"
)

func startCall(c *fiber.Ctx) error {
	u := new(schemas.StartCall)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	rec := &utils.Recorder{
		Channel: u.Channel,
		UID:     u.Uid,
	}

	_, err := rec.Acquire()

	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	_, err = rec.Start()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
		"data": map[string]interface{}{
			"rid":     rec.RID,
			"sid":     rec.SID,
			"token":   rec.Token,
			"channel": rec.Channel,
			"uid":     rec.UID,
		},
	})
}

func stopCall(c *fiber.Ctx) error {
	u := new(schemas.StopCall)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	_, err := utils.Stop(u.Channel, u.Uid, u.Rid, u.Sid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
	})
}

func createRTCToken(c *fiber.Ctx) error {
	channel := c.Params("channel")
	uid := int(rand.Uint32())
	rtcToken, err := utils.GetRtcToken(channel, uid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtc_token": rtcToken,
	})
}

func createRTMToken(c *fiber.Ctx) error {
	uid := c.Params("uid")
	rtmToken, err := utils.GetRtmToken(fmt.Sprint(uid))
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtm_token": rtmToken,
	})
}

func createTokens(c *fiber.Ctx) error {
	channel := c.Params("channel")
	uid := int(rand.Uint32())
	rtcToken, err := utils.GetRtcToken(channel, uid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	rtmToken, err := utils.GetRtmToken(fmt.Sprint(uid))
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"code":      http.StatusOK,
		"rtc_token": rtcToken,
		"rtm_token": rtmToken,
	})
}

// MountRoutes mounts all routes declared here
func MountRoutes(app *fiber.App) {
	app.Post("/api/start/call", startCall)
	app.Post("/api/stop/call", stopCall)
	app.Get("/api/get/rtc/:channel", createRTCToken)
	app.Get("/api/get/rtm/:uid", createRTMToken)
	app.Get("/api/tokens/:channel", createTokens)
}
