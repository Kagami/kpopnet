/**
 * Profile types/interfaces and accompaying functions.
 *
 * @module kpopnet/api/profiles
 */

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
  for (const { idols } of bandMap.values()) {
    idols.sort(compareIdols);
  }
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
  case "zodiac":
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

function renderLineCtx(idol: Idol, [key, val]: InfoLine): RenderLine {
  val = denormalizeVal(key, val, idol);
  key = denormalizeKey(key);
  return [key, val];
}

// Convert YYYY-MM-DD to JS Date object.
// Backend guarantees to return birtday only in this format.
function birtday2date(birthday: string): Date {
  const [y, m, d] = birthday.split(/-/, 3);
  return new Date(+y, +m - 1, +d);
}

const MILLISECONDS_IN_YEAR = 1000 * 365 * 24 * 60 * 60;

export function getAge(birthday: string): number {
  const now = Date.now();
  const born = birtday2date(birthday).getTime();
  return Math.floor((now - born) / MILLISECONDS_IN_YEAR);
}

function compareIdols(a: Idol, b: Idol): number {
  const ageA = a.birth_date ? getAge(a.birth_date as string) : 0;
  const ageB = b.birth_date ? getAge(b.birth_date as string) : 0;
  return ageB - ageA;
}
