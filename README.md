# Go Fiber + Next.js Docker Starter Kit

> Full-stack dev environment in one command: `docker compose up -d --build`

Go Fiber API · GORM · MySQL · Next.js (App Router, TypeScript, Tailwind) · Air hot reload · phpMyAdmin

Starter kit untuk development berbasis **Go Fiber + GORM + MySQL** (API) dan
**Next.js** (App Router, TypeScript, Tailwind) sebagai frontend. Sudah termasuk hot
reload (Air untuk Go, `next dev` untuk Next.js) dan phpMyAdmin.

Hasil konversi dari _Gin MVC Docker Starter Kit_ (Gin + template HTML + Knockout.js).

## Daftar Isi

- [Prasyarat](#prasyarat)
- [Struktur Proyek](#struktur-proyek)
- [Memulai](#memulai)
- [Perintah Sehari-hari](#perintah-sehari-hari)
- [Menambah Modul Baru](#menambah-modul-baru)
- [Contoh Modul Tanpa Database: Test Algoritma](#contoh-modul-tanpa-database-test-algoritma)
- [Aturan Penulisan Go di Proyek Ini](#aturan-penulisan-go-di-proyek-ini)
- [Troubleshooting](#troubleshooting)
- [Padanan Gin/Knockout ke Fiber/Next.js](#padanan-ginknockout-ke-fibernextjs)
- [Permission dan Keamanan](#permission-dan-keamanan)
- [Catatan](#catatan)
- [Credits](#credits)

---

## Prasyarat

- Docker Desktop (atau Docker Engine) dengan **Docker Compose v2** (perintah `docker compose`, bukan `docker-compose`).
- Port host yang bebas untuk web, API, dan phpMyAdmin (nilai default ada di tabel [Memulai](#memulai)).
- Tidak perlu memasang Go atau Node.js di komputer; semuanya berjalan di dalam container.

---

## Struktur Proyek

```
.
├── docker-compose.yml       <- satu-satunya file compose
├── .env.example             <- salin menjadi .env
├── rebuild.bat
├── backend/                 Go Fiber API
│   ├── Dockerfile, entrypoint.sh, .air.toml, go.mod
│   ├── main.go
│   ├── config/              env.go, database.go, auth.go
│   ├── routes/              api.go           (URL -> controller, tanpa logika)
│   ├── controllers/         units_controller.go, auth_controller.go, helpers.go
│   ├── services/            logika bisnis: units_service.go, id.go, seed_units.go
│   ├── models/              struct tabel GORM
│   ├── resources/           bentuk respons JSON
│   ├── middleware/          cors.go
│   └── seed_units.sql       data dummy (opsional, via phpMyAdmin)
└── frontend/                Next.js
    ├── Dockerfile, package.json, next.config.mjs
    ├── app/                 layout.tsx, page.tsx, master-data/units/page.tsx
    ├── components/          Navbar, Pagination
    ├── context/             AppContext.tsx (auth, permission, loading, toast)
    ├── features/units/      UnitsPage.tsx, api.ts
    └── lib/                 api.ts (fetch wrapper), types.ts
```

Alur request:

```
Browser -> Next.js (:3100) --/api/* rewrite--> Fiber (:8000 di dalam container)
                                               Route -> Controller -> Service -> Model (GORM) -> MySQL
```

Browser hanya berbicara dengan Next.js. Request `/api/*` dan `/health` diteruskan
ke Fiber lewat `rewrites()` di `frontend/next.config.mjs`, sehingga **tidak perlu CORS**.

Pembagian tanggung jawab tiap lapisan:

| Lapisan    | Folder         | Tugas                                                        |
| ---------- | -------------- | ------------------------------------------------------------ |
| Route      | `routes/`      | Memetakan URL ke controller. Tidak berisi logika.            |
| Controller | `controllers/` | Membaca request, validasi, memanggil service, mengirim JSON. |
| Service    | `services/`    | Logika bisnis dan fungsi murni. Tidak mengenal HTTP.         |
| Model      | `models/`      | Struct tabel GORM.                                           |

---

## Memulai

**Linux / macOS**

```bash
cp .env.example .env
docker compose up -d --build
docker compose logs -f backend frontend
```

**Windows (PowerShell)**

```powershell
copy .env.example .env
docker compose up -d --build
docker compose logs -f backend frontend
```

Build pertama lebih lama karena mengunduh dependency Go dan menjalankan `npm install`.
Tunggu sampai muncul `listening on :8000` (backend) dan `Ready` (frontend).

> Proyek ini hanya memakai **satu** file compose, yaitu `docker-compose.yml`.
> Tidak ada `docker-compose.dev.yml`, jadi tidak perlu opsi `-f` tambahan.

| Layanan              | URL (default)                                                  |
| -------------------- | -------------------------------------------------------------- |
| Web (Next.js)        | http://localhost:3100                                          |
| Units (contoh modul) | http://localhost:3100/master-data/units                        |
| API Fiber langsung   | http://localhost:8100                                          |
| Health check         | http://localhost:8100/health atau http://localhost:3100/health |
| phpMyAdmin           | http://localhost:9192                                          |

Port di atas adalah port **host**. Di dalam container, backend tetap mendengarkan
`:8000`. Nilai port diatur lewat `.env`; ubah di sana kalau bentrok dengan project lain.

### Data dummy

Ada dua cara mengisi 16 data units:

1. Set `SEED=true` di `.env`, jalankan `docker compose up -d`, lalu kembalikan ke `false`.
2. Atau impor `backend/seed_units.sql` lewat phpMyAdmin.

### Cek cepat bahwa API hidup

```powershell
Invoke-RestMethod http://localhost:8100/health
Invoke-RestMethod "http://localhost:8100/api/master-data/units/data?start=0&length=10"
```

Hasil `/health` yang sehat: `status: ok`, `database: connected`.

---

## Perintah Sehari-hari

```bash
docker compose ps
docker compose logs -f backend
docker compose exec backend bash
docker compose exec frontend sh
docker compose restart backend
docker compose up -d --build frontend    # bangun ulang satu layanan saja
```

Tambah dependency:

```bash
docker compose exec backend go get github.com/xxx/yyy
docker compose exec backend go mod tidy
docker compose exec frontend npm install nama-paket
```

Apa yang harus dilakukan setelah mengubah sesuatu:

| Perubahan                         | Yang dilakukan                             |
| --------------------------------- | ------------------------------------------ |
| File `.go` di `backend/`          | Otomatis, Air rebuild                      |
| File di `frontend/`               | Otomatis, Fast Refresh                     |
| `.env` atau `docker-compose.yml`  | `docker compose up -d`                     |
| `Dockerfile` atau `entrypoint.sh` | `docker compose up -d --build`             |
| `package.json` diedit manual      | `docker compose exec frontend npm install` |

Reset database (**menghapus semua data**): `docker compose down -v`

---

## Menambah Modul Baru

Contoh: modul `categories`. Pola paling mudah adalah menyalin file `units` lalu
mengganti namanya.

**Backend**

1. `models/categories.go` berisi struct, lalu daftarkan di `AutoMigrate` (`main.go`).
2. `services/categories_service.go`
3. `resources/categories_resource.go`
4. `controllers/categories_controller.go`
5. Daftarkan route di `routes/api.go`.
6. Tambah key permission `categories` di `config/auth.go`.

**Frontend**

1. `features/categories/api.ts` dan `CategoriesPage.tsx` (salin dari `features/units`).
2. `app/master-data/categories/page.tsx`
3. Tambah item di `MENU` pada `components/Navbar.tsx`.

---

## Contoh Modul Tanpa Database: Test Algoritma

Kalau hanya ingin mencoba fungsi murni (sorting, searching, graph) lewat HTTP, tidak
perlu model maupun database. Pola ini berguna sebagai latihan alur
`route -> controller -> service`.

```
backend/
├── services/algoritma.go                 fungsi murni: BubbleSort, QuickSort, ...
├── controllers/algorithm_controller.go   handler HTTP
└── routes/api.go                         daftarkan route
```

**Service** (`services/algoritma.go`): pakai `package services` dan nama fungsi
berawalan huruf besar supaya bisa dipanggil dari `controllers`.

**Controller** (`controllers/algorithm_controller.go`):

```go
package controllers

import (
	"github.com/gofiber/fiber/v2"

	"myApp/go-api/services"
)

type dataRequest struct {
	Data []int `json:"data"`
}

// POST /api/test/algorithm/bubble-sort   body: {"data":[5,3,8,1]}
func BubbleSort(c *fiber.Ctx) error {
	var req dataRequest
	if err := c.BodyParser(&req); err != nil {
		return validationError(c, "Body JSON tidak valid.")
	}
	return c.JSON(fiber.Map{
		"result": "OK",
		"input":  req.Data,
		"output": services.BubbleSort(req.Data),
	})
}
```

**Route** (`routes/api.go`, di dalam `Register`):

```go
algo := api.Group("/test/algorithm")
algo.Post("/bubble-sort", controllers.BubbleSort)
```

**Mencoba** (PowerShell, sesuaikan port host API):

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8100/api/test/algorithm/bubble-sort" `
  -ContentType "application/json" -Body '{"data":[64,34,25,12,22]}'
```

Respons yang diharapkan: `output` berisi `12, 22, 25, 34, 64`. Untuk data banyak,
simpan JSON ke file lalu kirim dengan `-InFile .\data.json`.

Catatan: Binary Search hanya benar pada data yang **sudah terurut**. Urutkan dulu
(misalnya dengan `services.QuickSort`) sebelum mencari.

---

## Aturan Penulisan Go di Proyek Ini

Sebagian besar error build di proyek ini berasal dari aturan dasar Go berikut.

| Aturan                                                                                                  | Akibat kalau dilanggar                             |
| ------------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| Semua file dalam satu folder memakai **nama package yang sama** (`services`, `controllers`, `routes`).  | `found packages service and services`              |
| Nama yang dipakai dari package lain harus **berawalan huruf besar** (`BubbleSort`, bukan `bubbleSort`). | `undefined` atau `cannot refer to unexported name` |
| Setiap import harus dipakai, dan setiap variabel lokal juga.                                            | `imported and not used`, `declared and not used`   |
| Satu nama fungsi hanya boleh dideklarasikan **sekali per package**, apa pun nama filenya.               | `redeclared in this block`                         |
| Hanya boleh ada satu `func main()` per package `main`.                                                  | `main redeclared in this block`                    |
| Teks memakai kutip ganda `"..."` atau backtick. Kutip tunggal `'x'` hanya untuk satu karakter.          | `empty rune literal or unescaped '`                |
| Di file `.go`, jangan menempel catatan atau teks biasa di luar fungsi.                                  | `non-declaration statement outside function body`  |
| Saat mengganti nama fungsi, ganti juga pemanggilan rekursifnya.                                         | `undefined: namaLama`                              |

---

## Troubleshooting

### Build gagal di Air (`failed to build`)

Lihat log backend:

```bash
docker compose logs --tail 40 backend
```

Cari baris `failed to build, error: exit status 1` lalu baca error compiler **di
atasnya**. Rujuk tabel [Aturan Penulisan Go](#aturan-penulisan-go-di-proyek-ini).
Baris `/app/tmp/main: not found` atau `Process Exit with Code: 127` yang muncul di
bawahnya hanya akibat build yang gagal, bukan masalah terpisah.

### `Connection reset` atau `Connection Failed` saat memanggil API

Port terbuka tetapi tidak ada program Go yang berjalan di dalam container. Hampir
selalu penyebabnya build Air yang gagal. Perbaiki error build (lihat bagian di atas),
tunggu sampai muncul `listening on :8000`, lalu coba lagi.

### `404 page not found` pada route baru

Pastikan route sudah didaftarkan di `routes/api.go` dan path-nya persis sama (tanda
hubung, bukan garis bawah). Setelah build berhasil, daftar route baru akan tampil di
log. Kalau tidak muncul, build terakhir kemungkinan gagal dan yang berjalan masih
binary lama.

### `The system cannot find the file specified` untuk file compose

Docker tidak menemukan file yang disebut lewat `-f`. Proyek ini hanya punya
`docker-compose.yml`, jadi jalankan tanpa `-f`:

```bash
docker compose up -d --build frontend
```

Pastikan juga terminal berada di folder yang berisi `docker-compose.yml`.

### Mengetik `.\file.go` di PowerShell tidak menjalankan apa-apa

Itu hanya membuka file dengan program bawaan Windows. Kode Go dijalankan lewat
container (otomatis oleh Air), atau `go run .` di dalam container:

```bash
docker compose exec backend go run .
```

### Port sudah dipakai

Ubah port host di `.env`, lalu jalankan `docker compose up -d`.

### Peringatan `punycode DeprecationWarning`

Berasal dari Node.js atau editor, tidak berbahaya, dan bisa diabaikan.

---

## Padanan Gin/Knockout ke Fiber/Next.js

| Versi lama                                     | Versi baru                                                |
| ---------------------------------------------- | --------------------------------------------------------- |
| `gin.Context`, `ctx.JSON`                      | `fiber.Ctx`, `c.JSON`                                     |
| `ShouldBindJSON` + `binding:"required,max=50"` | `c.BodyParser` + `validate()` manual per request          |
| `ctx.Request.FormValue("start")`               | `c.QueryInt("start", 0)`                                  |
| `r.Group("/api")`, `g.POST("")`                | `app.Group("/api")`, `g.Post("/")`                        |
| Middleware CORS buatan sendiri                 | `fiber/middleware/cors` (aktif jika `CORS_ORIGIN` diisi)  |
| `config/view.go` + `views/*.html`              | Dihapus; UI ada di `frontend/`                            |
| `window.AUTH_USER` (di-inject ke HTML)         | `GET /api/auth/me` -> `useApp().user`                     |
| `model.permission('update')`                   | `can("units", "update")`                                  |
| `model.Processing(true/false)`                 | `track(promise)` -> preloader otomatis                    |
| `ko.observable` / `ko.mapping`                 | `useState`                                                |
| jQuery DataTables server side                  | `<table>` + `<Pagination>` + `useEffect`                  |
| `Swal.fire` / `alert` / `confirm`              | `toast()` / `window.confirm`                              |
| Template + Air (`template/` -> `src/`)         | Dihapus; `backend/` dan `frontend/` langsung jadi project |

### Perubahan perilaku yang perlu diketahui

- **Respons daftar** berubah dari `{draw, RecordsTotal, RecordsFiltered, Data}`
  menjadi `{draw, total, filtered, data}`. Endpoint `POST /data` dihapus; sekarang
  hanya `GET /data?start=&length=&filtertext=`.
- **Validasi** mengembalikan `errors` per field, misalnya
  `{"code": "The code field is required."}`, dan ditampilkan di bawah input.
- **Tema Hope UI** (Bootstrap, `assets/`, `assetsadmin/`) tidak dibawa karena file
  vendor-nya tidak ada di zip. UI dibangun ulang dengan Tailwind dengan struktur yang
  sama (tab Form/List, pencarian, paginasi, tombol nonaktif sesuai permission).
- `container_name`, nama volume, dan nama network tidak lagi ditulis permanen; Compose
  menambahkan awalan nama folder. Yang perlu dibedakan antar project hanya port host
  di `.env`.

---

## Permission dan Keamanan

`backend/config/auth.go` berisi data contoh (semua akses diizinkan). Ganti dengan
user dan permission dari session, JWT, atau database. Permission di frontend hanya
mengatur tampilan; **pembatasan sebenarnya harus dipasang di server** (middleware
Fiber per route).

---

## Catatan

- Jangan commit `.env` berisi password sungguhan. Password di `.env.example` hanya
  untuk development lokal.
- Aplikasi contoh belum memiliki login; tambahkan autentikasi sebelum dipakai sungguhan.
- Image ini untuk **development** (toolchain Go + Air, `next dev`). Untuk production,
  gunakan build multi-stage (biner Go + `next build`). `BACKEND_URL` di
  `next.config.mjs` dibaca saat build, jadi set nilainya saat membangun image production.
- `backend/go.mod` sengaja tanpa daftar `require`; `go mod tidy` di `entrypoint.sh`
  mengisinya (dan membuat `go.sum`) saat start pertama.

---

## Credits

Berdasarkan **Gin MVC Docker Starter Kit** oleh **Arif Efendi** (© 2021-2026).
