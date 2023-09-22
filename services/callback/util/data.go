package util

type progressData struct {
	WorkUUID  string  `json:"work_uuid"`
	SeverName string  `json:"sever_name"`
	Progress  float64 `json:"progress"`
}

type resultData struct {
	WorkUUID  string                 `json:"work_uuid"`
	SeverName string                 `json:"sever_name"`
	Result    map[string]interface{} `json:"result"`
}
