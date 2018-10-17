package main

import (
	"log"

	envstruct "github.com/mxschmitt/golang-env-struct"
)

// Settings is a example struct which can be used with envstruct
type Settings struct {
	StringVar  string `env:"STRING_VAR"`
	IntegerVar int    `env:"INTEGER_VAR"`
	Nested     struct {
		StringVar  string `env:"STRING_VAR"`
		IntegerVar int    `env:"INTEGER_VAR"`
	} `env:"NESTED"`
}

func main() {
	settings := Settings{}
	if err := envstruct.ApplyEnvVars(&settings, "EXAMPLE_APP"); err != nil {
		log.Fatalf("could not apply environment variables: %v", err)
	}
	log.Printf("Settings: %+v\n", settings)
}
