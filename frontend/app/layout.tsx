import type { Metadata } from "next";
import type { ReactNode } from "react";

import Navbar from "@/components/Navbar";
import { AppProvider } from "@/context/AppContext";
import "./globals.css";

export const metadata: Metadata = {
  title: "Fiber + Next.js Starter",
  description: "Starter kit Go Fiber + Next.js + MySQL",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="id">
      <body className="min-h-screen bg-slate-50 text-slate-800 antialiased">
        <AppProvider>
          <Navbar />
          <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
        </AppProvider>
      </body>
    </html>
  );
}
