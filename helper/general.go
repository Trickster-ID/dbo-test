package helper

func Ifelse(param1 interface{}, param2 interface{}) interface{} {
	if param1 == 0 || param1 == "" || param1 == float64(0) {
		return param2
	}
	return param1
}
