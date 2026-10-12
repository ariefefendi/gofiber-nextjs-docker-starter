/** @type {import('next').NextConfig} */

// Alamat API Go Fiber.
// Di Docker : http://backend:8000 (nama service compose)
// Di lokal  : http://localhost:8000
const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8000";

const nextConfig = {
  reactStrictMode: true,

  // Browser memanggil /api/... di origin Next.js, lalu Next meneruskannya
  // ke Fiber. Hasilnya: tidak perlu CORS dan URL backend tidak bocor ke client.
  async rewrites() {
    return [
      { source: "/api/:path*", destination: `${BACKEND_URL}/api/:path*` },
      { source: "/health", destination: `${BACKEND_URL}/health` },
    ];
  },
};

export default nextConfig;
