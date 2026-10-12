@echo off
cd /d "%~dp0"

rem Build dulu: kalau gagal, container yang sedang jalan tidak disentuh
docker compose build || goto :fail

rem Hanya container yang berubah yang dibuat ulang; orphan dibersihkan
docker compose up -d --remove-orphans || goto :fail

rem Hapus image lama yang sudah tidak terpakai
docker image prune -f

echo Deploy selesai.
exit /b 0

:fail
echo Deploy gagal, lihat pesan error di atas.
exit /b 1
