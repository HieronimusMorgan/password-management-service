# 🔐 Password Management Service

A secure, end-to-end encrypted password manager built with Golang. This service allows users to manage and share credentials securely using AES-256 and RSA public-key encryption.

---

## ✨ Features

* **End-to-End Encryption** with AES-GCM (256-bit)
* **RSA Key Pair per User** for encrypted private key storage
* **Secure Password Sharing** using key wrapping
* **Password Grouping** for organizing entries
* **Tag Filtering** on entries
* **Password History Logging** for audit trails

---

## 🧱 Tech Stack

* **Golang 1.21+**
* **PostgreSQL** (with native `text[]` support for tags)
* **GORM** ORM
* **Gin** HTTP framework
* **Argon2id** for key derivation
* **RSA-2048**, **AES-GCM-256** for encryption

---

## 🗃 Database Schema

Key tables:

* `users` – core user info, hashed credentials
* `user_keys` – per-user RSA key pair
* `password_entries` – encrypted password entries
* `password_entry_keys` – wrapped AES key
* `shared_passwords` – encrypted shared access
* `password_groups` – groupings of entries
* `password_history` – historical changes

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/your-username/password-management-service.git
cd password-management-service
```

### 2. Setup environment

Create a `.env` file with the following:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=passman
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run migrations

Use GORM auto migration or embed migrations:

```go
// In main.go or setup:
db.AutoMigrate(&models.User{}, &models.UserKey{}, &models.PasswordEntry{}, ...)
```

### 5. Run the server

```bash
go run main.go
```

---

## 🔒 Encryption Flow

1. **RSA key pair** is generated per user
2. Private key is encrypted with AES key derived from user password + salt (Argon2id)
3. For each password entry:

   * AES key encrypts the password/notes
   * AES key is encrypted using RSA public key
4. Decryption requires:

   * Unwrapping AES key with RSA private key
   * Decrypting data with AES-GCM

---


## 👥 Contributing

PRs and suggestions welcome! Please open issues for bugs or feature requests.

---


## 📫 Contact

Created by \[Hieronimus Fredy Morgan] – feel free to reach out on GitHub!
