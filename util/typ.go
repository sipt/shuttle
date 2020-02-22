package util

// []interface{} => []string
func InterfaceSliceToStringSlice(params []interface{}) []string {
	reply := make([]string, len(params))
	for i, v := range params {
		reply[i], _ = v.(string)
	}
	return reply
}
