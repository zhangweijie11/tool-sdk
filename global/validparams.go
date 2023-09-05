package global

import "fmt"

type ParamsInterface interface {
	CheckParamsInterface()
	ValidWorkCreateParams(map[string]interface{}) error
	//ValidWorkDeleteParams(map[string]interface{}) error
	//ValidWorkUpdateParams(map[string]interface{}) error
	//ValidWorkGetInfoParams(map[string]interface{}) error
	//ValidWorkGetStatusParams(map[string]interface{}) error
}

type paramsIns struct{}

func (pi *paramsIns) ValidWorkCreateParams(params map[string]interface{}) error {
	fmt.Println("------------>本地参数", params)
	return nil
}

func (pi *paramsIns) CheckParamsInterface() {
	fmt.Println("------------>", "检查是否实现参数校验接口")
}

func CheckParamsInterface(obj ParamsInterface) {
	obj.CheckParamsInterface()
}
