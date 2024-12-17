package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HashPassword takes a string and hashes it using SHA-256, returning the
// resulting bytes as a hexadecimal string.
func HashPassword(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

// CheckPassHash takes a plaintext password and a hashed password, and
// returns true if the hashed password matches the given plaintext
// password when hashed with the same algorithm, and false otherwise.
//
// This function is safe to use with user-submitted passwords, as it
// uses a cryptographically secure hashing algorithm.
func CheckPassHash(pass string, hashed string) bool {
	attempPass := HashPassword(pass)
	log.Println(attempPass, "   :    ", hashed)
	return attempPass == hashed
}

// Sanitize takes a string and returns a sanitized version of it, where
// any invalid UTF-8 bytes are replaced with a replacement character,
// any leading or trailing whitespace is trimmed, any special or
// dangerous characters are removed, and any multiple spaces are
// replaced with a single space.
func Sanitize(input string) string {
	// Step 1: Trim leading and trailing whitespace
	input = strings.TrimSpace(input)

	// Step 2: Ensure valid UTF-8 encoding
	input = strings.ToValidUTF8(input, "")

	// Step 3: Remove special or dangerous characters using a regex
	// Allow only some special characters
	re := regexp.MustCompile(`[^\w\s,.-@]`)
	input = re.ReplaceAllString(input, "")

	// Step 4: Replace multiple spaces with a single space
	reSpaces := regexp.MustCompile(`\s+`)
	input = reSpaces.ReplaceAllString(input, " ")

	return input
}

// Generate a simple random token
func GenerateToken() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ParseDate takes a string in the format "yyyy-mm-dd" and parses it into a time.Time object.
// If the string is not in the correct format, an error is returned.
func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

// ParseBranchIDs takes a string of comma-separated branch IDs and
// returns a slice of the parsed IDs. If the string is empty, or if any
// of the IDs are invalid, an error is returned. The IDs are expected to
// be integers, and any leading or trailing whitespace is trimmed before
// parsing.
func ParseBranchIDs(branchIDs string) ([]int, error) {
	if branchIDs == "" {
		return nil, errors.New("branch_ids cannot be empty")
	}

	log.Println("branch_ids: ", branchIDs)
	idStrings := strings.Split(branchIDs, ",")
	ids := make([]int, len(idStrings))

	for i, idStr := range idStrings {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			return nil, errors.New("invalid branch_id: " + idStr)
		}
		ids[i] = id
	}
	log.Println("ids: ", ids)
	return ids, nil
}
