package upload

import "fmt"

// GetTarget makes target parameter for uploading.
func GetTarget(dirId string) string {
	return fmt.Sprintf("U_1_%s", dirId)
}
