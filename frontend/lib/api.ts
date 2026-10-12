/**
 * Client tipis untuk API Go Fiber.
 * Base "/api" diteruskan ke backend oleh rewrite di next.config.mjs.
 */
const BASE = process.env.NEXT_PUBLIC_API_BASE ?? "/api";

export class ApiError extends Error {
  status: number;
  /** Pesan per field, mis. { code: "The code has already been taken." } */
  fieldErrors: Record<string, string>;

  constructor(status: number, message: string, fieldErrors: Record<string, string> = {}) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.fieldErrors = fieldErrors;
  }
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: body !== undefined ? { "Content-Type": "application/json" } : undefined,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let data: any = null;
  try {
    data = await res.json();
  } catch {
    /* respons kosong / bukan JSON */
  }

  if (!res.ok) {
    // Bentuk error Fiber: { result, message } atau { result, errors: string | {field: pesan} }
    const fieldErrors: Record<string, string> =
      data?.errors && typeof data.errors === "object" ? data.errors : {};

    const message: string =
      data?.message ||
      (typeof data?.errors === "string" ? data.errors : "") ||
      Object.values(fieldErrors).join(" ") ||
      "Terjadi kesalahan.";

    throw new ApiError(res.status, message, fieldErrors);
  }

  return data as T;
}

export const api = {
  get: <T>(path: string) => request<T>("GET", path),
  post: <T>(path: string, body: unknown) => request<T>("POST", path, body),
  put: <T>(path: string, body: unknown) => request<T>("PUT", path, body),
  delete: <T>(path: string) => request<T>("DELETE", path),
};
