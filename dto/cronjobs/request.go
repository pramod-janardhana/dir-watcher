package cronjobs

type UpdateConfigReq struct {
	CronExpresion string `json:"cronExpresion"`
	MagicString   string `json:"magicString"`
}
