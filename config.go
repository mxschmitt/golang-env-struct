package envstruct

import (
	"os"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// ApplyEnvVars overrides the values in the given struct recursively by the
// represent environment variable which is given in the struct tag 'env'.
// For nested structs the struct tag 'env' from the overlying struct will be
// added as a prefix. Currently it can handle the types: string, int, bool
func ApplyEnvVars(c interface{}, prefix string) error {
	return applyEnvVar(reflect.ValueOf(c), reflect.TypeOf(c), -1, prefix)
}

// applyEnvVar does the same as ApplyEnvVars for recursion purposes.
func applyEnvVar(v reflect.Value, t reflect.Type, counter int, prefix string) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}
	f := reflect.StructField{}
	if counter != -1 {
		f = t.Field(counter)
	}
	v = reflect.Indirect(v)
	fieldEnv, exists := f.Tag.Lookup("env")
	env := os.Getenv(prefix + fieldEnv)
	if exists && env != "" {
		switch v.Kind() {
		case reflect.Int:
			envI, err := strconv.Atoi(env)
			if err != nil {
				return errors.Wrap(err, "could not parse to int")
			}
			v.SetInt(int64(envI))
		case reflect.String:
			v.SetString(env)
		case reflect.Bool:
			envB, err := strconv.ParseBool(env)
			if err != nil {
				return errors.Wrap(err, "could not parse bool")
			}
			v.SetBool(envB)
		}
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			if err := applyEnvVar(v.Field(i).Addr(), v.Type(), i, prefix+fieldEnv+"_"); err != nil {
				return errors.Wrap(err, "could not apply env var")
			}
		}
	}
	return nil
}
