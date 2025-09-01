package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type HostConfig struct {
	Name           string `mapstructure:"name"`
	IP             string `mapstructure:"ip"`
	Port           int    `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	AuthMethod     string `mapstructure:"auth_method"`
	Password       string `mapstructure:"password"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
}

type GroupConfig struct {
	Default []string `mapstructure:"default"`
}
type AppConfig struct {
	Groups    GroupConfig  `mapstructure:"groups"`
	Hosts     []HostConfig `mapstructure:"hosts"`
	Command   string
	Nodes     []string // 从命令行获取的节点列表
	Debug     bool     // 是否启用调试模式
	LogFormat string   // 日志格式: text或json
	LogFile   string   // 日志文件路径
}

func ShowHelp() {
	fmt.Printf("SSH Task Runner (ssht)\n\n")
	fmt.Printf("Usage:\n")
	pflag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  # Basic usage\n")
	fmt.Printf("  ./ssht --command \"hostname\"\n\n")
	fmt.Printf("  # Run on specific nodes\n")
	fmt.Printf("  ./ssht --command \"hostname\" --nodes node1,node2\n\n")
	fmt.Printf("  # Debug mode with JSON logging\n")
	fmt.Printf("  ./ssht --command \"hostname\" --nodes node1 --debug --log-format json\n\n")
	fmt.Printf("  # Write logs to file\n")
	fmt.Printf("  ./ssht --command \"hostname\" --log-file output.log\n")
}

func Load() (*AppConfig, error) {
	// 初始化Viper
	v := viper.New()
	v.SetConfigName("config.toml") // 配置文件名称
	v.SetConfigType("toml")        // 配置文件类型
	v.AddConfigPath(".")           // 配置文件路径
	v.AddConfigPath("./config")

	// 设置命令行参数
	pflag.String("command", "", "SSH command to execute")
	pflag.StringSlice("nodes", []string{}, "List of nodes to execute on (default: all nodes)")
	pflag.Bool("debug", false, "Enable debug logging")
	pflag.String("log-format", "text", "Log format: text or json")
	pflag.String("log-file", "", "Log file path (default: stdout)")
	pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()

	if help, _ := pflag.CommandLine.GetBool("help"); help {
		ShowHelp()
		os.Exit(0)
	}

	v.BindPFlags(pflag.CommandLine)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从命令行获取参数
	cfg.Command = v.GetString("command")
	if cfg.Command == "" {
		log.Fatal("必须通过--command参数指定要执行的命令")
	}
	cfg.Nodes = v.GetStringSlice("nodes")
	cfg.Debug = v.GetBool("debug")
	cfg.LogFormat = v.GetString("log-format")
	cfg.LogFile = v.GetString("log-file")

	return &cfg, nil
}
