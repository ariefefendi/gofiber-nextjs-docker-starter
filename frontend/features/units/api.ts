import { api } from "@/lib/api";
import type { Paginated, Unit } from "@/lib/types";

const BASE = "/master-data/units";

export interface UnitForm {
  id: string;
  reference: string;
  code: string;
  name: string;
  description: string;
  is_active: boolean;
}

export const emptyUnit: UnitForm = {
  id: "",
  reference: "",
  code: "",
  name: "",
  description: "",
  is_active: true,
};

export const unitsApi = {
  list: (start: number, length: number, filtertext: string) =>
    api.get<Paginated<Unit>>(
      `${BASE}/data?start=${start}&length=${length}&filtertext=${encodeURIComponent(filtertext)}`,
    ),

  select: (reference: string) =>
    api.get<Unit>(`${BASE}/select?reference=${encodeURIComponent(reference)}`),

  save: (form: UnitForm) => {
    const body = {
      id: form.id,
      code: form.code,
      name: form.name,
      description: form.description === "" ? null : form.description,
      is_active: form.is_active,
    };
    // Insert = POST, Update = PUT (id terisi)
    return form.id ? api.put(BASE, body) : api.post(BASE, body);
  },

  remove: (reference: string) =>
    api.delete(`${BASE}?reference=${encodeURIComponent(reference)}`),
};
