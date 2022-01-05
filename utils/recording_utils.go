package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/raysandeep/Estimator-App/schemas"

	"github.com/spf13/viper"
)

type Recorder struct {
	http.Client
	Channel string
	Token   string
	UID     uint32
	RID     string
	SID     string
}

// Acquire runs the acquire endpoint for Cloud Recording
func (rec *Recorder) Acquire() error {
	creds, err := GenerateUserCredentials(rec.Channel)
	if err != nil {
		return err
	}

	rec.UID = uint32(creds.UID)
	rec.Token = creds.Rtc

	requestBody, _ := json.Marshal(&schemas.AcquireRequest{
		Cname: rec.Channel,
		UID:   strconv.Itoa(int(rec.UID)),
		ClientRequest: schemas.AcquireClientRequest{
			ResourceExpiredHour: 24,
		},
	})

	log.Println(string(requestBody))

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/acquire",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	resp, err := rec.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	log.Printf("%v", result)

	rec.RID = result["resourceId"].(string)

	return nil
}

// Start starts the recording
func (rec *Recorder) Start(channelTitle string, secret *string) error {
	// currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	currentTimeStamp := time.Now().In(location)
	currentDate := currentTimeStamp.Format("20060102")
	currentTime := currentTimeStamp.Format("150405")

	transcodingConfig := schemas.TranscodingConfig{
		Height:           720,
		Width:            1280,
		Bitrate:          2260,
		Fps:              15,
		MixedVideoLayout: 1,
		BackgroundColor:  "#000000",
	}
	var recordingConfig schemas.RecordingConfig
	if secret != nil && *secret != "" {
		recordingConfig = schemas.RecordingConfig{
			MaxIdleTime:       30,
			StreamTypes:       2,
			ChannelType:       1,
			DecryptionMode:    1,
			Secret:            *secret,
			TranscodingConfig: transcodingConfig,
		}
	} else {
		recordingConfig = schemas.RecordingConfig{
			MaxIdleTime:       30,
			StreamTypes:       2,
			ChannelType:       1,
			TranscodingConfig: transcodingConfig,
		}
	}

	recordingRequest := schemas.StartRecordRequest{
		Cname: rec.Channel,
		UID:   strconv.Itoa(int(rec.UID)),
		ClientRequest: schemas.ClientRequest{
			Token: rec.Token,
			StorageConfig: schemas.StorageConfig{
				Vendor:    viper.GetInt("RECORDING_VENDOR"),
				Region:    viper.GetInt("RECORDING_REGION"),
				Bucket:    viper.GetString("BUCKET_NAME"),
				AccessKey: viper.GetString("BUCKET_ACCESS_KEY"),
				SecretKey: viper.GetString("BUCKET_ACCESS_SECRET"),
				FileNamePrefix: []string{
					channelTitle, currentDate, currentTime,
				},
			},
			RecordingFileConfig: schemas.RecordingFileConfig{
				AVFileType: []string{"hls", "mp4"},
			},
			RecordingConfig: recordingConfig,
		},
	}

	requestBody, err := json.Marshal(&recordingRequest)
	if err != nil {
		return err
	}

	log.Print("https://api.agora.io/v1/apps/" + viper.GetString("APP_ID") + "/cloud_recording/resourceid/" + rec.RID + "/mode/mix/start")

	req, err := http.NewRequest("POST", "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rec.RID+"/mode/mix/start",
		bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	resp, err := rec.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	log.Printf("%v", result)
	rec.SID = result["sid"].(string)

	return nil
}

func (rec *Recorder) Stop() error {
	recordingRequest := schemas.AcquireRequest{
		Cname:         rec.Channel,
		UID:           strconv.Itoa(int(rec.UID)),
		ClientRequest: schemas.AcquireClientRequest{},
	}

	requestBody, _ := json.Marshal(&recordingRequest)

	url := "https://api.agora.io/v1/apps/"+viper.GetString("APP_ID")+"/cloud_recording/resourceid/"+rec.RID+"/sid/"+rec.SID+"/mode/mix/stop"
	req, err := http.NewRequest("POST",url ,
		bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(viper.GetString("CUSTOMER_ID"), viper.GetString("CUSTOMER_CERTIFICATE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func ImportEnv() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetDefault("PORT", 3000)
	viper.SetDefault("ENVIRONMENT", "developement")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found ignoring error
		} else {
			log.Panicln(fmt.Errorf("fatal error config file: %s", err))
		}
	}
}
