// +build !amd64

package native

func GetSIMDMode() string {
	return "NONE"
}
