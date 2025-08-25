# Tahap build: Menggunakan image Go resmi untuk mengompilasi aplikasi
FROM golang:1.24-alpine AS builder

# Atur direktori kerja
WORKDIR /app

# Salin go.mod dan go.sum untuk mengunduh dependensi
COPY go.mod ./
COPY go.sum ./

# Unduh semua dependensi
RUN go mod download

# Salin semua file sumber
COPY . .

# Bangun (compile) aplikasi ke file biner bernama 'main'
RUN go build -o main .

# Tahap final: Menggunakan image yang sangat kecil untuk menjalankan aplikasi
FROM alpine:3.18

# Mengatur variabel lingkungan untuk port
ENV PORT=8080

# Salin file biner dari tahap builder
COPY --from=builder /app/main .

# Atur direktori kerja
WORKDIR /

# Jalankan aplikasi
CMD ["./main"]
