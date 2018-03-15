/**
 * Reusable kpopnet components.
 *
 * Can be used to add profiles/face recognition functionality to
 * third-party sites (intended for cutechan).
 *
 * @module kpopnet/api
 */
// NOTE(Kagami): Make sure to import only essential modules here to keep
// build size small.

import { Profiles } from "./profiles";

export {
  ProfileValue, Band, Idol, Profiles,
  RenderLine, Rendered, renderIdol,
  BandMap, getBandMap,
} from "./profiles";

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
