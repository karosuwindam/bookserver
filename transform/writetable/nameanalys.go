package writetable

// 文字列から()で囲まれた部分を取り出し削除する
func removeParentheses(str string) ([]string, string) {
	output := []string{}
	// ()で囲まれた部分を取得する
	if len(str) > 0 {
		for i := 0; i < len(str); i++ {
			if str[i] == '(' {
				for j := i + 1; j < len(str); j++ {
					if str[j] == ')' {
						output = append(output, str[i+1:j])
						str = str[:i] + str[j+1:]
						i--
						break
					}
				}
			}
		}
	}
	return output, str
}

// 文字列から[]で囲まれた部分を取り出し削除する
func removeBrackets(str string) ([]string, string) {
	output := []string{}
	// []で囲まれた部分を取得する
	if len(str) > 0 {
		for i := 0; i < len(str); i++ {
			if str[i] == '[' {
				for j := i + 1; j < len(str); j++ {
					if str[j] == ']' {
						output = append(output, str[i+1:j])
						str = str[:i] + str[j+1:]
						i--
						break
					}
				}
			}
		}
	}
	return output, str
}
