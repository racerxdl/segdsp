// +build !amd64,!arm64

package native

func GetSIMDMode() string {
	return "NONE"
}
