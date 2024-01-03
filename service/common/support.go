package common

func MapStatusFpt(status string) string {
	if status == "1" {
		return "success"
	} else if status == "2" || status == "-11" {
		return "pending"
	} else if status == "0" {
		return "fail"
	}
	return "fail"
}
