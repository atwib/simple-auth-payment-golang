package utils

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var (
	blacklistMutex sync.Mutex
	fileData       = "data/blacklist.json"
)

type BlacklistEntry struct {
	Token         string    `json:"token"`
	BlacklistedAt time.Time `json:"blacklisted_at"`
}

// AddTokenToBlacklist menambahkan token ke file blacklist.json
func AddTokenToBlacklist(token string) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Baca file blacklist.json jika ada
	var blacklist []BlacklistEntry
	if _, err := os.Stat(fileData); errors.Is(err, os.ErrNotExist) {
		// Jika file tidak ada, inisialisasi blacklist kosong
		blacklist = []BlacklistEntry{}
	} else {
		file, err := os.ReadFile(fileData)
		if err != nil {
			return err
		}
		// Unmarshal isi file jika file valid
		if err := json.Unmarshal(file, &blacklist); err != nil {
			return err
		}
	}

	// Tambahkan token baru
	blacklist = append(blacklist, BlacklistEntry{
		Token:         token,
		BlacklistedAt: time.Now(),
	})

	// Marshal data kembali dan tulis ke file
	data, err := json.MarshalIndent(blacklist, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileData, data, 0644)
}

func IsTokenBlacklisted(token string) (bool, error) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	// Baca file blacklist.json jika ada
	file, err := os.ReadFile(fileData)
	if err != nil {
		if os.IsNotExist(err) {
			// Jika file tidak ada, berarti token tidak di-blacklist
			return false, nil
		}
		return false, err
	}

	// Parse isi file JSON
	var blacklist []BlacklistEntry
	if err := json.Unmarshal(file, &blacklist); err != nil {
		return false, err
	}

	// Periksa apakah token ada dalam daftar blacklist
	for _, entry := range blacklist {
		if entry.Token == token {
			return true, nil
		}
	}

	return false, nil
}
