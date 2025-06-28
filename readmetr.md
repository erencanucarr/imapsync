# IMAPSYNC CLI

![Go Sürümü](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Lisans](https://img.shields.io/badge/Lisans-MIT-green)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)

> Posta kutusu taşıma işini **güvenli, izlenebilir ve çoklu platform** hâline getiren, [imapsync](https://imapsync.lamiral.info/) etrafında çok dilli bir CLI sarmalayıcı.

---

## ✨ Öne Çıkanlar

| Kategori | Detaylar |
|----------|----------|
| Etkileşim | İngilizce, Türkçe, İspanyolca, Almanca (genişletilebilir) kılavuzlu menüler |
| Güvenlik | Şifreler gizli giriş ile alınır (`golang.org/x/term`) |
| Taşınabilirlik | Ubuntu, Debian, CentOS/RHEL, Arch, macOS için hazır kurulum betikleri |
| İzlenebilirlik | `imapsync` çıktısını ayrıştıran gerçek zamanlı ilerleme çubuğu (`schollz/progressbar`) |
| Genişletilebilirlik | Modüler Go kodu, JSON i18n, temiz klasör düzeni |

---

## 📂 Depo Yapısı

```
├── cmd/               # CLI giriş noktası
│   └── imapsync/      # main.go
├── internal/
│   ├── app/           # iş mantığı (kurulum, transfer, geliştirici)
│   ├── i18n/          # yerelleştirme yükleyici + JSON dosyaları (en, tr, es, de)
│   └── ui/            # renk yardımcıları
├── install/           # İşletim sistemi özel imapsync kurulum betikleri
├── README.md          # İngilizce sürüm
└── readmetr.md        # Türkçe sürüm (buradasınız)
```

---

## 🚀 Hızlı Başlangıç

### 1. Derle & Çalıştır

```bash
git clone https://github.com/yourname/imapsync-cli.git
cd imapsync-cli

go build -o imapsync ./cmd/imapsync
./imapsync -lang=tr   # dili -lang=<kod> ile değiştirebilirsiniz
```

### 2. İlk Kurulum

Menüden **1) Sistem Kurulumu** seçin. Python / imapsync eksikse şöyle bir ekran görürsünüz:

```
Python ✗
imapsync ✗
Yerel kurulum betikleri ./install klasöründe:
 - ubuntu  - debian  - centos  - arch  - darwin
Kurulum için dağıtım anahtarını girin veya Enter ile atlayın:
```

Dağıtım anahtarını yazarak **sudo/root** ile çalışan betiği başlatın. Kayıtlar `/var/log/imapsync-install.log` dosyasına yazılır.

### 3. E-posta Taşıma

1. **2) Posta Transferi**’ni seçin.
2. Kaynak & hedef sunucu/e-posta/şifre bilgilerini girin (şifre gizli).
3. Program `imapsync --justlogin` ile oturum doğrular, sonra senkronizasyonu başlatır.
4. İlerleme çubuğunu izleyin – `Ctrl+C` ile güvenle durdurup `--useuid` sayesinde kaldığınız yerden devam edebilirsiniz.

---

## ⚙️ Önerilen imapsync Bayrakları

Uygulama üretimde kanıtlanmış varsayılanlarla imapsync çağırır:

```
--ssl1 --ssl2 \
--exclude "^Junk E-Mail" --exclude "^Trash" --exclude "^Deleted( Items)?$" \
--regextrans2 's#^Sent$#Sent Items#' --regextrans2 's#^Spam$#Junk E-Mail#' \
--useuid --usecache --tmpdir ./tmp --syncinternaldates --progress
```
`internal/app/transfer.go` içinde ihtiyaçlarınıza göre düzenleyebilirsiniz.

---

## 🖥️ Manuel imapsync Kurulumu

Kendi paket yöneticinizi tercih ediyorsanız `install/` betiklerine veya resmi dokümana bakın:
<https://imapsync.lamiral.info/INSTALL.d/>.

---

## 🛠️ Geliştirme

```bash
go vet ./...
go test ./...
```

Hızlı deneme için `go run ./cmd/imapsync -lang=tr` kullanın.

### Yeni Dil Ekleme

1. `internal/i18n/locales/en.json` dosyasını kopyalayıp `internal/i18n/locales/fr.json` gibi adlandırın.
2. Değerleri çevirin.
3. `-lang=fr` ile çalıştırın.

### Kurulum Betiği Genişletme

`install/` klasörüne `<distro>.txt` ekleyin ve `internal/app/setup.go` içindeki `scripts` haritasına anahtarını tanımlayın.

---

## 🙋‍♂️ SSS

* **Şifreler saklanıyor mu?**  Hayır, doğrudan `imapsync` sürecine argüman olarak iletilir.
* **Yarım kalan senk. devam eder mi?**  Evet – `--useuid` sayesinde tekrarlanabilir.
* **Grafik arayüz?**  Yol haritasında; katkılara açığız!

---

## 🤝 Katkı

1. Depoyu çatallayın, dal açın.
2. Go standartlarına uyun (`go vet`, `golint`).
3. PR gönderin; CI `go test` çalıştırır.

Çeviri, hata bildirimi ve geliştirmeler memnuniyetle karşılanır.

---

## 📅 Yol Haritası

- [ ] Yapılandırma dosyası / profil kaydetme
- [ ] OAuth2 desteği (Gmail, Outlook 365)
- [ ] BubbleTea ile TUI
- [ ] Derleme hattı & ikili sürümler

---

## 📝 Lisans

MIT © 2025 Your Name
