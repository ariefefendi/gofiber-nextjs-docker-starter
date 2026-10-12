"use client";

interface Props {
  page: number; // 1-based
  pageSize: number;
  total: number; // jumlah data setelah filter
  onChange: (page: number) => void;
}

/** Paginasi sederhana (padanan paging DataTables). */
export default function Pagination({ page, pageSize, total, onChange }: Props) {
  const pages = Math.max(1, Math.ceil(total / pageSize));
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);

  const btn =
    "rounded border border-slate-300 px-3 py-1 text-sm disabled:cursor-not-allowed disabled:opacity-40 hover:bg-slate-50";

  return (
    <div className="mt-3 flex flex-col items-center justify-between gap-2 text-sm text-slate-600 md:flex-row">
      <span>
        Showing {from} to {to} of {total} entries
      </span>
      <div className="flex items-center gap-2">
        <button className={btn} disabled={page <= 1} onClick={() => onChange(page - 1)}>
          Previous
        </button>
        <span>
          {page} / {pages}
        </span>
        <button className={btn} disabled={page >= pages} onClick={() => onChange(page + 1)}>
          Next
        </button>
      </div>
    </div>
  );
}
