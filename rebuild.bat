@echo off

echo ===============================
echo Fiber + Next.js Starter Kit
echo ===============================

REM Hanya membersihkan project INI (tidak menyentuh project Docker lain)
docker compose down --remove-orphans

docker compose build --no-cache

docker compose up -d

REM Baca port dari .env (default jika tidak ada)
set HOST_WEB_PORT=3100
set HOST_API_PORT=8100
set HOST_PMA_PORT=9192
for /f "usebackq tokens=1,* delims==" %%a in (`findstr /b "HOST_WEB_PORT= HOST_API_PORT= HOST_PMA_PORT=" .env`) do set %%a=%%b

echo.
echo ====================================
echo Web        : http://localhost:%HOST_WEB_PORT%
echo API Health : http://localhost:%HOST_API_PORT%/health
echo phpMyAdmin : http://localhost:%HOST_PMA_PORT%
echo ====================================

pause
