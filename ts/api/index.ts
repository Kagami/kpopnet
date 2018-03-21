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

import { Idol, Profiles } from "./profiles";

export {
  ProfileValue, Band, Idol, Profiles,
  RenderLine, Rendered, renderIdol,
  BandMap, getBandMap,
  searchIdols,
} from "./profiles";

export interface ApiOpts {
  prefix?: string;
}

/**
 * Get all profiles. ~47kb gzipped currently.
 */
export function getProfiles(opts: ApiOpts = {}): Promise<Profiles> {
  const prefix = opts.prefix || "/api";
  return fetch(`${prefix}/profiles`, {credentials: "same-origin"})
    .then((res) => res.json());
}

export interface FileOpts {
  small?: boolean;
  prefix?: string;
  fallback?: string;
}

/**
 * Get URL of the idol's preview image. Safe to use in <img> element
 * right away.
 */
export function getIdolPreviewUrl(idol: Idol, opts: FileOpts = {}): string {
  const prefix = opts.prefix || "/uploads";
  const fallback = opts.fallback || "/static/img/no-preview.svg";
  const sizeDir = opts.small ? "thumb" : "src";
  const sha1 = idol.image_id;
  // NOTE(Kagami): This assumes that filetype of the preview image is
  // always JPEG. This must be ensured by Idol API service.
  return sha1
    ? `${prefix}/${sizeDir}/${sha1.slice(0, 2)}/${sha1.slice(2)}.jpg`
    : fallback;
}
