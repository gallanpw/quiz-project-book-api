# sanbercode-golang-batch-69
quiz project book api golang bootcamp sanbercode

# Quiz Project Sanbercode Golang Batch 69
Kelas : Golang\
Batch : 69\
Teknis Pengerjaan : Individu

---

# Proyek API Toko Buku Golang
Proyek ini adalah implementasi RESTful API untuk sistem manajemen toko buku yang dibuat menggunakan bahasa pemrograman Golang. Aplikasi ini menggunakan framework Gin Gonic untuk web server, **PostgreSQL** sebagai database, dan menerapkan **JSON Web Token (JWT)** untuk otentikasi.

## 🚀 Cara Menjalankan Proyek
Ikuti langkah-langkah di bawah ini untuk menginstal dan menjalankan proyek secara lokal.

### Prasyarat

* **Golang**: Pastikan Go sudah terinstal di sistem kamu.
* **PostgreSQL**: Pastikan server PostgreSQL sudah berjalan dan kamu memiliki kredensial akses.
* **Go Modules**: Proyek ini menggunakan Go Modules.

### 1. Klon Repositori
```sh
git clone <URL_REPOSITORI_KAMU>
cd <nama-proyek-golang>
```

### 2. Konfigurasi Database
Buat file `.env` di direktori utama dan isi dengan detail koneksi database kamu.
```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=nama_database
DB_SSLMODE=disable
JWT_SECRET_KEY=kunci_rahasia_jwt_kamu
```

### 3. Buat Skema Database
Gunakan skrip SQL berikut untuk membuat tabel yang diperlukan:
```sql
-- Tabel User
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

-- Tabel Kategori
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
);

-- Tabel Buku
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    image_url VARCHAR(255),
    release_year INTEGER,
    price INTEGER,
    total_page INTEGER,
    thickness VARCHAR(255),
    category_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
);

-- Tabel Blacklist Token (untuk logout)
CREATE TABLE blacklisted_tokens (
    token VARCHAR(255) PRIMARY KEY,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 4. Jalankan Aplikasi
Jalankan perintah berikut untuk mengunduh semua dependencies dan menjalankan server.
```sh
go mod tidy
go run main.go
```
Server local akan berjalan di `http://localhost:8080`.
Server production akan berjalan di `https://quiz-project-book-api-production.up.railway.app`.

---

## 🔑 Autentikasi

Semua endpoint kecuali `/api/users/login` dilindungi oleh **JWT Middleware**. Kamu harus mendapatkan token terlebih dahulu dari endpoint login dan menyertakannya di setiap permintaan melalui header `Authorization` dengan format `Bearer <token>`.

### Endpoint
<!-- ```sh
Method	Path               Deskripsi
POST	/api/users/login   Otentikasi: Mengambil token JWT.
POST	/api/users/logout  Memasukkan token ke daftar hitam (blacklist) untuk logout.
``` -->

| Method | Path                | Deskripsi |
|--------|---------------------|-----------|
| `POST` | `/api/users/login`  | Otentikasi: Mengambil token JWT.|
| `POST` | `/api/users/logout` | Memasukkan token ke daftar hitam (blacklist) untuk logout. |

---

## 📦 Fitur API

### 📚 API Kategori (`/api/categories`)
Mengelola data kategori buku. Semua endpoint memerlukan token JWT.

| Method | Path                          | Deskripsi      |
|--------|-------------------------------|----------------|
| `GET`    | `/api/categories`           | Menampilkan semua kategori yang aktif. |
| `POST`   | `/api/categories`           | Menambahkan kategori baru. |
| `GET`    | `/api/categories/:id`       | Menampilkan detail kategori berdasarkan ID. |
| `PUT`    | `/api/categories/:id`       | Memperbarui data kategori berdasarkan ID. |
| `DELETE` | `/api/categories/:id`       | Soft Delete: Menghapus kategori secara logis. |
| `GET`    | `/api/categories/:id/books` | Menampilkan buku berdasarkan kategori tertentu. |

### 📖 API Buku (`/api/books`)

Mengelola data buku. Semua endpoint memerlukan token JWT.

| Method | Path           | Deskripsi      |
|--------|----------------|----------------|
| `GET`    | `/api/books`     | Menampilkan semua buku yang aktif. |
| `POST`   | `/api/books`     | Menambahkan buku baru. |
| `GET`    | `/api/books/:id` | Menampilkan detail buku berdasarkan ID. |
| `PUT`    | `/api/books/:id` | Memperbarui data buku berdasarkan ID. |
| `DELETE` | `/api/books/:id` | **Soft Delete**: Menghapus buku secara logis. |

### Aturan Validasi

* `release_year` pada buku dibatasi antara **1980** dan **2024**.
* `thickness` dihitung berdasarkan `total_page`:
    * `total_page > 100`: `thickness` diisi **"tebal"**
    * `total_page <= 100`: `thickness` diisi **"tipis"**