package dnhtutils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadDictionary(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening dictionary file: %w", err)
	}
	defer file.Close()

	var dictionary []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary = append(dictionary, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dictionary file: %w", err)
	}

	return dictionary, nil
}

func ReplacePlaceholders(input string, values []string) string {
	for i, value := range values {
		placeholder := fmt.Sprintf("{%d}", i)
		input = strings.ReplaceAll(input, placeholder, value)
	}
	return input
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("invalid integer: %s", s))
	}
	return i
}
