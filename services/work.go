package services

type WorkApi interface {
	CreateWork(workParams interface{}) error
	DeleteWork(workParams interface{}) error
	UpdateWork(workParams interface{}) error
	GetWorkInfo(workParams interface{}) error
	GetWorkStatus(workParams interface{}) error
	CallbackWork(workParams interface{}) error
}
