# Audio-Steganografi

# Deskripsi Program
Program ini berisi implementassi steganografi pada file audio .mp3 dengan algoritma n-lsb. Program ini juga dilengkapi dengan fitur enkripsi Extended Vigenere Chiper. Untuk mendukung kemudahan pemakaian program ini, disediakan interface berupa website sederhana. Program dapat melakukan penyisipan dengan cara memberikan cover audio, konfigurasi steganografi (Enkripsi atau tidak, Random position, jumlah bit, Key), kemudian hasilnya berupa kalkulasi PNSR dan audio hasil steganografi. Audio sebelum dan sesudah proses steganografi dapat diputar di website. File stego audio dapat diunduh oleh user. Untuk proses ekstraksi stego audio, pertama user akan memasukan file stego audio. Kemudian akan diminta key dan pengaturan random position diaktifkan atau tidak. Hasilnya berupa secret message yang dapat diunduh.

# Teknologi dan Dependensi
- React
- Typescript
- Axios
- Shadcn UI
- Tailwind CSS
- Golang
- GIN
- Node Js
- air (untuk hot reload)

# Tata Cara Menjalankan
- Pastikan anda sudah menginstall Node Js, Golang, Package manager (kami menggunakan pnpm) 
- Untuk menjalankan frontend, pertama lakukan instalasi semua package dengan cara ```pnpm i```
- Kemudian jalankan frontend dengan cara ```pnpm dev```
- Kemudian lakukan instalasi semua package golang dengan cara ```go mod download```
- Jalankan proyek dengan memanggil ```air```
- Pastikan frontend berjalan di localhost:5173 dan backend di localhost:8080
