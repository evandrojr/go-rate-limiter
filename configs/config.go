package configs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/spf13/viper"
)

type config struct {
	// DBDriver          string `mapstructure:"DB_DRIVER"`
	// DBHost            string `mapstructure:"DB_HOST"`
	// DBPort            string `mapstructure:"DB_PORT"`
	// DBUser            string `mapstructure:"DB_USER"`
	// DBPassword        string `mapstructure:"DB_PASSWORD"`
	// DBName            string `mapstructure:"DB_NAME"`
	// WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	// GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	// GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	Tokens map[string]int
}

var Config config

// func LoadConfig(path string) (*conf, error) {
// 	var cfg *conf
// 	viper.SetConfigName("app_config")
// 	viper.SetConfigType("env")
// 	viper.AddConfigPath(path)
// 	viper.SetConfigFile(".env")
// 	viper.AutomaticEnv()
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = viper.Unmarshal(&cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cfg, err
// }

func LoadConfig() {
	// Configurar o Viper para ler o .env
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar o .env: %s", err))
	}

	// Ler o valor do MY_MAP
	rawMap := viper.GetString("TOKENS")

	// Converter para map[string]int
	parsedMap, parseErr := parseStringToMapInt(rawMap)
	if parseErr != nil {
		panic(fmt.Errorf("erro ao converter IPs: %s", parseErr))
	}

	// fmt.Printf("Mapa: %+v\n", parsedMap)
	Config.Tokens = parsedMap
	pretty.Println(Config)
}

// Função para converter string delimitada para map[string]int
func parseStringToMapInt(input string) (map[string]int, error) {
	result := make(map[string]int)
	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("valor inválido '%s' para chave '%s'", parts[1], key)
			}
			result[key] = value
		}
	}
	return result, nil
}
