package config

// GetAppName returns the app name (read at call time)
func GetAppName() string {
	return envOr("APP_NAME", "Velocity")
}

// GetAppEnv returns the app environment (read at call time)
func GetAppEnv() string {
	return envOr("APP_ENV", "development")
}

// GetPort returns the port (read at call time)
func GetPort() string {
	return envOr("PORT", "4000")
}
