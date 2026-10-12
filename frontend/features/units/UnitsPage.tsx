"use client";

/**
 * Halaman Units: padanan views/master_data/units.html + public/js/units.js
 * (Knockout + jQuery DataTables) dalam React.
 *
 *   units.Mode()              -> mode ("create" | "update")
 *   units.Recordunits         -> form (useState)
 *   units.FilterText          -> filterInput / query
 *   units.grid (DataTable)    -> rows + <Pagination> + useEffect(load)
 *   units.back / selectdata / save / remove -> back / selectData / save / remove
 */
import { useCallback, useEffect, useState, type FormEvent } from "react";

import Pagination from "@/components/Pagination";
import { useApp } from "@/context/AppContext";
import { ApiError } from "@/lib/api";
import type { Unit } from "@/lib/types";
import { emptyUnit, unitsApi, type UnitForm } from "./api";

const MODULE = "units";
const PAGE_SIZE = 10;

type Tab = "form" | "list";

export default function UnitsPage() {
  const { can, track, toast, ready } = useApp();

  const [tab, setTab] = useState<Tab>("list");
  const [form, setForm] = useState<UnitForm>(emptyUnit);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const isUpdate = form.id !== "";

  // ---- data list ----
  const [rows, setRows] = useState<Unit[]>([]);
  const [total, setTotal] = useState(0);
  const [filtered, setFiltered] = useState(0);
  const [page, setPage] = useState(1);
  const [filterInput, setFilterInput] = useState("");
  const [query, setQuery] = useState(""); // filter yang sedang diterapkan
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;

    track(unitsApi.list((page - 1) * PAGE_SIZE, PAGE_SIZE, query))
      .then((res) => {
        if (cancelled) return;
        setRows(res.data ?? []);
        setTotal(res.total);
        setFiltered(res.filtered);
      })
      .catch((e: ApiError) => !cancelled && toast("error", e.message));

    return () => {
      cancelled = true;
    };
  }, [page, query, reloadKey, track, toast]);

  const reload = useCallback(() => setReloadKey((k) => k + 1), []);

  // ---- actions ----
  const back = useCallback(() => {
    setForm(emptyUnit);
    setErrors({});
    setTab("list");
    reload();
  }, [reload]);

  const newData = () => {
    setForm(emptyUnit);
    setErrors({});
    setTab("form");
  };

  const selectData = async (reference: string) => {
    try {
      const unit = await track(unitsApi.select(reference));
      setForm({
        id: unit.id,
        reference: unit.reference,
        code: unit.code,
        name: unit.name,
        description: unit.description ?? "",
        is_active: !!unit.is_active,
      });
      setErrors({});
      setTab("form");
    } catch (e) {
      toast("error", e instanceof ApiError ? e.message : "Gagal mengambil data");
    }
  };

  const save = async (e: FormEvent) => {
    e.preventDefault();
    if (!(isUpdate ? can(MODULE, "update") : can(MODULE, "create"))) return;

    // validasi sisi client (server tetap memvalidasi ulang)
    const clientErrors: Record<string, string> = {};
    if (!form.code.trim()) clientErrors.code = "Code wajib diisi.";
    if (!form.name.trim()) clientErrors.name = "Name wajib diisi.";
    setErrors(clientErrors);
    if (Object.keys(clientErrors).length > 0) return;

    try {
      await track(unitsApi.save(form));
      toast("success", "Data berhasil disimpan.");
      back();
    } catch (err) {
      if (err instanceof ApiError) {
        setErrors(err.fieldErrors);
        if (Object.keys(err.fieldErrors).length === 0) toast("error", err.message);
      } else {
        toast("error", "Terjadi kesalahan.");
      }
    }
  };

  const remove = async (reference: string) => {
    if (!can(MODULE, "delete")) return;
    if (!confirm("Yakin hapus data?")) return;

    try {
      await track(unitsApi.remove(reference));
      toast("success", "Data berhasil dihapus.");
      // hapus data terakhir di halaman > 1 -> mundur satu halaman
      if (rows.length === 1 && page > 1) setPage(page - 1);
      back();
    } catch (e) {
      toast("error", e instanceof ApiError ? e.message : "Gagal menghapus data");
    }
  };

  const search = () => {
    setPage(1);
    setQuery(filterInput.trim());
  };

  // ---- UI helpers ----
  const canSave = isUpdate ? can(MODULE, "update") : can(MODULE, "create");
  const tabCls = (active: boolean, disabled = false) =>
    `px-4 py-2 text-sm font-medium border-b-2 ${
      active ? "border-indigo-600 text-indigo-600" : "border-transparent text-slate-500 hover:text-slate-800"
    } ${disabled ? "permission-disabled" : ""}`;
  const inputCls = (err?: string) =>
    `mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-300 ${
      err ? "border-rose-500" : "border-slate-300"
    }`;

  if (!ready) return null;

  return (
    <div className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
      {/* Tabs */}
      <div className="flex border-b border-slate-200">
        <button
          type="button"
          className={tabCls(tab === "form", !can(MODULE, "create"))}
          disabled={!can(MODULE, "create") && !isUpdate}
          onClick={() => (tab === "form" ? undefined : newData())}
        >
          Form
        </button>
        <button type="button" className={tabCls(tab === "list")} onClick={() => setTab("list")}>
          List
        </button>
      </div>

      {/* ============ FORM ============ */}
      {tab === "form" && (
        <form onSubmit={save} className="mt-4" noValidate>
          <div className="mb-4 flex flex-wrap items-center gap-2">
            {isUpdate && (
              <button
                type="button"
                onClick={back}
                className="rounded-md border border-slate-300 px-3 py-1.5 text-sm hover:bg-slate-50"
              >
                ← Kembali
              </button>
            )}

            <button
              type="submit"
              disabled={!canSave}
              className={`rounded-md bg-indigo-600 px-4 py-1.5 text-sm text-white hover:bg-indigo-700 ${
                canSave ? "" : "permission-disabled"
              }`}
            >
              Save
            </button>

            {isUpdate && (
              <button
                type="button"
                disabled={!can(MODULE, "delete")}
                onClick={() => remove(form.reference)}
                className={`rounded-md border border-rose-500 px-4 py-1.5 text-sm text-rose-600 hover:bg-rose-50 ${
                  can(MODULE, "delete") ? "" : "permission-disabled"
                }`}
              >
                Delete
              </button>
            )}

            {isUpdate && (
              <span className="rounded-md bg-slate-100 px-3 py-1.5 text-sm text-slate-600">
                ref : {form.reference}
              </span>
            )}
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <label className="text-sm font-medium">
                Code <span className="text-rose-600">*</span>
              </label>
              <input
                className={inputCls(errors.code)}
                value={form.code}
                maxLength={50}
                onChange={(e) => setForm({ ...form, code: e.target.value })}
              />
              {errors.code && <p className="mt-1 text-xs text-rose-600">{errors.code}</p>}
            </div>

            <div>
              <label className="text-sm font-medium">
                Name <span className="text-rose-600">*</span>
              </label>
              <input
                className={inputCls(errors.name)}
                value={form.name}
                maxLength={100}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
              />
              {errors.name && <p className="mt-1 text-xs text-rose-600">{errors.name}</p>}
            </div>

            <div>
              <label className="text-sm font-medium">Description</label>
              <textarea
                className={inputCls()}
                rows={3}
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
              />
            </div>

            <div>
              <label className="text-sm font-medium">Status</label>
              <label className="mt-2 flex cursor-pointer items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  className="h-4 w-4"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                />
                {form.is_active ? "Active" : "Inactive"}
              </label>
            </div>
          </div>
        </form>
      )}

      {/* ============ LIST ============ */}
      {tab === "list" && (
        <div className="mt-4">
          <div className="mb-3 flex gap-2">
            <input
              className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-300"
              placeholder="Searching.."
              value={filterInput}
              onChange={(e) => setFilterInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && search()}
            />
            <button
              type="button"
              onClick={search}
              className="rounded-md border border-indigo-600 px-4 text-sm text-indigo-600 hover:bg-indigo-50"
            >
              Search
            </button>
            {can(MODULE, "create") && (
              <button
                type="button"
                onClick={newData}
                className="whitespace-nowrap rounded-md bg-indigo-600 px-4 text-sm text-white hover:bg-indigo-700"
              >
                + New
              </button>
            )}
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-slate-200 bg-slate-50 text-slate-600">
                <tr>
                  <th className="px-3 py-2">No</th>
                  <th className="px-3 py-2">Code</th>
                  <th className="px-3 py-2">Name</th>
                  <th className="px-3 py-2">Description</th>
                  <th className="px-3 py-2">Status</th>
                  <th className="px-3 py-2">Action</th>
                </tr>
              </thead>
              <tbody>
                {rows.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-3 py-6 text-center text-slate-400">
                      No data available
                    </td>
                  </tr>
                )}

                {rows.map((row, i) => (
                  <tr key={row.id} className="border-b border-slate-100 odd:bg-white even:bg-slate-50/50">
                    <td className="px-3 py-2">{(page - 1) * PAGE_SIZE + i + 1}</td>
                    <td className="px-3 py-2">
                      <button
                        type="button"
                        disabled={!can(MODULE, "update")}
                        onClick={() => selectData(row.reference)}
                        className="text-indigo-600 hover:underline disabled:cursor-not-allowed disabled:text-slate-400 disabled:no-underline"
                      >
                        {row.code}
                      </button>
                    </td>
                    <td className="px-3 py-2">{row.name}</td>
                    <td className="px-3 py-2">{row.description}</td>
                    <td className="px-3 py-2">
                      <span
                        className={`rounded px-2 py-0.5 text-xs text-white ${
                          row.is_active ? "bg-emerald-600" : "bg-slate-400"
                        }`}
                      >
                        {row.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="whitespace-nowrap px-3 py-2">
                      <button
                        type="button"
                        disabled={!can(MODULE, "update")}
                        onClick={() => selectData(row.reference)}
                        className={`mr-3 text-slate-600 hover:text-indigo-600 ${
                          can(MODULE, "update") ? "" : "permission-disabled"
                        }`}
                      >
                        Edit
                      </button>
                      <button
                        type="button"
                        disabled={!can(MODULE, "delete")}
                        onClick={() => remove(row.reference)}
                        className={`text-slate-600 hover:text-rose-600 ${
                          can(MODULE, "delete") ? "" : "permission-disabled"
                        }`}
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Pagination page={page} pageSize={PAGE_SIZE} total={filtered} onChange={setPage} />
          {query && (
            <p className="mt-1 text-xs text-slate-400">
              (filtered from {total} total entries)
            </p>
          )}
        </div>
      )}
    </div>
  );
}
