//go:build legacy_i18n
// +build legacy_i18n

package i18n

// T returns the message for the given key and language.
// If the key is missing for that language, it falls back to English.
func T(lang, key string) string {
    if m, ok := messages[lang]; ok {
        if v, ok2 := m[key]; ok2 {
            return v
        }
    }
    return messages["en"][key]
}

var messages = map[string]map[string]string{
    "en": {
        "title":          "IMAPSYNC Mail Transfer Tool",
        "menu":           "Please select an option:",
        "menu_setup":     "1 - Setup System",
        "menu_transfer":  "2 - Transfer Mail",
        "menu_developer": "3 - Developer",
        "menu_exit":      "4 - Exit",
        "choice":         "Choice (1/2/3/4): ",
        "invalid":        "Invalid choice. Please enter 1, 2, 3 or 4.",
        "exit":           "Exiting program...",
        "transfer_start": "Starting mail transfer...",
        "transfer_success": "Mail transfer completed successfully!",
        "transfer_fail": "Mail transfer failed.",
        "error": "Error:",
    },
    "tr": {
        "title":          "IMAPSYNC Mail Taşıma Aracı",
        "menu":           "Lütfen bir seçenek seçin:",
        "menu_setup":     "1 - Sistemi Kur",
        "menu_transfer":  "2 - Mail Taşı",
        "menu_developer": "3 - Geliştirici",
        "menu_exit":      "4 - Çıkış",
        "choice":         "Seçiminiz (1/2/3/4): ",
        "invalid":        "Geçersiz seçenek. Lütfen 1, 2, 3 veya 4 girin.",
        "exit":           "Programdan çıkılıyor...",
        "transfer_start": "Mail transferi başlatılıyor...",
        "transfer_success": "Mail transferi başarıyla tamamlandı!",
        "transfer_fail": "Mail transferi başarısız.",
        "error": "Hata:",
    },
}
