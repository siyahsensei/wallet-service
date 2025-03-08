# Wallet Service - Finansal Varlık Yönetim Uygulama Servisleri

Bu uygulama, farklı finansal varlıklarınızı (banka hesapları, borsa yatırımları, kripto paralar vb.) tek bir platformda yönetmenize olanak tanır.

## Özellikler

- **Nakit Varlıklar**: Banka hesapları, nakit para
- **Yatırım Varlıkları**: Hisse senetleri, yatırım fonları, tahviller, VİOP
- **Kripto Varlıklar**: Kripto para birimleri, NFT'ler, DeFi varlıkları
- **Diğer Varlıklar**: Altın/gümüş, gayrimenkul, borçlar/alacaklar
- **Varlık Takibi**: Tüm varlıklarınızın toplam değerini ve zaman içindeki performansını görüntüleme
- **İşlem Kaydı**: Tüm finansal işlemlerinizin kaydını tutma
- **API Entegrasyonu**: Banka, borsa ve kripto borsası API'leri ile otomatik veri senkronizasyonu

## Teknolojiler

- **Backend**: Go (Golang)
- **Web Framework**: Fiber
- **Veritabanı**: PostgreSQL
- **Kimlik Doğrulama**: JWT

## Kurulum

### Ön Koşullar

- Go 1.21+
- PostgreSQL
- git

### Veritabanı Kurulumu

```bash
# PostgreSQL veritabanı oluşturma
createdb wallet

# Migrasyonları çalıştırmak için golang-migrate kurulumu
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Migrasyonları çalıştırma
make migrate-up
```

### Uygulamayı Çalıştırma

1. Repo'yu klonlayın:

```bash
git clone https://github.com/siyahsensei/wallet-service.git
cd wallet-service
```

2. Bağımlılıkları yükleyin:

```bash
go mod download
```

3. Konfigürasyon için `.env` dosyasını oluşturun:

```bash
cp .env.example .env
# .env dosyasını düzenleyerek gerekli ayarları yapın
```

4. API Sunucusunu Çalıştırın:

```bash
make run-api
# veya doğrudan:
go run cmd/api/main.go
```

5. Worker Servisini Çalıştırın (opsiyonel, arka plan görevleri için):

```bash
make run-worker
# veya doğrudan:
go run cmd/worker/main.go
```

## API Endpoint'leri

### Kimlik Doğrulama

- `POST /api/v1/auth/register` - Yeni kullanıcı kaydı
- `POST /api/v1/auth/login` - Kullanıcı girişi
- `GET /api/v1/auth/me` - Mevcut kullanıcı bilgilerini görüntüleme

### Hesaplar

- `GET /api/v1/accounts` - Tüm hesapları listele
- `POST /api/v1/accounts` - Yeni hesap oluştur
- `GET /api/v1/accounts/{id}` - Belirli bir hesabı görüntüle
- `PUT /api/v1/accounts/{id}` - Hesap bilgilerini güncelle
- `DELETE /api/v1/accounts/{id}` - Hesabı sil
- `POST /api/v1/accounts/{id}/credentials` - Hesap API kimlik bilgilerini ayarla
- `GET /api/v1/accounts/types` - Mevcut hesap tiplerini listele

### Varlıklar

- `GET /api/v1/assets` - Tüm varlıkları listele
- `POST /api/v1/assets` - Yeni varlık ekle
- `GET /api/v1/assets/{id}` - Belirli bir varlığı görüntüle
- `PUT /api/v1/assets/{id}` - Varlık bilgilerini güncelle
- `DELETE /api/v1/assets/{id}` - Varlığı sil
- `GET /api/v1/assets/types` - Mevcut varlık tiplerini listele

### İşlemler

- `GET /api/v1/transactions` - Tüm işlemleri listele
- `POST /api/v1/transactions` - Yeni işlem ekle
- `GET /api/v1/transactions/{id}` - Belirli bir işlemi görüntüle
- `PUT /api/v1/transactions/{id}` - İşlem bilgilerini güncelle
- `DELETE /api/v1/transactions/{id}` - İşlemi sil
- `GET /api/v1/transactions/types` - Mevcut işlem tiplerini listele

## Proje Yapısı

