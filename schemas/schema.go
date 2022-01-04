package schemas

import "time"

type StartCall struct {
	Channel string `json:"channel"`
}

type StopCall struct {
	Channel string `json:"channel"`
}

type UserCredentials struct {
	Rtc string `json:"rtc"`
	UID int    `json:"uid"`
}

type CallStatus struct {
	Rid string `json:"rid"`
	Sid string `json:"sid"`
}

type AcquireClientRequest struct {
	ResourceExpiredHour int `json:"resourceExpiredHour,omitempty"`
}

type AcquireRequest struct {
	Cname         string               `json:"cname"`
	UID           string               `json:"uid"`
	ClientRequest AcquireClientRequest `json:"clientRequest"`
}

type TranscodingConfig struct {
	Height           int    `json:"height"`
	Width            int    `json:"width"`
	Bitrate          int    `json:"bitrate"`
	Fps              int    `json:"fps"`
	MixedVideoLayout int    `json:"mixedVideoLayout"`
	MaxResolutionUID string `json:"maxResolutionUid,omitempty"`
	BackgroundColor  string `json:"backgroundColor"`
}

type RecordingConfig struct {
	MaxIdleTime       int               `json:"maxIdleTime"`
	StreamTypes       int               `json:"streamTypes"`
	ChannelType       int               `json:"channelType"`
	DecryptionMode    int               `json:"decryptionMode,omitempty"`
	Secret            string            `json:"secret,omitempty"`
	TranscodingConfig TranscodingConfig `json:"transcodingConfig"`
}

type StorageConfig struct {
	Vendor         int      `json:"vendor"`
	Region         int      `json:"region"`
	Bucket         string   `json:"bucket"`
	AccessKey      string   `json:"accessKey"`
	SecretKey      string   `json:"secretKey"`
	FileNamePrefix []string `json:"fileNamePrefix"`
}

type RecordingFileConfig struct {
	AVFileType []string `json:"avFileType"`
}

type ClientRequest struct {
	Token               string              `json:"token"`
	RecordingConfig     RecordingConfig     `json:"recordingConfig"`
	RecordingFileConfig RecordingFileConfig `json:"recordingFileConfig"`
	StorageConfig       StorageConfig       `json:"storageConfig"`
}

type StartRecordRequest struct {
	Cname         string        `json:"cname"`
	UID           string        `json:"uid"`
	ClientRequest ClientRequest `json:"clientRequest"`
}

type BubbleRecord struct {
	CreatedBy         string    `json:"Created By"`
	CreatedDate       time.Time `json:"Created Date"`
	ModifiedDate      time.Time `json:"Modified Date"`
	CustomerFirstName string    `json:"Customer First Name"`
	CustomerLastName  string    `json:"Customer Last Name"`
	CustomerEmail     string    `json:"Customer email"`
	CustomerPhone     string    `json:"Customer Phone"`
	ShortID           string    `json:"short_id"`
	Rid               string    `json:"rid"`
	Sid               string    `json:"sid"`
	UID               string    `json:"uid"`
	ID                string    `json:"_id"`
}

type BubbleResponse struct {
	Response struct {
		Cursor    int            `json:"cursor"`
		Results   []BubbleRecord `json:"results"`
		Remaining int            `json:"remaining"`
		Count     int            `json:"count"`
	} `json:"response"`
}
