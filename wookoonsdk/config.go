package wookoonsdk

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// AuthConfig is the core configuration.
// The first step to use this SDK is to use the InitConfig function to initialize the global authConfig.
type AuthConfig struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	Organization string
	Application  string
	Certificate  string
}

type AppConfig struct {
	ServerPort string
	ServerMode string
}

type CorsConfig struct {
	AllowOrigins []string
}

type JwtConfig struct {
	JWTSecret string
	JWTExpire time.Duration
	JWTIssuer string
}

// 加载 .env
func (s *Server) loadConfigs() {
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  未找到 .env 文件，将使用系统环境变量: %v", err)
		return
	} else {
		log.Println("✅ 成功加载 .env 文件")
	}

	s.AppConfig = &AppConfig{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerMode: getEnv("SERVER_MODE", "release"),
	}

	s.CorsConfig = &CorsConfig{
		AllowOrigins: getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
	}

	expireHours := getEnvAsInt("JWT_EXPIRE_HOURS", 24)
	s.JwtConfig = &JwtConfig{
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpire: time.Duration(expireHours) * time.Hour,
		JWTIssuer: getEnv("JWT_ISSUER", "wookoon-go-cas"),
	}

	s.AuthConfig = &AuthConfig{
		Endpoint:     getEnv("CASDOOR_ENDPOINT", ""),
		ClientID:     getEnv("CASDOOR_CLIENT_ID", ""),
		ClientSecret: getEnv("CASDOOR_CLIENT_SECRET", ""),
		Organization: getEnv("CASDOOR_ORGANIZATION", ""),
		Application:  getEnv("CASDOOR_APPLICATION", ""),
	}

	s.loadCertificate()
	s.validate()
}

// 加载证书
func (s *Server) loadCertificate() {
	// 优先从文件读取
	certFile := getEnv("CASDOOR_CERTIFICATE_FILE", "")
	if certFile != "" {
		certBytes, err := os.ReadFile(certFile)
		if err != nil {
			log.Fatalf("❌ 读取证书文件失败 [%s]: %v", certFile, err)
		}

		s.AuthConfig.Certificate = string(certBytes)
		log.Printf("✅ 从文件加载证书: %s (%d bytes)", certFile, len(certBytes))
		return
	}

	// 其次从环境变量读取（单行格式，\n 需要还原）
	certEnv := getEnv("CASDOOR_CERTIFICATE", "")
	if certEnv != "" {
		// 将字面 \n 转为真实换行符
		s.AuthConfig.Certificate = strings.ReplaceAll(certEnv, "\\n", "\n")
		log.Println("✅ 从环境变量加载证书")
		return
	}

	log.Println("⚠️  未配置证书（CASDOOR_CERTIFICATE_FILE 或 CASDOOR_CERTIFICATE）")
}

func (s *Server) validate() {
	if s.AuthConfig.Endpoint == "" {
		log.Fatal("❌ CASDOOR_ENDPOINT 未配置")
	}
	if s.AuthConfig.ClientID == "" {
		log.Fatal("❌ CASDOOR_CLIENT_ID 未配置")
	}
	if s.AuthConfig.ClientSecret == "" {
		log.Fatal("❌ CASDOOR_CLIENT_SECRET 未配置")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		// 去除每个元素的首尾空格
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}
