package integrationtest

//IsInArrayString will return true if check exists in arr
func IsInArrayString(arr []string, check string) bool {
	for _, v := range arr {
		if v == check {
			return true
		}
	}
	return false
}

//CombineArrayString sum of two array of string
func CombineArrayString(arr1 []string, arr2 []string) []string {
	for _, str := range arr2 {
		arr1 = append(arr1, str)
	}
	return arr1
}
