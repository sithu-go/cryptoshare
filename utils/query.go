package utils

import "fmt"

func IdsIntToInCon(ids []uint64) string {
	inCon := ""
	for k, v := range ids {
		if k == 0 {
			inCon += fmt.Sprintf("%v", v)
		} else {
			inCon += fmt.Sprintf(",%v", v)
		}
	}
	if inCon == "" {
		inCon = "0"
	}
	return inCon
}
