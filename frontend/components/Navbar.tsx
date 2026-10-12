"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const APP_NAME = process.env.NEXT_PUBLIC_APP_NAME ?? "Fiber + Next.js App";

const MENU = [
  { href: "/", label: "Dashboard" },
  { href: "/master-data/units", label: "Units" },
];

/**
 * Padanan placeholder navbar di views/layouts/main.html.
 * Ganti/tambah item MENU saat menambah modul baru.
 */
export default function Navbar() {
  const pathname = usePathname();

  return (
    <nav className="sticky top-0 z-40 border-b border-slate-200 bg-white">
      <div className="mx-auto flex max-w-6xl items-center gap-6 px-4 py-3">
        <Link href="/" className="font-bold text-slate-900">
          {APP_NAME}
        </Link>

        <div className="flex gap-4 text-sm">
          {MENU.map((m) => {
            const active = m.href === "/" ? pathname === "/" : pathname.startsWith(m.href);
            return (
              <Link
                key={m.href}
                href={m.href}
                className={active ? "font-semibold text-indigo-600" : "text-slate-600 hover:text-slate-900"}
              >
                {m.label}
              </Link>
            );
          })}
        </div>
      </div>
    </nav>
  );
}
