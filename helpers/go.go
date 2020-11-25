package helpers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
)

// StringBetween Get String Between
func StringBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// StringBefore Get String Before
func StringBefore(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

// StringAfter Get String After
func StringAfter(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

// Contains : Check if element in array.
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// UniqueStrings make string array unique
func UniqueStrings(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

// IsValidUUID Check if given string is uuid
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

//EncodeMessageUTF16 EncodeMessageUTF16
func EncodeMessageUTF16(message string) string {
	runeByte := []rune(message)
	encodedByte := utf16.Encode(runeByte)
	var buf bytes.Buffer
	defer buf.Reset()
	for _, num := range encodedByte {
		buf.WriteString(fmt.Sprintf("%04X", num))
	}
	return buf.String()
}

func IsNewer(time1 string, time2 string) bool {
	obj1, err := time.Parse(time.RFC3339, time1)
	if err != nil {
		return false
	}

	obj2, err := time.Parse(time.RFC3339, time2)
	if err != nil {
		return false
	}
	return obj1.After(obj2)
}
