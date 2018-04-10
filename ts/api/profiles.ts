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
  alt_names?: string[];
  label_icon?: string;
  label_name?: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Idol {
  // Known props.
  id: string;
  name: string;
  name_hangul?: string;
  alt_names?: string[];
  band_id: string;
  alt_band_ids?: string[];
  image_id?: string;
  birth_name?: string;
  birth_name_hangul?: string;
  korean_name?: string;
  // Other props.
  [key: string]: ProfileValue;
}

export interface Profiles {
  bands: Band[];
  idols: Idol[];
}

export type BandMap = Map<string, Band>;

export function getBandMap(profiles: Profiles): BandMap {
  const bandMap = new Map();
  profiles.bands.forEach((band) => {
    bandMap.set(band.id, band);
  });
  return bandMap;
}

export type IdolMap = Map<string, Idol>;

export function getIdolMap(profiles: Profiles): IdolMap {
  const idolMap = new Map();
  profiles.idols.forEach((idol) => {
    idolMap.set(idol.id, idol);
  });
  return idolMap;
}

function tryConcat<T>(a: T, b?: T[]): T[] {
  return b ? [a].concat(b) : [a];
}

// TODO(Kagami): Profile/optimize.
function getBandNames(idol: Idol, bandMap: BandMap): string[] {
  // NOTE(Kagami): Backend doesn't currently guarantee that alt_band_ids
  // contain correct existing band ids, so this may potentially fail.
  const bandIds = tryConcat(idol.band_id, idol.alt_band_ids);
  const bands = bandIds.map((id) => bandMap.get(id));
  const bnames = [] as string[];
  bands.forEach((b) => {
    tryConcat(b.name, b.alt_names).forEach((bname) => {
      bnames.push(bname);
    });
  });
  return bnames;
}

export type RenderedLine = [string, string];
export type Rendered = RenderedLine[];

export function renderIdol(idol: Idol, bandMap: BandMap): Rendered {
  const renderLine = renderLineCtx.bind(null, idol);
  const bnames = getBandNames(idol, bandMap);
  // Main name of the main band goes first.
  let bname = bnames[0];
  if (bnames.length > 1) {
    bname += " (";
    bname += bnames.slice(1).join(", ");
    bname += ")";
  }
  const lines = getLines(idol).concat([["bands", bname]]);
  return lines.filter(keepLine).sort(compareLines).map(renderLine);
}

type InfoLine = [string, ProfileValue];

function getLines(idol: Idol): InfoLine[] {
  return Object.entries(idol);
}

const knownKeys = [
  "name",
  "birth_name",
  "bands",
  "birth_date",
  "height",
  "weight",
  "positions",
];

const keyPriority = new Map(knownKeys
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

function renderLineCtx(idol: Idol, [key, val]: InfoLine): RenderedLine {
  val = denormalizeVal(key, val, idol);
  key = denormalizeKey(key);
  return [key, val];
}

const MILLISECONDS_IN_YEAR = 1000 * 365 * 24 * 60 * 60;

function getAge(birthday: string): number {
  const now = Date.now();
  // Birthday is always in YYYY-MM-DD form and can be parsed as
  // simplified ISO 8601 format.
  const born = new Date(birthday).getTime();
  return Math.floor((now - born) / MILLISECONDS_IN_YEAR);
}

// Remove symbols which doesn't make sense for fuzzy search.
function normalize(s: string): string {
  return s.replace(/[' .-]/g, "").toLowerCase();
}

type SearchProp = [string, string];

interface Query {
  name: string;
  props: SearchProp[];
}

// Split query into main component and property-tagged parts.
// Example: "name words prop1:prop words prop2:more words"
// TODO(Kagami): Profile/optimize.
function parseQuery(query: string): Query {
  let name = "";
  const props = [] as SearchProp[];
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
          props.push([lastKey, lastVal]);
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
        props.push([lastKey, lastVal]);
      } else {
        name = lastVal;
      }
      break;
    }
  }
  return {name, props};
}

// Match all possible names of the idol.
// TODO(Kagami): Search for hangul?
function matchIdolName(idol: Idol, val: string): boolean {
  if (normalize(idol.name).includes(val)) {
    return true;
  }
  if (idol.birth_name && normalize(idol.birth_name).includes(val)) {
    return true;
  }
  if (idol.korean_name && normalize(idol.korean_name).includes(val)) {
    return true;
  }
  if (idol.alt_names && idol.alt_names.some((n) => normalize(n).includes(val))) {
    return true;
  }
  return false;
}

// Match all possible band names.
function matchBandName(idol: Idol, bandMap: BandMap, val: string): boolean {
  const bnames = getBandNames(idol, bandMap);
  return bnames.some((bname) => normalize(bname).includes(val));
}

/**
 * Find idols matching given query.
 */
// TODO(Kagami): Profile/optimize.
export function searchIdols(
  query: string, profiles: Profiles, bandMap: BandMap,
): Idol[] {
  if (query.length < 3) return [];
  // console.time("parseQuery");
  const q = parseQuery(query);
  // console.timeEnd("parseQuery");
  if (!q.name && !q.props.length) return [];

  // TODO(Kagam): Sort idols?
  // console.time("searchIdols");
  const result = profiles.idols.filter((idol) => {
    // Fuzzy name matching.
    // TODO(Kagami): Allow combinations like "Orange Caramel lizzy"
    if (q.name
        && !matchIdolName(idol, q.name)
        && !matchBandName(idol, bandMap, q.name)) {
      return false;
    }
    // Match for exact properties if user requested.
    return q.props.every(([key, val]) => {
      switch (key) {
      case "n":
      case "name":
        if (normalize(idol.name).includes(val)) {
          return true;
        }
        break;
      case "rn":
        if (normalize(idol.birth_name || "").includes(val)) {
          return true;
        }
        break;
      case "b":
      case "band":
        if (normalize(bandMap.get(idol.band_id).name).includes(val)) {
          return true;
        }
        break;
      }
      return false;
    });
  });
  // console.timeEnd("searchIdols");
  return result;
}
