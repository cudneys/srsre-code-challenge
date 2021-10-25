package password

func GetDefaultValue(val, length int) int {
	if val == 0 {
		return length / 4
	}
	return val
}
