package util

//去除首尾大括号字符{}
func Substr(str string) string {

	rs := []rune(str)
	length := len(rs)

	// if start < 0 || start > length {
	// 	return ""
	// }

	// if end < 0 || end > length {
	// 	return ""
	// }
	return string(rs[1 : length-1])
}

// 某个值在数组中出现次数
func RepeatCount(search string, array []string) int {
	count := 0
	for _, v := range array {
		if search == v {
			count++
		}
	}
	return count
}

// 如果字符串有引号就去除
func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
