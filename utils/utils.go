package utils

func Contains(sli []string, str string) bool {
	for _, v := range sli {
		if v == str {
			return true
		}
	}
	return false
}
