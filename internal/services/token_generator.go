package services

// TokenGenerator интерфейс для работы с токенами
type TokenGenerator interface {
	GenerateAccessToken(guid string, timeCreated int64) (string, error)
	GenerateRefreshToken(guid string, ip string, timeCreated int64) (string, error)
	StoreRefreshToken(guid string, refreshToken string) error
}

// Реальная реализация TokenGenerator
type RealTokenGenerator struct{}

func (r *RealTokenGenerator) GenerateAccessToken(guid string, timeCreated int64) (string, error) {
	return GenerateAccessToken(guid, timeCreated) // Используем текущую функцию
}

func (r *RealTokenGenerator) GenerateRefreshToken(guid string, ip string, timeCreated int64) (string, error) {
	return GenerateRefreshToken(guid, ip, timeCreated) // Используем текущую функцию
}

func (r *RealTokenGenerator) StoreRefreshToken(guid string, refreshToken string) error {
	return StoreRefreshToken(guid, refreshToken) // Используем текущую функцию
}
