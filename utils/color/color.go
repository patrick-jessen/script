package color

import "fmt"

func Red(str string) string {
	return fmt.Sprintf("\033[1;31m%v\033[0m", str)
}

func Yellow(str string) string {
	return fmt.Sprintf("\033[1;33m%v\033[0m", str)
}

func Green(str string) string {
	return fmt.Sprintf("\033[1;32m%v\033[0m", str)
}
