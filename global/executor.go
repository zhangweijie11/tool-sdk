package global

type ExecutorInterface interface {
	ValidWorkCreateParams(map[string]interface{}) (string, error)
	ExecutorMainFunc(interface{}) error
}

type ExecutorIns struct{}

func NewExecutorIns() *ExecutorIns {
	return &ExecutorIns{}
}

func (ei *ExecutorIns) ExecutorMainFunc(params interface{}) error {
	return nil
}

func (ei *ExecutorIns) ValidWorkCreateParams(params map[string]interface{}) (string, error) {
	return "", nil
}
