package model

type PodEnv struct {
	ID     int64  `json:"id"`
	PodId  int64  `json:"pod_id"`
	EnvKey string `json:"env_key"`
	EnvVal string `json:"env_val"`
}
