package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// StringProperties holds the computed properties of an analyzed string
type StringProperties struct {
	Length               int            `json:"length"`
	IsPalindrome         bool           `json:"is_palindrome"`
	UniqueCharacters     int            `json:"unique_characters"`
	WordCount            int            `json:"word_count"`
	SHA256Hash           string         `json:"sha256_hash"`
	CharacterFrequencyMap map[string]int `json:"character_frequency_map"`
}

// AnalyzedString represents a stored analyzed string with metadata
type AnalyzedString struct {
	ID         string            `json:"id"`
	Value      string            `json:"value"`
	Properties StringProperties  `json:"properties"`
	CreatedAt  time.Time         `json:"created_at"`
}

// CreateStringRequest is the request body for creating/analyzing a string
type CreateStringRequest struct {
	Value string `json:"value" binding:"required"`
}

// FilterParams holds query parameters for filtering
type FilterParams struct {
	IsPalindrome     *bool   `form:"is_palindrome"`
	MinLength        *int    `form:"min_length"`
	MaxLength        *int    `form:"max_length"`
	WordCount        *int    `form:"word_count"`
	ContainsCharacter *string `form:"contains_character"`
}

// NaturalLanguageQuery holds parsed natural language query results
type NaturalLanguageQuery struct {
	Original      string                 `json:"original"`
	ParsedFilters map[string]interface{} `json:"parsed_filters"`
}

// FilterResponse wraps filtered results with metadata
type FilterResponse struct {
	Data          []AnalyzedString         `json:"data"`
	Count         int                      `json:"count"`
	FiltersApplied map[string]interface{}  `json:"filters_applied,omitempty"`
	InterpretedQuery *NaturalLanguageQuery `json:"interpreted_query,omitempty"`
}

// StringStore is a simple in-memory storage for analyzed strings
var stringStore = make(map[string]*AnalyzedString)

func main() {
	router := gin.Default()

	// Routes
	router.POST("/strings", createString)
	router.GET("/strings", getAllStrings)
	router.GET("/strings/filter-by-natural-language", filterByNaturalLanguage)
	router.GET("/strings/:value", getString)
	router.DELETE("/strings/:value", deleteString)

	router.Run(":8080")
}

// createString handles POST /strings
func createString(c *gin.Context) {
	var req CreateStringRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing 'value' field"})
		return
	}

	// Validate that value is a string
	if len(req.Value) == 0 && req.Value != "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid data type for 'value' (must be string)"})
		return
	}

	// Analyze the string
	props := analyzeString(req.Value)
	hash := props.SHA256Hash

	// Check if string already exists
	if _, exists := stringStore[hash]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "String already exists in the system"})
		return
	}

	// Create and store the analyzed string
	analyzed := &AnalyzedString{
		ID:        hash,
		Value:     req.Value,
		Properties: props,
		CreatedAt: time.Now().UTC(),
	}

	stringStore[hash] = analyzed

	c.JSON(http.StatusCreated, analyzed)
}

// getString handles GET /strings/:value
func getString(c *gin.Context) {
	value := c.Param("value")

	// Calculate hash of the provided value to lookup
	hash := calculateSHA256(value)

	// Find the string in store
	analyzed, exists := stringStore[hash]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system"})
		return
	}

	c.JSON(http.StatusOK, analyzed)
}

// getAllStrings handles GET /strings with optional filtering
func getAllStrings(c *gin.Context) {
	var params FilterParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
		return
	}

	// Collect all strings and apply filters
	var results []AnalyzedString
	for _, analyzed := range stringStore {
		if matchesFilters(analyzed, params) {
			results = append(results, *analyzed)
		}
	}

	// Sort by creation time for consistent ordering
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	// Build filters applied map
	filtersApplied := make(map[string]interface{})
	if params.IsPalindrome != nil {
		filtersApplied["is_palindrome"] = *params.IsPalindrome
	}
	if params.MinLength != nil {
		filtersApplied["min_length"] = *params.MinLength
	}
	if params.MaxLength != nil {
		filtersApplied["max_length"] = *params.MaxLength
	}
	if params.WordCount != nil {
		filtersApplied["word_count"] = *params.WordCount
	}
	if params.ContainsCharacter != nil {
		filtersApplied["contains_character"] = *params.ContainsCharacter
	}

	response := FilterResponse{
		Data:          results,
		Count:         len(results),
		FiltersApplied: filtersApplied,
	}

	c.JSON(http.StatusOK, response)
}

// filterByNaturalLanguage handles GET /strings/filter-by-natural-language
func filterByNaturalLanguage(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'query' parameter"})
		return
	}

	// Parse natural language query
	parsedFilters, err := parseNaturalLanguageQuery(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to parse natural language query: %v", err)})
		return
	}

	// Convert parsed filters to FilterParams
	filterParams := convertParsedFiltersToParams(parsedFilters)

	// Collect all strings and apply filters
	var results []AnalyzedString
	for _, analyzed := range stringStore {
		if matchesFilters(analyzed, filterParams) {
			results = append(results, *analyzed)
		}
	}

	// Sort by creation time
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	interpretedQuery := &NaturalLanguageQuery{
		Original:      query,
		ParsedFilters: parsedFilters,
	}

	response := FilterResponse{
		Data:                results,
		Count:               len(results),
		InterpretedQuery:    interpretedQuery,
	}

	c.JSON(http.StatusOK, response)
}

// deleteString handles DELETE /strings/:value
func deleteString(c *gin.Context) {
	value := c.Param("value")

	// Calculate hash of the provided value to lookup
	hash := calculateSHA256(value)

	// Check if string exists
	if _, exists := stringStore[hash]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system"})
		return
	}

	// Delete the string
	delete(stringStore, hash)

	c.JSON(http.StatusNoContent, nil)
}

