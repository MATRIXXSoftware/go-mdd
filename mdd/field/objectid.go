package field

// todo move to cmdc decode encode?
// type ObjectID string

// func NewObjectID(value string) (ObjectID, error) {
// 	var colonCount, dashCount int
// 	for _, c := range value {
// 		if c == ':' {
// 			colonCount++
// 		} else if c == '-' {
// 			dashCount++
// 		} else if !unicode.IsDigit(c) {
// 			return "", fmt.Errorf("invalid ObjectId format '%s'. Parts must be numeric", value)
// 		}
// 	}
// 	if !(dashCount == 3 && colonCount == 0) && !(dashCount == 0 && colonCount == 3) {
// 		return "", fmt.Errorf("invalid ObjectId format '%s'", value)
// 	}
// 	return ObjectID(value), nil
// }
