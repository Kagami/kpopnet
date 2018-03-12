import os
from os import path
from contextlib import suppress
import json
import hashlib
from datetime import datetime


INDEX_NAME = 'index'


# TODO(Kagami): Make it configurable.
def get_data_path():
    dpath = path.join(path.dirname(__file__), '..', '..', 'data')
    dpath = path.abspath(dpath)
    return dpath


def get_profiles_path():
    return path.join(get_data_path(), 'profiles')


def get_images_path():
    return path.join(get_data_path(), 'images')


def check_name(name):
    assert (name != '.'
            and name != '..'
            and '/' not in name
            and name != INDEX_NAME), 'Bad name'


def get_band_path(band):
    check_name(band['name'])
    fname = '{}.json'.format(INDEX_NAME)
    return path.join(get_profiles_path(), band['name'], fname)


def get_band_path_by_name(name):
    check_name(name)
    fname = '{}.json'.format(INDEX_NAME)
    return path.join(get_profiles_path(), name, fname)


def get_idol_path(band, idol):
    check_name(band['name'])
    check_name(idol['name'])
    fname = '{}.json'.format(idol['name'])
    return path.join(get_profiles_path(), band['name'], fname)


def strip_json_ext(fname):
    assert fname.endswith('.json'), 'Bad filename'
    return fname[:-5]


def load_json(b):
    # TODO(Kagami): Convert to native python types, e.g. birthday date.
    return json.loads(b, encoding='utf-8')


def default_encoder(o):
    if isinstance(o, datetime):
        return o.strftime('%Y-%m-%d')
    raise TypeError('Unknown type')


def dump_json(d):
    s = json.dumps(d, ensure_ascii=False, sort_keys=True, indent=4,
                   default=default_encoder)
    s += '\n'
    return s.encode('utf-8')


def update_profile(a, b):
    for k, v in b.items():
        if isinstance(v, list):
            try:
                s = set(a[k])
                s.update(v)
                a[k] = sorted(s)
            except (KeyError, TypeError):
                a[k] = sorted(v)
        elif k in ('id', 'band_id') and k in a:
            continue
        else:
            a[k] = v


_collected_band_urls = False
_all_band_urls = set()


def has_band_by_url(url):
    global _collected_band_urls
    if not _collected_band_urls:
        _collected_band_urls = True
        try:
            bnames = os.listdir(get_profiles_path())
        except OSError:
            return False
        else:
            for name in bnames:
                bpath = get_band_path_by_name(name)
                with suppress(OSError, KeyError):
                    band = load_json(open(bpath, 'rb').read())
                    _all_band_urls.update(band['urls'])
    return url in _all_band_urls


def save_band(updates):
    _all_band_urls.update(updates['urls'])
    bpath = get_band_path(updates)
    os.makedirs(path.dirname(bpath), exist_ok=True)
    try:
        band = load_json(open(bpath, 'rb').read())
    except OSError:
        band = {}
    update_profile(band, updates)
    with open(bpath, 'wb') as f:
        f.write(dump_json(band))


def save_idol(band, updates):
    ipath = get_idol_path(band, updates)
    os.makedirs(path.dirname(ipath), exist_ok=True)
    try:
        idol = load_json(open(ipath, 'rb').read())
    except OSError:
        idol = {}
    update_profile(idol, updates)
    with open(ipath, 'wb') as f:
        f.write(dump_json(idol))


def get_all_idol_names():
    try:
        bnames = os.listdir(get_profiles_path())
    except OSError:
        return []
    for bname in bnames:
        if bname == INDEX_NAME:
            continue
        try:
            inames = os.listdir(path.join(get_profiles_path(), bname))
        except OSError:
            continue
        for iname in inames:
            iname = strip_json_ext(iname)
            if iname == INDEX_NAME:
                continue
            yield bname, iname


_collected_idol_images = False
_all_idols_with_images = set()


def has_images_by_name(bname, iname):
    global _collected_idol_images
    if not _collected_idol_images:
        _collected_idol_images = True
        bnames = []
        with suppress(OSError):
            bnames = os.listdir(get_images_path())
        for bname in bnames:
            inames = []
            with suppress(OSError):
                inames = os.listdir(path.join(get_images_path(), bname))
            for iname in inames:
                _all_idols_with_images.add((bname, iname))
    return (bname, iname) in _all_idols_with_images


def save_image_by_name(bname, iname, data):
    _all_idols_with_images.add((bname, iname))
    md5 = hashlib.md5(data).hexdigest()
    # We only use JPEG files for simplicity.
    fname = '{}.jpg'.format(md5)
    fpath = path.join(get_images_path(), bname, iname, fname)
    os.makedirs(path.dirname(fpath), exist_ok=True)
    try:
        open(fpath, 'xb').write(data)
    except FileExistsError:
        return False
    else:
        return True
