# 🎓 Aksara Batak API

API untuk project lomba pembelajaran aksara batak
---

## 📁 Struktur Project

```bash
.
├── apispec.json                  # Dokumentasi Swagger API
├── config/
│   └── env.go                    # Konfigurasi environment (load API key, token )
├── controllers/
│   └── controller.go      # Handler Gin untuk setiap endpoint 
├── main.go                       # Entry point aplikasi
├── middleware/
│   └── auth_middleware.go        # Middleware autentikasi API Key
├── model/
│   └── web/
│       ├── request.go
│       ├── response.go
│       └── web_response.go       # Struktur umum response API
|       └── dst.....
├── README.md                     # Dokumentasi proyek
├── routes/
│   └── routes.go                 # Definisi semua endpoint dan group route Gin
├── services/
│   ├── service.go         # Interface untuk Service
│   └── service_impl.go    # Implementasi business logic Service
├── repositories/
│   ├── repository.go         # Interface untuk repository
│   └── repository_impl.go    # implementasi untuk database ORM

````

---

## ⚙️ Requirement

* Go versi **1.24+**
* Golang Migrate

---

## 🛠️ Instalasi

### 1. Clone Project

```bash
git clone https://github.com/rizkycahyono97/aksara_batak_api
cd aksara_batak_api
```

### 2. Buat File `.env`

Contoh `cp .env.example .env`:

```env
MOODLE_TOKEN=your_moodle_token_here
APP_PORT=YOUR_PORT
JWT_SECRET_KEY=your_custom_jwt_key
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Jalankan Aplikasi

```bash
go run main.go
```

---

## 📦 Dependency Penting

* [`https://docs.gofiber.io/`](https://docs.gofiber.io/) - HTTP Web Framework
* [`github.com/joho/godotenv`](https://github.com/joho/godotenv) - Load environment variables
* [`github.com/swaggo/gin-swagger`](https://github.com/swaggo/gin-swagger) - Swagger UI untuk dokumentasi API

---

## 📌 Endpoint Penting

| Method | Endpoint                             | Deskripsi         |
| ------ |--------------------------------------|-------------------|
| `GET`  | `/api/v1/login`                      | Login Aplikasi    |
| `POST` | `/api/v1/register`                   | Register Aplikasi |
---
