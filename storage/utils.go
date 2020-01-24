package storage

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	// PATTERNFILE ...
	PATTERNFILE = "data:([a-z]+/[a-z]+);base64,(.{5,})"
)

func getFileExtension(base64String string) (string, error) {

	re := regexp.MustCompile(PATTERNFILE)
	match := re.FindStringSubmatch(base64String)
	fmt.Println(match)
	if len(match) == 3 {
		if match[1] == "application/pdf" || match[1] == "text/plain" || match[1] == "image/jpeg" || match[1] == "image/png" || match[1] == "image/svg+xml" {
			log.Println("Tipo de arquivo -", match[1])
			return match[1], nil
		}
		return "", ErrTipoArquivoInvalido
	}
	return "", ErrParametroArquivo
}

// data:application/pdf;base64,BASE_64_STRING

func base64ToByteFile(base64 string) ([]byte, uintptr, string, error) {

	fileType, errType := getFileExtension(base64)
	if errType != nil {
		return nil, 0, "", errType
	}
	file, errorFile := b64.StdEncoding.DecodeString(strings.Split(base64, ",")[1])
	if errorFile != nil {
		return nil, 0, "", errorFile
	}
	size := unsafe.Sizeof(file)
	// data:application/pdf;base64,
	return file, size, fileType, nil
}

func generateUniqueName() string {
	current := time.Now()
	timestamp := current.Unix()
	s := strconv.FormatInt(timestamp, 10)

	return s
}
