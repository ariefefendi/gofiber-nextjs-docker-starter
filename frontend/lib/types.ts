export type PermissionAction = "view" | "read" | "create" | "update" | "delete";

export type ModulePermissions = Partial<Record<PermissionAction, boolean>>;

/** Respons GET /api/auth/me (padanan window.AUTH_USER di versi Gin). */
export interface AuthUser {
  user_id: number;
  role_id: number;
  role_name: string;
  permissions: Record<string, ModulePermissions>;
}

/** Respons daftar dengan paginasi server side. */
export interface Paginated<T> {
  draw: number;
  total: number;
  filtered: number;
  data: T[];
}

export interface Unit {
  id: string;
  reference: string;
  code: string;
  name: string;
  description: string | null;
  is_active: boolean;
  created_at: string;
}
