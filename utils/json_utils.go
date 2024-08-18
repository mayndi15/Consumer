package utils

import (
	"strings"
)

func CleanJsonString(jsonString string) string {
    if len(jsonString) > 0 && jsonString[0] == '"' && jsonString[len(jsonString)-1] == '"' {
        jsonString = jsonString[1 : len(jsonString)-1]
    }
    unescapedString := strings.ReplaceAll(jsonString, `\"`, `"`)
    unescapedString = strings.ReplaceAll(unescapedString, `\\`, `\`)
    unescapedString = strings.ReplaceAll(unescapedString, `\n`, "\n")
    unescapedString = strings.ReplaceAll(unescapedString, `\t`, "\t")
    unescapedString = strings.ReplaceAll(unescapedString, `\r`, "\r")

    unescapedString = strings.TrimSpace(unescapedString)

    return unescapedString
}
