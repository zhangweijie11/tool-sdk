package global

type ParamsInterface interface {
	ValidWorkCreateParams(map[string]interface{}) (string, error)
}

type ParamsIns struct{}

func NewParamsIns() *ParamsIns {
	return &ParamsIns{}
}

func (pi *ParamsIns) ValidWorkCreateParams(params map[string]interface{}) (string, error) {
	return "", nil
}
