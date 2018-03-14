/**
 * Reusable kpopnet components.
 *
 * May be used to add profiles/face recognition functionality to
 * third-party sites (intended for cutechan).
 *
 * @module kpopnet/api
 */

// Profile structs, reflect go's kpopnet/profile/.

export interface Band {
  // Mandatory props.
  id: string;
  name: string;
  // Other props.
  [key: string]: any;
}

export interface Idol {
  // Mandatory props.
  id: string;
  band_id: string;
  name: string;
  // Other props.
  [key: string]: any;
}

export interface Profiles {
  bands: Band[];
  idols: Idol[];
}

// Defined in webpack's config.
declare const API_PREFIX: string;

function get(resource: string): Promise<Response> {
  return fetch(`${API_PREFIX}/api/${resource}`, {credentials: "same-origin"});
}

/**
 * Get all profiles. ~47kb gzipped currently.
 */
export function getProfiles(): Promise<Profiles> {
  return get("profiles").then((res) => res.json());
}

/**
 * Get URL of the idol's preview image. Safe to use in <img> element
 * right away.
 */
export function getIdolPreviewUrl(id: string): string {
  return `${API_PREFIX}/api/idols/${id}/preview`;
}