// analyzeString computes all properties for a given string
func analyzeString(value string) StringProperties {
	length := utf8.RuneCountInString(value)
	isPalindrome := checkPalindrome(value)
	uniqueChars := countUniqueCharacters(value)
	wordCount := countWords(value)
	sha256Hash := calculateSHA256(value)
	charFreqMap := buildCharacterFrequencyMap(value)

	return StringProperties{
		Length:               length,
		IsPalindrome:         isPalindrome,
		UniqueCharacters:     uniqueChars,
		WordCount:            wordCount,
		SHA256Hash:           sha256Hash,
		CharacterFrequencyMap: charFreqMap,
	}
}

// checkPalindrome checks if a string is a palindrome (case-insensitive)
func checkPalindrome(s string) bool {
	// Remove spaces and convert to lowercase for comparison
	cleaned := strings.ToLower(regexp.MustCompile(`\s+`).ReplaceAllString(s, ""))
	runes := []rune(cleaned)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// countUniqueCharacters counts distinct characters in a string
func countUniqueCharacters(s string) int {
	charSet := make(map[rune]bool)
	for _, r := range strings.ToLower(s) {
		if r != ' ' && r != '\t' && r != '\n' {
			charSet[r] = true
		}
	}
	return len(charSet)
}

// countWords counts whitespace-separated words
func countWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// calculateSHA256 computes the SHA256 hash of a string
func calculateSHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// buildCharacterFrequencyMap creates a map of character frequencies
func buildCharacterFrequencyMap(s string) map[string]int {
	freqMap := make(map[string]int)
	for _, r := range strings.ToLower(s) {
		if r != ' ' && r != '\t' && r != '\n' {
			freqMap[string(r)]++
		}
	}
	return freqMap
}

// matchesFilters checks if an analyzed string matches all applied filters
func matchesFilters(analyzed *AnalyzedString, params FilterParams) bool {
	if params.IsPalindrome != nil && analyzed.Properties.IsPalindrome != *params.IsPalindrome {
		return false
	}

	if params.MinLength != nil && analyzed.Properties.Length < *params.MinLength {
		return false
	}

	if params.MaxLength != nil && analyzed.Properties.Length > *params.MaxLength {
		return false
	}

	if params.WordCount != nil && analyzed.Properties.WordCount != *params.WordCount {
		return false
	}

	if params.ContainsCharacter != nil {
		char := strings.ToLower(*params.ContainsCharacter)
		if _, exists := analyzed.Properties.CharacterFrequencyMap[char]; !exists {
			return false
		}
	}

	return true
}

// parseNaturalLanguageQuery parses natural language filter queries
func parseNaturalLanguageQuery(query string) (map[string]interface{}, error) {
	lowerQuery := strings.ToLower(query)
	filters := make(map[string]interface{})

	// Check for word count patterns
	if strings.Contains(lowerQuery, "single word") {
		filters["word_count"] = 1
	} else if strings.Contains(lowerQuery, "two word") || strings.Contains(lowerQuery, "2 word") {
		filters["word_count"] = 2
	} else if strings.Contains(lowerQuery, "three word") || strings.Contains(lowerQuery, "3 word") {
		filters["word_count"] = 3
	}

	// Check for palindrome pattern
	if strings.Contains(lowerQuery, "palindrom") {
		filters["is_palindrome"] = true
	}

	// Check for length patterns
	lengthPattern := regexp.MustCompile(`longer than (\d+)`)
	if matches := lengthPattern.FindStringSubmatch(lowerQuery); matches != nil {
		minLength := 0
		fmt.Sscanf(matches[1], "%d", &minLength)
		filters["min_length"] = minLength + 1
	}

	shorterPattern := regexp.MustCompile(`shorter than (\d+)`)
	if matches := shorterPattern.FindStringSubmatch(lowerQuery); matches != nil {
		maxLength := 0
		fmt.Sscanf(matches[1], "%d", &maxLength)
		filters["max_length"] = maxLength - 1
	}

	// Check for character patterns
	charPattern := regexp.MustCompile(`(?:contain|with) (?:the |letter |character )?'?([a-z])'?`)
	if matches := charPattern.FindStringSubmatch(lowerQuery); matches != nil {
		filters["contains_character"] = matches[1]
	}

	// Check for vowel patterns
	if strings.Contains(lowerQuery, "first vowel") {
		filters["contains_character"] = "a"
	} else if strings.Contains(lowerQuery, "last vowel") {
		filters["contains_character"] = "u"
	}

	// If no filters could be parsed, return error
	if len(filters) == 0 {
		return nil, fmt.Errorf("unable to parse any filters from query")
	}

	return filters, nil
}

// convertParsedFiltersToParams converts parsed filters map to FilterParams struct
func convertParsedFiltersToParams(parsed map[string]interface{}) FilterParams {
	params := FilterParams{}

	if isPalin, ok := parsed["is_palindrome"].(bool); ok {
		params.IsPalindrome = &isPalin
	}

	if minLen, ok := parsed["min_length"].(int); ok {
		params.MinLength = &minLen
	}

	if maxLen, ok := parsed["max_length"].(int); ok {
		params.MaxLength = &maxLen
	}

	if wordCnt, ok := parsed["word_count"].(int); ok {
		params.WordCount = &wordCnt
	}

	if char, ok := parsed["contains_character"].(string); ok {
		params.ContainsCharacter = &char
	}

	return params
}