import os
from os import path
from contextlib import suppress
import json


# TODO(Kagami): Make it configurable.
def get_data_path():
    dpath = path.join(path.dirname(__file__), '..', 'data')
    dpath = path.abspath(dpath)
    return dpath


def get_profiles_path():
    return path.join(get_data_path(), 'profiles')


def check_name(name):
    assert (name != '.'
            and name != '..'
            and '/' not in name
            and name != 'index'), 'Bad name'


def get_band_path(band):
    check_name(band['name'])
    return path.join(get_profiles_path(), band['name'], 'index.json')


def get_band_path_by_name(name):
    check_name(name)
    return path.join(get_profiles_path(), name, 'index.json')


def get_member_path(band, member):
    check_name(band['name'])
    check_name(member['name'])
    fname = '{}.json'.format(member['name'])
    return path.join(get_profiles_path(), band['name'], fname)


def load_json(b):
    # TODO(Kagami): Convert to native python types, e.g. birthday date.
    return json.loads(b, encoding='utf-8')


def dump_json(d):
    s = json.dumps(d, ensure_ascii=False, sort_keys=True, indent=4,
                   default=str)
    s += '\n'
    return s.encode('utf-8')


def deep_update(a, b):
    for k, v in b.items():
        if isinstance(v, list):
            try:
                s = set(a[k])
                s.update(v)
                a[k] = sorted(s)
            except (KeyError, TypeError):
                a[k] = sorted(v)
        else:
            a[k] = v


_collected_bands_urls = False
_all_bands_urls = set()


def has_band_by_url(url):
    global _collected_bands_urls
    if not _collected_bands_urls:
        _collected_bands_urls = True
        try:
            names = os.listdir(get_profiles_path())
        except OSError:
            return False
        else:
            for name in names:
                bpath = get_band_path_by_name(name)
                with suppress(OSError, KeyError):
                    band = load_json(open(bpath, 'rb').read())
                    _all_bands_urls.update(band['urls'])
    return url in _all_bands_urls


def save_band(updates):
    _all_bands_urls.update(updates['urls'])
    bpath = get_band_path(updates)
    os.makedirs(path.dirname(bpath), exist_ok=True)
    try:
        band = load_json(open(bpath, 'rb').read())
    except OSError:
        band = {}
    deep_update(band, updates)
    with open(bpath, 'wb') as f:
        f.write(dump_json(band))


def save_member(band, updates):
    mpath = get_member_path(band, updates)
    os.makedirs(path.dirname(mpath), exist_ok=True)
    try:
        member = load_json(open(mpath, 'rb').read())
    except OSError:
        member = {}
    deep_update(member, updates)
    with open(mpath, 'wb') as f:
        f.write(dump_json(member))
