# Go Wallet App

Go ile geliştirilmiş örnek **Wallet (cüzdan)** uygulaması.  
Proje; **REST API (Gin)** ve **gRPC** arayüzleri üzerinden cüzdan oluşturma, bakiye sorgulama, para yatırma/çekme işlemleri sağlar.  
Veri katmanında **Ent ORM** ve **PostgreSQL** kullanılır, migration yönetimi için **Atlas** entegre edilmiştir.  
Kimlik doğrulama için **JWT Authorization** uygulanmıştır.

---

## 🚀 Özellikler

- Kullanıcı kaydı (signup) ile otomatik cüzdan oluşturma
- JWT tabanlı kimlik doğrulama (login sonrası token ile erişim)
- Para yatırma / çekme (validation ve hata yönetimi ile)
- Bakiye ve işlem dökümü sorgulama
- REST + gRPC API (aynı iş kurallarını kullanır)
- Eşzamanlı işlem güvenliği (mutex / transaction)
- Test senaryoları (unit + integration)

---

## 📂 Proje Yapısı

├─ cmd/server/main.go # Uygulama giriş noktası
├─ internal/
│ ├─ http/ # Gin router, REST handler
│ ├─ service/ # İş kuralları (business logic)
│ ├─ repo/ # Ent repository katmanı
│ ├─ db/ # DB init, Ent client, migrations
│ ├─ config/ # Viper config
│ └─ validation/ # Request validation
├─ ent/ # Ent codegen çıktıları
├─ migrations/ # Atlas migration dosyaları
├─ proto/ # gRPC proto ve üretilmiş kodlar
├─ api/ # Swagger dokümantasyonu
├─ tests/ # Integration & unit testler

├─ docker-compose.yml # PostgreSQL + App için docker
└─ Makefile
---

## ⚙️ Kurulum

### 1. Gereksinimler
- Go `>=1.22`
- Docker & Docker Compose
- Atlas CLI (migration için)
- Protoc (gRPC için)

### 2. Kurulum
```bash
git clone https://github.com/<username>/go-wallet-app.git
cd go-wallet-app
go mod tidy
````
###3. Docker ile DB başlat
```bash
docker-compose up -d
````

### 4. Migration çalıştır
```bash
atlas migrate apply --dir "file://migrations" --url "postgres://user:pass@localhost:5432/wallet?sslmode=disable"
````
###5. Server başlat
```bash
go run ./cmd/server/main.go
````
## 🌐 API Kullanımı
## 🔑 Authorization (JWT)

Uygulamada kimlik doğrulama ve yetkilendirme için **JWT (JSON Web Token)** kullanılmaktadır.  

- **POST /signup**  
  Yeni kullanıcı kaydı yapılır.  
  → Kullanıcı kaydı tamamlandığında cüzdan da otomatik olarak oluşturulur.  

- **POST /login**  
  Kullanıcı giriş yapar.  
  → Dönen JWT token, sonraki tüm `wallet` endpoint’lerinde **Authorization** header’ı ile gönderilmelidir:  



### REST (Gin)

- **POST /signup** -> kayıt oluşturma, kayıt ile beraber wallet de oluşur
- -**POST /login** -> login yapılır 
- **POST /wallet/deposit** → para yatır  
- **POST /wallet/withdraw** → para çek  
- **GET /wallet/balance** → bakiye sorgula  
- **GET /wallet/transaction** -> transaction dökümü
#### Örnek

```bash
# Kayıt ol
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Rabia","email":"test@example.com","password":"123456"}'

# Login ol
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Para yatır (Authorization header ile)
curl -X POST http://localhost:8080/wallet/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"amount": 50}'

````






gRPC

WalletService.GetBalance

WalletService.Deposit

WalletService.Withdraw

Örnek (evans CLI)
evans --host localhost --port 50051 -r repl
> call GetBalance
userId (TYPE_STRING) => <UUID>

✅ Testler

Testler tests/ klasöründe yer alıyor.
Çalıştırmak için:

go test ./... -v

Örnek Senaryolar

Bakiye 100 TL, çekim 150 TL → Hata beklenir

Bakiye 100 TL, eşzamanlı 2×60 TL çekim → Yalnızca biri başarılı olmalı

REST ile para yatır, gRPC ile bakiye sorgula → Aynı sonucu dönmeli

🔧 Teknolojiler

Gin
 → REST framework

gRPC
 → yüksek performanslı RPC

Ent
 → Go ORM

Atlas
 → DB migration tool

PostgreSQL
 → Veritabanı

Docker
 → container
