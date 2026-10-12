"use client";

/**
 * Pengganti objek global `model` (model_head.js) di versi Knockout:
 *
 *   model.permission(module, action)  ->  can(module, action)
 *   model.Processing(true/false)      ->  track(promise) / processing
 *   window.AUTH_USER                  ->  user  (GET /api/auth/me)
 *   Swal.fire(...)                    ->  toast(type, message)
 */
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";

import { api } from "@/lib/api";
import type { AuthUser, PermissionAction } from "@/lib/types";

type ToastType = "success" | "error";
interface ToastItem {
  id: number;
  type: ToastType;
  message: string;
}

interface AppContextValue {
  user: AuthUser | null;
  /** true setelah /api/auth/me selesai (berhasil atau gagal) */
  ready: boolean;
  can: (module: string, action: PermissionAction) => boolean;
  /** true selama ada request yang dibungkus track() */
  processing: boolean;
  track: <T>(promise: Promise<T>) => Promise<T>;
  toast: (type: ToastType, message: string) => void;
}

const AppContext = createContext<AppContextValue | null>(null);

export function AppProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [ready, setReady] = useState(false);
  const [pending, setPending] = useState(0);
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const toastId = useRef(0);

  const toast = useCallback((type: ToastType, message: string) => {
    const id = ++toastId.current;
    setToasts((t) => [...t, { id, type, message }]);
    setTimeout(() => setToasts((t) => t.filter((x) => x.id !== id)), 4000);
  }, []);

  const track = useCallback(async <T,>(promise: Promise<T>): Promise<T> => {
    setPending((n) => n + 1);
    try {
      return await promise;
    } finally {
      setPending((n) => n - 1);
    }
  }, []);

  useEffect(() => {
    api
      .get<AuthUser>("/auth/me")
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setReady(true));
  }, []);

  const can = useCallback(
    (module: string, action: PermissionAction) =>
      user?.permissions?.[module.replace(/-/g, "_")]?.[action] === true,
    [user],
  );

  const value = useMemo<AppContextValue>(
    () => ({ user, ready, can, processing: pending > 0, track, toast }),
    [user, ready, can, pending, track, toast],
  );

  return (
    <AppContext.Provider value={value}>
      {children}

      {/* Preloader (padanan div.preloader + visible: model.Processing) */}
      {pending > 0 && (
        <div className="fixed inset-0 z-[9999] flex items-center justify-center bg-white/40">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-slate-300 border-t-indigo-600" />
        </div>
      )}

      {/* Toast */}
      <div className="fixed right-4 top-4 z-[10000] flex w-80 flex-col gap-2">
        {toasts.map((t) => (
          <div
            key={t.id}
            role="status"
            className={`rounded-lg px-4 py-3 text-sm text-white shadow-lg ${
              t.type === "success" ? "bg-emerald-600" : "bg-rose-600"
            }`}
          >
            {t.message}
          </div>
        ))}
      </div>
    </AppContext.Provider>
  );
}

export function useApp(): AppContextValue {
  const ctx = useContext(AppContext);
  if (!ctx) throw new Error("useApp harus dipakai di dalam <AppProvider>");
  return ctx;
}
