# Go Fiber + Next.js Docker Starter Kit

> Full-stack dev environment in one command: `docker compose up -d --build`

Go Fiber API · GORM · MySQL · Next.js (App Router, TypeScript, Tailwind) · Air hot reload · phpMyAdmin

A development starter kit built on **Go Fiber + GORM + MySQL** (API) and
**Next.js** (App Router, TypeScript, Tailwind) as the frontend. It includes hot
reload (Air for Go, `next dev` for Next.js) and phpMyAdmin out of the box.

Converted from the _Gin MVC Docker Starter Kit_ (Gin + HTML templates + Knockout.js).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Everyday Commands](#everyday-commands)
- [Adding a New Module](#adding-a-new-module)
- [Example Module Without a Database: Algorithm Tests](#example-module-without-a-database-algorithm-tests)
- [Go Writing Rules for This Project](#go-writing-rules-for-this-project)
- [Troubleshooting](#troubleshooting)
- [Gin/Knockout to Fiber/Next.js Mapping](#ginknockout-to-fibernextjs-mapping)
- [Permissions and Security](#permissions-and-security)
- [Notes](#notes)
- [Credits](#credits)

---

## Prerequisites

- Docker Desktop (or Docker Engine) with **Docker Compose v2** (the `docker compose` command, not `docker-compose`).
- Free host ports for the web app, the API, and phpMyAdmin (defaults are listed in the [Getting Started](#getting-started) table).
- You do not need Go or Node.js installed on your machine; everything runs inside containers.

---

## Project Structure

```
.
├── docker-compose.yml       <- the only compose file
├── .env.example             <- copy to .env
├── rebuild.bat
├── backend/                 Go Fiber API
│   ├── Dockerfile, entrypoint.sh, .air.toml, go.mod
│   ├── main.go
│   ├── config/              env.go, database.go, auth.go
│   ├── routes/              api.go           (URL -> controller, no logic)
│   ├── controllers/         units_controller.go, auth_controller.go, helpers.go
│   ├── services/            business logic: units_service.go, id.go, seed_units.go
│   ├── models/              GORM table structs
│   ├── resources/           JSON response shapes
│   ├── middleware/          cors.go
│   └── seed_units.sql       dummy data (optional, via phpMyAdmin)
└── frontend/                Next.js
    ├── Dockerfile, package.json, next.config.mjs
    ├── app/                 layout.tsx, page.tsx, master-data/units/page.tsx
    ├── components/          Navbar, Pagination
    ├── context/             AppContext.tsx (auth, permission, loading, toast)
    ├── features/units/      UnitsPage.tsx, api.ts
    └── lib/                 api.ts (fetch wrapper), types.ts
```

Request flow:

```
Browser -> Next.js (:3100) --/api/* rewrite--> Fiber (:8000 inside the container)
                                               Route -> Controller -> Service -> Model (GORM) -> MySQL
```

The browser only talks to Next.js. Requests to `/api/*` and `/health` are forwarded
to Fiber through `rewrites()` in `frontend/next.config.mjs`, so **CORS is not needed**.

Responsibility of each layer:

| Layer      | Folder         | Responsibility                                               |
| ---------- | -------------- | ------------------------------------------------------------ |
| Route      | `routes/`      | Maps URLs to controllers. Contains no logic.                 |
| Controller | `controllers/` | Reads the request, validates, calls the service, sends JSON. |
| Service    | `services/`    | Business logic and pure functions. Knows nothing about HTTP. |
| Model      | `models/`      | GORM table structs.                                          |

---

## Getting Started

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

The first build takes longer because it downloads Go dependencies and runs `npm install`.
Wait until you see `listening on :8000` (backend) and `Ready` (frontend).

> This project uses only **one** compose file: `docker-compose.yml`.
> There is no `docker-compose.dev.yml`, so no extra `-f` option is needed.

| Service                | URL (default)                                                |
| ---------------------- | ------------------------------------------------------------ |
| Web (Next.js)          | http://localhost:3100                                        |
| Units (example module) | http://localhost:3100/master-data/units                      |
| Fiber API (direct)     | http://localhost:8100                                        |
| Health check           | http://localhost:8100/health or http://localhost:3100/health |
| phpMyAdmin             | http://localhost:9192                                        |

The ports above are **host** ports. Inside the container the backend still listens on
`:8000`. Port values are set in `.env`; change them there if they clash with another project.

### Dummy data

There are two ways to load the 16 sample units:

1. Set `SEED=true` in `.env`, run `docker compose up -d`, then set it back to `false`.
2. Or import `backend/seed_units.sql` through phpMyAdmin.

### Quick check that the API is up

```powershell
Invoke-RestMethod http://localhost:8100/health
Invoke-RestMethod "http://localhost:8100/api/master-data/units/data?start=0&length=10"
```

A healthy `/health` response is `status: ok`, `database: connected`.

---

## Everyday Commands

```bash
docker compose ps
docker compose logs -f backend
docker compose exec backend bash
docker compose exec frontend sh
docker compose restart backend
docker compose up -d --build frontend    # rebuild a single service
```

Add a dependency:

```bash
docker compose exec backend go get github.com/xxx/yyy
docker compose exec backend go mod tidy
docker compose exec frontend npm install package-name
```

What to do after changing something:

| Change                          | What to do                                 |
| ------------------------------- | ------------------------------------------ |
| `.go` files in `backend/`       | Automatic, Air rebuilds                    |
| Files in `frontend/`            | Automatic, Fast Refresh                    |
| `.env` or `docker-compose.yml`  | `docker compose up -d`                     |
| `Dockerfile` or `entrypoint.sh` | `docker compose up -d --build`             |
| `package.json` edited manually  | `docker compose exec frontend npm install` |

Reset the database (**deletes all data**): `docker compose down -v`

---

## Adding a New Module

Example: a `categories` module. The easiest approach is to copy the `units` files and
rename them.

**Backend**

1. `models/categories.go` with the struct, then register it in `AutoMigrate` (`main.go`).
2. `services/categories_service.go`
3. `resources/categories_resource.go`
4. `controllers/categories_controller.go`
5. Register the routes in `routes/api.go`.
6. Add a `categories` permission key in `config/auth.go`.

**Frontend**

1. `features/categories/api.ts` and `CategoriesPage.tsx` (copy from `features/units`).
2. `app/master-data/categories/page.tsx`
3. Add an item to `MENU` in `components/Navbar.tsx`.

---

## Example Module Without a Database: Algorithm Tests

If you only want to try pure functions (sorting, searching, graphs) over HTTP, you do
not need a model or a database. This pattern is also a good exercise for the
`route -> controller -> service` flow.

```
backend/
├── services/algoritma.go                 pure functions: BubbleSort, QuickSort, ...
├── controllers/algorithm_controller.go   HTTP handlers
└── routes/api.go                         register the routes
```

**Service** (`services/algoritma.go`): use `package services` and start function names
with an uppercase letter so they can be called from `controllers`.

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
		return validationError(c, "Invalid JSON body.")
	}
	return c.JSON(fiber.Map{
		"result": "OK",
		"input":  req.Data,
		"output": services.BubbleSort(req.Data),
	})
}
```

**Route** (`routes/api.go`, inside `Register`):

```go
algo := api.Group("/test/algorithm")
algo.Post("/bubble-sort", controllers.BubbleSort)
```

**Trying it** (PowerShell, adjust the API host port):

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8100/api/test/algorithm/bubble-sort" `
  -ContentType "application/json" -Body '{"data":[64,34,25,12,22]}'
```

Expected response: `output` contains `12, 22, 25, 34, 64`. For larger inputs,
save the JSON to a file and send it with `-InFile .\data.json`.

Note: Binary Search is only correct on data that is **already sorted**. Sort first
(for example with `services.QuickSort`) before searching.

---

## Go Writing Rules for This Project

Most build errors in this project come from the basic Go rules below.

| Rule                                                                                                      | What happens if you break it                      |
| --------------------------------------------------------------------------------------------------------- | ------------------------------------------------- |
| All files in one folder use the **same package name** (`services`, `controllers`, `routes`).              | `found packages service and services`             |
| Names used from another package must **start with an uppercase letter** (`BubbleSort`, not `bubbleSort`). | `undefined` or `cannot refer to unexported name`  |
| Every import must be used, and so must every local variable.                                              | `imported and not used`, `declared and not used`  |
| A function name may be declared **only once per package**, regardless of the file name.                   | `redeclared in this block`                        |
| A `package main` may contain only one `func main()`.                                                      | `main redeclared in this block`                   |
| Strings use double quotes `"..."` or backticks. Single quotes `'x'` are for one character only.           | `empty rune literal or unescaped '`               |
| Do not paste notes or plain text outside functions in a `.go` file.                                       | `non-declaration statement outside function body` |
| When you rename a function, rename its recursive calls too.                                               | `undefined: oldName`                              |

---

## Troubleshooting

### Build fails in Air (`failed to build`)

Check the backend log:

```bash
docker compose logs --tail 40 backend
```

Look for the line `failed to build, error: exit status 1` and read the compiler
error **above** it. Refer to the [Go Writing Rules](#go-writing-rules-for-this-project)
table. The `/app/tmp/main: not found` or `Process Exit with Code: 127` lines below it are
only a consequence of the failed build, not a separate problem.

### `Connection reset` or `Connection Failed` when calling the API

The port is open but no Go program is running inside the container. This is almost
always caused by a failed Air build. Fix the build error (see the section above), wait
until `listening on :8000` appears, then try again.

### `404 page not found` on a new route

Make sure the route is registered in `routes/api.go` and the path matches exactly
(hyphens, not underscores). After a successful build, the new route is listed in the
log. If it is missing, the last build probably failed and the old binary is still running.

### `The system cannot find the file specified` for a compose file

Docker cannot find a file named with `-f`. This project only has `docker-compose.yml`,
so run without `-f`:

```bash
docker compose up -d --build frontend
```

Also make sure your terminal is in the folder that contains `docker-compose.yml`.

### Typing `.\file.go` in PowerShell does nothing

That only opens the file with the default Windows program. Go code runs through the
container (automatically via Air), or with `go run .` inside the container:

```bash
docker compose exec backend go run .
```

### Port already in use

Change the host port in `.env`, then run `docker compose up -d`.

### `punycode DeprecationWarning`

It comes from Node.js or the editor, is harmless, and can be ignored.

---

## Gin/Knockout to Fiber/Next.js Mapping

| Old version                                    | New version                                                   |
| ---------------------------------------------- | ------------------------------------------------------------- |
| `gin.Context`, `ctx.JSON`                      | `fiber.Ctx`, `c.JSON`                                         |
| `ShouldBindJSON` + `binding:"required,max=50"` | `c.BodyParser` + manual `validate()` per request              |
| `ctx.Request.FormValue("start")`               | `c.QueryInt("start", 0)`                                      |
| `r.Group("/api")`, `g.POST("")`                | `app.Group("/api")`, `g.Post("/")`                            |
| Custom CORS middleware                         | `fiber/middleware/cors` (enabled when `CORS_ORIGIN` is set)   |
| `config/view.go` + `views/*.html`              | Removed; the UI lives in `frontend/`                          |
| `window.AUTH_USER` (injected into HTML)        | `GET /api/auth/me` -> `useApp().user`                         |
| `model.permission('update')`                   | `can("units", "update")`                                      |
| `model.Processing(true/false)`                 | `track(promise)` -> automatic preloader                       |
| `ko.observable` / `ko.mapping`                 | `useState`                                                    |
| jQuery DataTables server side                  | `<table>` + `<Pagination>` + `useEffect`                      |
| `Swal.fire` / `alert` / `confirm`              | `toast()` / `window.confirm`                                  |
| Template + Air (`template/` -> `src/`)         | Removed; `backend/` and `frontend/` are the projects directly |

### Behavior changes to be aware of

- **List responses** changed from `{draw, RecordsTotal, RecordsFiltered, Data}`
  to `{draw, total, filtered, data}`. The `POST /data` endpoint was removed; only
  `GET /data?start=&length=&filtertext=` remains.
- **Validation** returns `errors` per field, for example
  `{"code": "The code field is required."}`, which are shown under each input.
- **The Hope UI theme** (Bootstrap, `assets/`, `assetsadmin/`) was not carried over
  because its vendor files were not in the zip. The UI was rebuilt with Tailwind using
  the same structure (Form/List tabs, search, pagination, buttons disabled by permission).
- `container_name`, volume names, and network names are no longer hardcoded; Compose
  prefixes them with the folder name. Only the host ports in `.env` need to differ
  between projects.

---

## Permissions and Security

`backend/config/auth.go` contains sample data (all access allowed). Replace it with
users and permissions from a session, JWT, or database. Frontend permissions only
control what is displayed; **the real restriction must be enforced on the server**
(Fiber middleware per route).

---

## Notes

- Never commit a `.env` file with real passwords. The password in `.env.example` is
  for local development only.
- The sample app has no login yet; add authentication before using it for real.
- This image is for **development** (Go toolchain + Air, `next dev`). For production,
  use a multi-stage build (Go binary + `next build`). `BACKEND_URL` in
  `next.config.mjs` is read at build time, so set it when building the production image.
- `backend/go.mod` intentionally has no `require` list; `go mod tidy` in
  `entrypoint.sh` fills it in (and creates `go.sum`) on first start.

---

## Credits

Based on the _*Gin MVC Docker Starter Kit*_ by **Arif Efendi** (© 2021-2026).
