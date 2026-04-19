package config

// GetCryptoKey returns the crypto key (read at call time, not init time)
func GetCryptoKey() string {
	return envOr("CRYPTO_KEY", envOr("APP_KEY", ""))
}

// GetCryptoCipher returns the crypto cipher
func GetCryptoCipher() string {
	return envOr("CRYPTO_CIPHER", "AES-256-CBC")
}

// Legacy variables for backwards compatibility (set after godotenv loads)
var (
	CryptoKey    string
	CryptoCipher string
)

func InitCrypto() {
	CryptoKey = GetCryptoKey()
	CryptoCipher = GetCryptoCipher()
}
