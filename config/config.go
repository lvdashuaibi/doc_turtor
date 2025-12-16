package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config 存储应用配置
type Config struct {
	DashScopeAPIKey string
	DashScopeURL    string
}

// LoadConfig 从.env文件加载配置
func LoadConfig() (*Config, error) {
	config := &Config{}

	// 首先尝试从环境变量读取
	config.DashScopeAPIKey = os.Getenv("DASHSCOPE_API_KEY")
	config.DashScopeURL = os.Getenv("DASHSCOPE_BASE_URL")

	// 如果环境变量为空，尝试从.env文件读取
	if config.DashScopeAPIKey == "" || config.DashScopeURL == "" {
		err := loadEnvFile(".env", config)
		if err != nil {
			return nil, fmt.Errorf("加载配置失败: %v", err)
		}
	}

	// 验证必要的配置
	if config.DashScopeAPIKey == "" {
		return nil, fmt.Errorf("缺少 DASHSCOPE_API_KEY 配置")
	}
	if config.DashScopeURL == "" {
		return nil, fmt.Errorf("缺少 DASHSCOPE_BASE_URL 配置")
	}

	// 将配置值设置到系统环境变量中
	os.Setenv("DASHSCOPE_API_KEY", config.DashScopeAPIKey)
	os.Setenv("DASHSCOPE_BASE_URL", config.DashScopeURL)

	return config, nil
}

// loadEnvFile 从.env文件读取配置
func loadEnvFile(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开 %s 文件: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除引号（如果有）
		value = strings.Trim(value, "\"'")

		switch key {
		case "DASHSCOPE_API_KEY":
			config.DashScopeAPIKey = value
		case "DASHSCOPE_BASE_URL":
			config.DashScopeURL = value
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 %s 文件出错: %v", filename, err)
	}

	return nil
}

// PrintConfig 打印配置信息（隐藏敏感信息）
func (c *Config) PrintConfig() {
	fmt.Println("🔧 DashScope 配置信息:")
	// 安全地处理 API Key 的显示，避免短字符串导致的 panic
	apiKey := c.DashScopeAPIKey
	if len(apiKey) > 14 {
		apiKey = apiKey[:10] + "***" + apiKey[len(apiKey)-4:]
	} else if len(apiKey) > 0 {
		apiKey = apiKey[:1] + "***" + apiKey[len(apiKey)-1:]
	} else {
		apiKey = "***"
	}
	fmt.Printf("  • API Key: %s\n", apiKey)
	fmt.Printf("  • Base URL: %s\n", c.DashScopeURL)
}
