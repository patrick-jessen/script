package color

import "fmt"

func Red(val interface{}) string {
	return fmt.Sprintf("\033[1;31m%v\033[0m", val)
}

func Yellow(val interface{}) string {
	return fmt.Sprintf("\033[1;33m%v\033[0m", val)
}

func Green(val interface{}) string {
	return fmt.Sprintf("\033[1;32m%v\033[0m", val)
}

func Blue(val interface{}) string {
	return fmt.Sprintf("\033[1;34m%v\033[0m", val)
}
