# Zero Dependency IMAPSYNC CLI

Bu proje **zero dependency** prensibiyle geliştirilmiştir. Tüm harici kütüphaneler kendi implementasyonlarımızla değiştirilmiştir.

## 🎯 Zero Dependency Hedefi

Proje, herhangi bir harici Go kütüphanesine bağımlı olmadan çalışır. Sadece Go'nun standart kütüphanesi kullanılır.

## 📦 Değiştirilen Kütüphaneler

### 1. **github.com/schollz/progressbar/v3** → `internal/app/progressbar.go`
- **Özellik**: Terminal progress bar
- **Kendi Implementasyonumuz**: 
  - Basit progress bar
  - ETA hesaplama
  - Performans optimizasyonu (100ms güncelleme aralığı)

### 2. **github.com/patrickmn/go-cache** → `internal/app/cache.go`
- **Özellik**: In-memory cache with expiration
- **Kendi Implementasyonumuz**:
  - Thread-safe cache
  - Otomatik expiration
  - Cleanup mekanizması

### 3. **github.com/sirupsen/logrus** → `internal/app/logger.go`
- **Özellik**: Structured logging
- **Kendi Implementasyonumuz**:
  - Basit log seviyeleri (DEBUG, INFO, WARN, ERROR, FATAL)
  - Timestamp formatı
  - Configurable output

### 4. **golang.org/x/sync/semaphore** → `internal/app/semaphore.go`
- **Özellik**: Semaphore for concurrency control
- **Kendi Implementasyonumuz**:
  - Weighted semaphore
  - Context cancellation support
  - TryAcquire functionality

### 5. **golang.org/x/term** → `internal/app/term.go`
- **Özellik**: Terminal input handling
- **Kendi Implementasyonumuz**:
  - Cross-platform password input
  - Basit terminal handling

## 🚀 Performans Optimizasyonları

### Paralel İşleme
- **Bağlantı Havuzu**: Semaphore ile eşzamanlı transfer sınırlaması
- **Cache Sistemi**: Başarılı transferlerin önbelleklenmesi
- **Bellek Yönetimi**: Otomatik bellek optimizasyonu

### İstatistikler
- Transfer başarı oranı
- Ortalama hız hesaplama
- Bellek kullanımı takibi
- Cache performansı

## 📊 Mevcut Özellikler

### 1. **Temel Transfer**
- Tek hesap transferi
- Progress bar ile ilerleme takibi
- Retry mekanizması

### 2. **Paralel Transfer**
- Birden fazla hesabı aynı anda transfer etme
- Job yönetimi
- Durum takibi

### 3. **Performans İstatistikleri**
- Gerçek zamanlı metrikler
- Bellek kullanımı
- Cache durumu

### 4. **Sistem Kurulumu**
- Otomatik imapsync kurulumu
- Platform desteği

## 🔧 Teknik Detaylar

### Cache Implementasyonu
```go
type Cache struct {
    items map[string]*CacheItem
    mu    sync.RWMutex
}

type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}
```

### Progress Bar Implementasyonu
```go
type ProgressBar struct {
    current     int
    total       int
    width       int
    description string
    startTime   time.Time
    lastUpdate  time.Time
}
```

### Semaphore Implementasyonu
```go
type Semaphore struct {
    permits int64
    mu      sync.Mutex
    cond    *sync.Cond
}
```

## 🎯 Avantajlar

1. **Bağımsızlık**: Harici kütüphane güncellemelerinden etkilenmez
2. **Güvenlik**: Sadece kendi kodumuz çalışır
3. **Performans**: Minimal overhead
4. **Bakım**: Tüm kod elimizde
5. **Özelleştirme**: İhtiyaçlara göre uyarlanabilir

## 📈 Gelecek Planları

- [ ] Gelişmiş terminal UI
- [ ] Web arayüzü
- [ ] Konfigürasyon dosyası desteği
- [ ] OAuth2 entegrasyonu
- [ ] Backup ve restore özellikleri

## 🛠️ Geliştirme

Proje tamamen Go standart kütüphanesi kullanır:

```bash
go build -o imapsync ./cmd/imapsync
./imapsync -lang=tr
```

## 📝 Lisans

MIT © 2025 - Zero Dependency IMAPSYNC CLI 