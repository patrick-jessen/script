import c "user32.dll"

func print(first string, second string) {
    c.MessageBoxA(0, first, second, 0)
}