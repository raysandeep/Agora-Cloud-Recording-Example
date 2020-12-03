package schemas

type StartCall struct {
	Uid     int    `json:"uid"`
	Channel string `json:"channel"`
}

type StopCall struct {
	Uid     int    `json:"uid"`
	Channel string `json:"channel"`
	Rid     string `json:"rid"`
	Sid     string `json:"sid"`
}

type UserCredentials struct {
	Rtc string `json:"rtc"`
	UID int    `json:"uid"`
}
