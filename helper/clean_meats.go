package helper

import (
	"regexp"
	"strings"
)

// Clean and count meats in the meat string
func CleanMeats(input string) []string {
	// Remove unwanted characters (., .., spaces)
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s-]`) //[^\w\s]
	cleaned := reg.ReplaceAllString(input, "")

	// Split the string by spaces to get individual meat names
	meats := strings.Fields(cleaned)
	return meats
}
