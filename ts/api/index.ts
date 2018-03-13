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

/**
 * Get all profiles. ~47kb gzipped currently.
 */
export function getProfiles(): Promise<Profiles> {
  return fetch("/api/profiles", {credentials: "same-origin"})
    .then((res) => res.json());
}
