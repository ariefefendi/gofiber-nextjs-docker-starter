import Link from "next/link";

const CARDS = [
  { href: "/master-data/units", title: "Units", desc: "Kelola master satuan barang.", external: false },
  { href: "/health", title: "Health", desc: "Cek status API dan database.", external: true },
];

export default function HomePage() {
  return (
    <div>
      <h1 className="mb-4 text-xl font-semibold">Dashboard</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {CARDS.map((c) => {
          const cls = "rounded-lg border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow";
          const body = (
            <>
              <h2 className="font-semibold">{c.title}</h2>
              <p className="text-sm text-slate-500">{c.desc}</p>
            </>
          );
          // /health diteruskan ke Fiber lewat rewrite -> pakai <a>, bukan <Link>
          return c.external ? (
            <a key={c.href} href={c.href} className={cls}>
              {body}
            </a>
          ) : (
            <Link key={c.href} href={c.href} className={cls}>
              {body}
            </Link>
          );
        })}
      </div>
    </div>
  );
}
