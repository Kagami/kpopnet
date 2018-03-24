/**
 * Profile types/interfaces and accompanying functions.
 *
 * @module kpopnet/api/profiles
 */

export type ProfileValue = string | number | string[];

export interface Band {
  // Known props.
  id: string;
  name: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Idol {
  // Known props.
  id: string;
  name: string;
  band_id: string;
  image_id?: string;
  birth_name?: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Profiles {
  bands: Band[];
  idols: Idol[];
}

export interface BandInfo {
  band: Band;
  idols: Idol[];
}
export type BandMap = Map<string, BandInfo>;

export function getBandMap(profiles: Profiles): BandMap {
  const bandMap = new Map();
  profiles.bands.forEach((band) => {
    bandMap.set(band.id, { band, idols: [] });
  });
  profiles.idols.forEach((idol) => {
    // Backend guarantees every idol has associated band.
    // But not the other way around: band can have no members.
    bandMap.get(idol.band_id).idols.push(idol);
  });
  return bandMap;
}

export type RenderLine = [string, string];
export type Rendered = RenderLine[];

export function renderIdol(idol: Idol, band: Band): Rendered {
  const renderLine = renderLineCtx.bind(null, idol);
  const lines = getLines(idol).concat([["band", band.name]]);
  return lines.filter(keepLine).sort(compareLines).map(renderLine);
}

type InfoLine = [string, ProfileValue];

function getLines(idol: Idol): InfoLine[] {
  return Object.entries(idol);
}

const knownKeys = [
  "name",
  "birth_name",
  "band",
  "birth_date",
  "height",
  "weight",
  "positions",
];

const keyPriority = new Map(knownKeys
  // https://github.com/Microsoft/TypeScript/issues/6574
  .map((k, idx) => [k, idx] as [string, number]));

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
    return `${val} (${getAge(val as string)})`;
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

function renderLineCtx(idol: Idol, [key, val]: InfoLine): RenderLine {
  val = denormalizeVal(key, val, idol);
  key = denormalizeKey(key);
  return [key, val];
}

const MILLISECONDS_IN_YEAR = 1000 * 365 * 24 * 60 * 60;

export function getAge(birthday: string): number {
  const now = Date.now();
  // Birthday is always in YYYY-MM-DD form and can be parsed as
  // simplified ISO 8601 format.
  const born = new Date(birthday).getTime();
  return Math.floor((now - born) / MILLISECONDS_IN_YEAR);
}

// Remove symbols which doesn't make sense for fuzzy search.
function normalize(s: string): string {
  return s.replace(/[ .-]/g, "").toLowerCase();
}

// Minimal meaningful length.
const MIN_QUERY_LENGTH = 3;

function checkQuery(query: string): boolean {
  return query.length >= MIN_QUERY_LENGTH;
}

interface Query {
  name: string;
  props: Array<[string, string]>;
}

// Split query into main component and property-tagged parts.
// Example: "name words prop1:prop words prop2:more words"
function parseQuery(query: string): Query {
  let name = "";
  const props = [] as Array<[string, string]>;
  let lastKey = "";
  while (true) {
    // Search for prop1[:]
    const colonIdx = query.indexOf(":");
    if (colonIdx >= 1) {
      // Search for [ ]prop1:
      const spaceIdx = query.lastIndexOf(" ", colonIdx);
      if (spaceIdx >= 0) {
        // [name words] prop1:
        const lastVal = normalize(query.slice(0, spaceIdx));
        if (lastKey) {
          if (checkQuery(lastVal)) {
            props.push([lastKey, lastVal]);
          }
        } else {
          name = lastVal;
        }
        // [prop1]:...
        lastKey = query.slice(spaceIdx + 1, colonIdx);
        // prop1:[name words...]
        query = query.slice(colonIdx + 1);
      } else {
        // prop1:word:prop2
        if (lastKey) break;
        // Allow to start with []prop1:word
        lastKey = query.slice(0, colonIdx);
        // prop1:[name words...]
        query = query.slice(colonIdx + 1);
      }
    } else {
      // prop2:[more words]
      const lastVal = normalize(query);
      if (lastKey) {
        if (checkQuery(lastVal)) {
          props.push([lastKey, lastVal]);
        }
      } else {
        name = lastVal;
      }
      break;
    }
  }
  return {name, props};
}

// All allowed search props.
const propsMap = new Map([
  ["b", "band"],
]);

/**
 * Find idols matching given query.
 */
export function searchIdols(
  query: string, profiles: Profiles, bandMap: BandMap,
): Idol[] {
  if (!checkQuery(query)) return [];
  const q = parseQuery(query);
  if (!checkQuery(q.name) && !q.props.length) return [];
  return profiles.idols.filter((idol) => {
    const { band } = bandMap.get(idol.band_id);
    if (q.name) {
      if (normalize(idol.name).includes(q.name)) {
        return true;
      }
      if (idol.birth_name && normalize(idol.birth_name).includes(q.name)) {
        return true;
      }
      if (normalize(band.name).includes(q.name)) {
        return true;
      }
    }
    // tslint:disable-next-line:prefer-for-of
    for (let i = 0; i < q.props.length; i++) {
      const [keyId, val] = q.props[i];
      const key = propsMap.get(keyId);
      switch (key) {
      case "band":
        if (normalize(band.name).includes(val)) {
          return true;
        }
        break;
      }
    }
    return false;
  });
}
