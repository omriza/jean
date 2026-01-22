package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
)

const (
	chromiumSalt       = "saltysalt"
	chromiumIterations = 1003
	chromiumKeyLength  = 16
)

var (
	uuidRegex       = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	sessionKeyRegex = regexp.MustCompile(`sk-ant-sid[A-Za-z0-9_-]+`)

	// Cache keychain password to avoid multiple prompts
	keychainCache     string
	keychainCacheOnce sync.Once
	keychainCacheErr  error
)

// Credentials holds the session key and org ID
type Credentials struct {
	SessionKey string
	OrgID      string
}

// GetCredentials retrieves both session key and org ID in a single operation
// to minimize keychain prompts
func GetCredentials() (*Credentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cookiesPath := filepath.Join(homeDir, "Library", "Application Support", "Claude", "Cookies")

	// Get keychain password once
	keychainPassword, err := getKeychainPasswordCached()
	if err != nil {
		return nil, fmt.Errorf("failed to get keychain password: %w", err)
	}

	// Open database once for both cookies
	db, err := sql.Open("sqlite3", cookiesPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("failed to open cookies database: %w", err)
	}
	defer db.Close()

	// Get session key
	sessionKey, err := getCookieValue(db, "sessionKey", keychainPassword, sessionKeyRegex)
	if err != nil {
		return nil, fmt.Errorf("failed to get session key: %w", err)
	}

	// Get org ID
	orgID, err := getCookieValue(db, "lastActiveOrg", keychainPassword, uuidRegex)
	if err != nil {
		return nil, fmt.Errorf("failed to get org ID: %w", err)
	}

	return &Credentials{
		SessionKey: sessionKey,
		OrgID:      orgID,
	}, nil
}

func getCookieValue(db *sql.DB, cookieName, keychainPassword string, pattern *regexp.Regexp) (string, error) {
	var encryptedValue []byte
	err := db.QueryRow(
		"SELECT encrypted_value FROM cookies WHERE name = ? AND host_key LIKE '%claude.ai'",
		cookieName,
	).Scan(&encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to read cookie %s: %w", cookieName, err)
	}

	decrypted, err := decryptChromiumCookie(encryptedValue, keychainPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt cookie %s: %w", cookieName, err)
	}

	match := pattern.FindString(decrypted)
	if match == "" {
		return "", fmt.Errorf("pattern not found in decrypted %s", cookieName)
	}

	return match, nil
}

func getKeychainPasswordCached() (string, error) {
	keychainCacheOnce.Do(func() {
		cmd := exec.Command("security", "find-generic-password", "-s", "Claude Safe Storage", "-a", "Claude Key", "-w")
		output, err := cmd.Output()
		if err != nil {
			keychainCacheErr = fmt.Errorf("failed to get keychain password (click 'Always Allow' when prompted): %w", err)
			return
		}
		keychainCache = strings.TrimSpace(string(output))
	})
	return keychainCache, keychainCacheErr
}

func decryptChromiumCookie(encryptedValue []byte, keychainPassword string) (string, error) {
	if len(encryptedValue) < 3 {
		return "", fmt.Errorf("encrypted value too short")
	}

	if string(encryptedValue[:3]) != "v10" {
		return "", fmt.Errorf("unexpected encryption version: %s", string(encryptedValue[:3]))
	}

	data := encryptedValue[3:]

	if len(data)%16 != 0 {
		return "", fmt.Errorf("data length not multiple of block size")
	}

	key := pbkdf2.Key([]byte(keychainPassword), []byte(chromiumSalt), chromiumIterations, chromiumKeyLength, sha1.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := []byte("                ") // 16 spaces

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	if len(decrypted) > 0 {
		padding := int(decrypted[len(decrypted)-1])
		if padding > 0 && padding <= 16 {
			decrypted = decrypted[:len(decrypted)-padding]
		}
	}

	return string(decrypted), nil
}
