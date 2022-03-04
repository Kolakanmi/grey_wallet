package envconfig

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var envPrefix = ""

//Load func
func Load(t interface{}) {
	LoadWithPrefix(envPrefix, t)
}

//LoadWithPrefix func
func LoadWithPrefix(prefix string, t interface{}) {
	if err := envconfig.Process(prefix, t); err != nil {
		log.Println("Unable to read: ", err)
	}
}

//SetEnvFromConfig - sets up env file and calls os.Setenv on each
func SetEnvFromConfig(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "#") || strings.TrimSpace(txt) == "" {
			continue
		}
		env := strings.SplitN(txt, "=", 2)
		if len(env) != 2 {
			return errors.New("must be key-value pair")
		}
		key := env[0]
		value := env[1]
		value = strings.Trim(value, `"`)
		_, ok := os.LookupEnv(key)
		if ok {
			continue
		}
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
