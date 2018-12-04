package native

func GetSIMDMode() string {
	// Neon is always available at AArch64
	// Disabled for now because we don't have AArch64 support on c2goasm
	// return "NEON"

	return "NONE"
}
