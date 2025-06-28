package i18n

import (
    "embed"
    "encoding/json"
    "log"
    "path/filepath"
    "strings"
)

//go:embed locales/*.json
var localeFS embed.FS

var messages = map[string]map[string]string{}

func init() {
    entries, err := localeFS.ReadDir("locales")
    if err != nil {
        log.Printf("i18n: cannot read embedded locales: %v", err)
        return
    }
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        data, err := localeFS.ReadFile(filepath.Join("locales", e.Name()))
        if err != nil {
            log.Printf("i18n: read file %s: %v", e.Name(), err)
            continue
        }
        var m map[string]string
        if err := json.Unmarshal(data, &m); err != nil {
            log.Printf("i18n: parse %s: %v", e.Name(), err)
            continue
        }
        lang := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
        messages[lang] = m
    }
}

// T returns translated message or fallback to English or key itself.
func T(lang, key string) string {
    if m, ok := messages[lang]; ok {
        if v, ok2 := m[key]; ok2 {
            return v
        }
    }
    if m, ok := messages["en"]; ok {
        if v, ok2 := m[key]; ok2 {
            return v
        }
    }
    return key
}
