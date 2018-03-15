/**
 * Reusable kpopnet components.
 *
 * May be used to add profiles/face recognition functionality to
 * third-party sites (intended for cutechan).
 *
 * @module kpopnet/api
 */

// NOTE(Kagami): Make sure to import only essential modules here to keep
// build size small.

// Profile structs, reflect JSON structure.

export type ProfileValue = string | number | string[];

export interface Band {
  // Mandatory props.
  id: string;
  name: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Idol {
  // Mandatory props.
  id: string;
  band_id: string;
  name: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Profiles {
  bands: Band[];
  idols: Idol[];
}

// Convert profile info to human-readable form, partial revert of
// python's spiders work.

type InfoLine = [string, ProfileValue];

function getLines(idol: Idol): InfoLine[] {
  return Object.entries(idol);
}

const knownKeys = [
  "name",
  "birth_name",
  "birth_date",
  "height",
  "weight",
  "positions",
];
const keyPriority = new Map(knownKeys
  // https://github.com/Microsoft/TypeScript/issues/6574
  .map((k, i) => [k, i] as [string, number]));

function keepLine([key, val]: InfoLine): boolean {
  return keyPriority.has(key);
}

function compareLines(a: InfoLine, b: InfoLine): number {
  const k1 = a[0];
  const k2 = b[0];
  return keyPriority.get(k1) - keyPriority.get(k2);
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function denormalizeKey(key: string): string {
  switch (key) {
  case "birth_date":
    return "Birthday";
  case "birth_name":
    return "Real name";
  }
  key = capitalize(key);
  key = key.replace(/_/g, " ");
  return key;
}

function denormalizeVal(key: string, val: ProfileValue, idol: Idol): string {
  switch (key) {
  case "name":
    const hangul = idol.name_hangul;
    return hangul ? `${val} (${hangul})` : val as string;
  case "birth_name":
    const hangul2 = idol.birth_name_hangul;
    return hangul2 ? `${val} (${hangul2})` : val as string;
  case "birth_date":
    const zodiac = idol.zodiac;
    return zodiac ? `${val} (${zodiac})` : val as string;
  case "height":
    return val + " cm";
  case "weight":
    return val + " kg";
  case "positions":
    return (val as string[]).join(", ");
  default:
    return val.toString();
  }
}

function showLineCtx(idol: Idol, [key, val]: InfoLine): [string, string] {
  val = denormalizeVal(key, val, idol);
  key = denormalizeKey(key);
  return [key, val];
}

export function showIdol(idol: Idol): Array<[string, string]> {
  const showLine = showLineCtx.bind(null, idol);
  return getLines(idol).filter(keepLine).sort(compareLines).map(showLine);
}

// Communicate with backend.

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
