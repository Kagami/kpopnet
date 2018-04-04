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
  RenderedLine, Rendered, renderIdol,
  BandMap, IdolMap, getBandMap, getIdolMap,
  searchIdols,
} from "./profiles";

function handleResponse(res: Response): Promise<any> {
  return res.ok ? res.json() : handleErrorCode(res);
}

function handleErrorCode(res: Response): Promise<any> {
  const unknown = "unknown error";
  const ctype = res.headers.get("Content-Type");
  const isHtml = ctype.startsWith("text/html");
  const isJson = ctype.startsWith("application/json");
  if (isHtml) {
    // Probably 404/500 page, not bother parsing.
    throw new Error(unknown);
  } else if (isJson) {
    // Probably standardly-shaped JSON error.
    return res.json().then((data) => {
      throw new Error(data && data.error || unknown);
    });
  } else {
    // Probably text/plain or something like this.
    return res.text().then((data) => {
      throw new Error(data || unknown);
    });
  }
}

function handleError(err: Error) {
  throw new Error(err.message || "unknown error");
}

export interface ApiOpts {
  prefix?: string;
}

/**
 * Get all profiles. ~47kb gzipped currently.
 */
export function getProfiles(opts: ApiOpts = {}): Promise<Profiles> {
  const prefix = opts.prefix || "/api";
  return fetch(`${prefix}/idols/profiles`, {credentials: "same-origin"})
    .then(handleResponse, handleError);
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
  // always JPEG. It must be ensured by Idol API service.
  return sha1
    ? `${prefix}/${sizeDir}/${sha1.slice(0, 2)}/${sha1.slice(2)}.jpg`
    : fallback;
}

export interface ImageIdData {
  SHA1: string;
}

/**
 * Set idol's preview.
 */
export function setIdolPreview(idol: Idol, file: File, opts: ApiOpts = {}): Promise<ImageIdData> {
  const prefix = opts.prefix || "/api";
  const form = new FormData();
  form.append("files[]", file);
  return fetch(`${prefix}/idols/${idol.id}/preview`, {
    credentials: "same-origin",
    method: "POST",
    body: form,
  }).then(handleResponse, handleError);
}

export interface IdolIdData {
  id: string;
}

/**
 * Recognize idol.
 */
export function recognizeIdol(file: File, opts: ApiOpts = {}): Promise<IdolIdData> {
  const prefix = opts.prefix || "/api";
  const form = new FormData();
  form.append("files[]", file);
  return fetch(`${prefix}/idols/recognize`, {
    credentials: "same-origin",
    method: "POST",
    body: form,
  }).then(handleResponse, handleError);
}