```
wallet-service/
├── cmd/                # Ana çalıştırılabilir dosyalar (main paketleri)
│   ├── api/            # API sunucusu (REST, gRPC vb.)
│   │   └── main.go
│   └── worker/         # Arka plan görevlerini işleyen worker'lar (örn. periyodik veri senkronizasyonu) (Daha tam planlanmadı. Gelecekte...)
│       └── main.go
├── internal/           # Uygulamaya özel, dışarıdan erişilemeyen (private) kodlar
│   ├── app/            # Uygulama katmanı (iş mantığı servisleri)
│   │   ├── users/       # Kullanıcı yönetimi servisi
│   │   ├── accounts/    # Hesap yönetimi servisi
│   │   ├── assets/      # Varlık yönetimi servisi
│   │   ├── transactions/ # İşlem yönetimi servisi
│   │   └── ...
│   ├── pkg/            # Uygulamanın farklı bölümleri tarafından kullanılabilecek yardımcı paketler (utilities)
│   │   ├── auth/       # Kimlik doğrulama ve yetkilendirme
│   │   ├── httpclient/  # HTTP client konfigürasyonu ve yardımcı fonksiyonlar
│   │   ├── logger/     # Loglama
│   │   ├── config/     # Konfigürasyon yönetimi
│   │   └── ...
│   ├── platform/       # Dış servislerle entegrasyonlar (3. parti API'ler, veritabanı)
│   │   ├── database/   # Veritabanı bağlantısı ve işlemleri (ORM kullanılıyorsa, model tanımları burada olabilir)
│   │   ├── bankapi/   # Banka API'leri ile entegrasyon
│   │   ├── exchangeapi/ # Borsa API'leri ile entegrasyon
│   │   ├── cryptoapi/  # Kripto para borsası API'leri ile entegrasyon
│   │   └── ...
├── pkg/                # Diğer projelerde de kullanılabilecek, genel amaçlı (public) paketler (opsiyonel)
│   ├── api/            # API tanımları (Protobuf, OpenAPI vb.)
│   └── ...
├── domain/             # İş alanı (domain) nesneleri ve kuralları (DDD)
│   ├── user/          # Kullanıcı modeli ve ilgili iş kuralları
│   │   ├── user.go
│   │   ├── repository.go  # Kullanıcı verilerine erişim arayüzü (interface)
│   │   └── service.go     # Kullanıcı ile ilgili iş mantığı (opsiyonel, `internal/app` ile birleştirilebilir)
│   ├── account/       # Hesap modeli ve ilgili iş kuralları
│   ├── asset/         # Varlık modeli ve ilgili iş kuralları
│   ├── transaction/  # İşlem modeli ve ilgili iş kuralları
│   └── ...
├── infrastructure/   # Altyapı katmanı (DDD) - Veritabanı, dış servisler vb. ile etkileşim (opsiyonel)
│   ├── persistence/ # Veritabanı işlemleri (Repository implementasyonları)
│   │   ├── userrepo/   # UserRepository implementasyonu
│   │   ├── accountrepo/
│   │   └── ...
│   ├── external/      # Dış servislerle etkileşim
│   │   ├── bank/      # Banka API client
│   │   ├── exchange/  # Borsa API client
│   │   └── ...
├── api/           # Sunum katmanı (Presentation Layer) - API handler'ları, istek/yanıt modelleri
│  ├── handlers/  # HTTP handler'ları
│  ├── models/   # İstek/yanıt (request/response) veri yapıları
│  └── middleware/ # Middleware'ler (kimlik doğrulama, loglama vb.)
├── scripts/            # Yardımcı script'ler (veritabanı migration'ları, deployment vb.)
├── deployments/        # Deployment konfigürasyonları (Docker, Kubernetes vb.)
├── configs/            # Konfigürasyon dosyaları (ortama özel ayarlar)
├── test/               # Testler (unit, integration, e2e)
├── Makefile            # Yaygın görevler için kısayollar (build, test, deploy vb.)
└── go.mod, go.sum    # Go modülleri
```

## Geliştirme

### Yeni Migration Oluşturma

```bash
make migrate-create name=migration_ismi
```

### Test Çalıştırma

```bash
make test
```

### Uygulamayı Derleme

```bash
make build
``` 