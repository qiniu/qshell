package utils

func GetNotEmptyStringIfExist(values ... string) string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return ""
}

func GetTrueBoolValueIfExist(values ... bool) bool {
	for _, value := range values {
		if value {
			return value
		}
	}
	return false
}

func GetNotZeroIntIfExist(values ... int) int {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroUIntIfExist(values ... uint) uint {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroInt64IfExist(values ... int64) int64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroUInt64IfExist(values ... uint64) uint64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroInt16IfExist(values ... int16) int16 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroUInt16IfExist(values ... uint16) uint16 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}


func GetNotZeroInt8IfExist(values ... int8) int8 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func GetNotZeroUInt8IfExist(values ... uint8) uint8 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}
