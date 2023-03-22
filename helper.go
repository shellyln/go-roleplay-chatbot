package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"os"
	"strings"
)

//go:embed .env
var dotenvContent []byte

// dotenvファイルを読み込む
func readDotenv() error {
	// ファイルをスキャンして、環境変数に設定する
	scanner := bufio.NewScanner(bytes.NewReader(dotenvContent))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue // コメント行は無視する
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // フォーマットが正しくない行は無視する
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = strings.Trim(value, "'")
		} else if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}
		os.Setenv(key, value)
	}

	return nil
}
