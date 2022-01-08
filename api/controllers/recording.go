package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/raysandeep/Estimator-App/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/raysandeep/Estimator-App/schemas"
)

func StartCall(c *fiber.Ctx) error {
	u := new(schemas.StartCall)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	uid := int(rand.Uint32())
	rec := &utils.Recorder{
		Channel: u.Channel,
		UID:     uint32(uid),
	}

	record, err := utils.FetchRecord(rec.Channel)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	err = rec.Acquire()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	err = rec.Start(u.Channel, nil)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	utils.PatchRecord(record.ID, map[string]interface{}{
		"rid":    rec.RID,
		"sid":    rec.SID,
		"uid":    strconv.Itoa(int(rec.UID)),
		"status": 1,
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
		"data": fiber.Map{
			"rid":     rec.RID,
			"sid":     rec.SID,
			"token":   rec.Token,
			"channel": rec.Channel,
			"uid":     rec.UID,
		},
	})
}

func StopCall(c *fiber.Ctx) error {
	u := new(schemas.StopCall)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}

	record, err := utils.FetchRecord(u.Channel)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	uid, err := strconv.Atoi(record.UID)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	rec := &utils.Recorder{
		Channel: u.Channel,
		UID:     uint32(uid),
		RID:     record.Rid,
		SID:     record.Sid,
	}
	err = rec.Stop()
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	utils.PatchRecord(record.ID, map[string]interface{}{
		"status": 2,
	})

	return c.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "successful",
	})
}

func CreateRTCToken(c *fiber.Ctx) error {
	channel := c.Params("channel")
	uid := int(0)
	rtcToken, err := utils.GetRtcToken(channel, uid)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":      200,
		"rtc_token": rtcToken,
		"uid":       uid,
	})
}

func Process(ctx *fiber.Ctx) error {
	u := new(schemas.StopCall)

	if err := ctx.BodyParser(u); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": "invalid json",
			"err": err.Error(),
		})
	}
	record, err := utils.FetchRecord(u.Channel)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"msg": http.StatusInternalServerError,
			"err": err.Error(),
		})
	}

	objects := utils.ListObjectInS3(u.Channel)

	m3u8File := ""
	mp4File := ""
	for _, i := range objects {
		if strings.HasSuffix(i, ".m3u8") {
			m3u8File = i
			break
		} else if strings.HasSuffix(i, ".mp4") {
			mp4File = i
			break
		}
	}

	utils.PatchRecord(record.ID, map[string]interface{}{
		"status":             3,
		"video_url_download": mp4File,
		"video_url_play":     m3u8File,
	})

	return ctx.JSON(fiber.Map{"m3u8File": m3u8File, "mp4File": mp4File})
}

func PlayVideo(ctx *fiber.Ctx) error {
	// ?file=m3u8File
	file := ctx.Params("file")
	id := ctx.Params("id")

	object, err := utils.GetObjectInS3(file, 12*time.Hour)
	if err != nil {
		log.Println("unable to fetch")
		return ctx.JSON(fiber.Map{
			"error": err,
		})
	}
	client := http.DefaultClient

	resp, err := client.Get(object)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err,
		})
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err,
		})
	}
	lines := strings.Split(string(bytes), "\n")

	for i, line := range lines {
		if strings.Contains(line, ".ts") {
			tsBatch, err := utils.GetObjectInS3(fmt.Sprintf("%s/%s", id, line), 12*time.Hour)
			if err != nil {
				log.Println("unable to fetch")
				return ctx.JSON(fiber.Map{
					"error": err,
				})
			}
			lines[i] = tsBatch
		}
	}
	output := strings.Join(lines, "\n")
	ctx.Context().SetContentType("application/x-mpegURL")
	return ctx.SendString(output)
}

func DownloadVideo(ctx *fiber.Ctx) error {
	file := ctx.Params("file")

	object, err := utils.GetObjectInS3(file, 12*time.Hour)
	if err != nil {
		log.Println("unable to fetch")
		return ctx.JSON(fiber.Map{
			"error": err,
		})
	}

	return ctx.JSON(fiber.Map{
		"url": object,
	})
}
