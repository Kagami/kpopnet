import re
from contextlib import suppress
from datetime import datetime

from ._profile import ProfileSpider


class KprofilesSpider(ProfileSpider):
    name = 'kprofiles'
    start_urls = ['http://kprofiles.com/k-pop-girl-groups/']

    keep_only = {
        'I.O.I': set(['Sohye', 'Somi']),
        'I.B.I': set(['Haein', 'Hyeri', 'Sohee', 'Suhyun']),
        'Orange Caramel': set(),
        'Girls Next Door': set(),
        '4Minute': set(['Gayoon', 'Jihyun', 'Jiyoon', 'Sohyun']),
    }

    def parse(self, response):
        urls = response.css('.entry-content > p > a::attr(href)').extract()
        for url in sorted(set(urls)):
            if not url.endswith('-profile/'):
                continue
            if not self.update_all:
                try:
                    band = self.get_band_by_url(url)
                except KeyError:
                    if self.bnames:
                        continue
                else:
                    found = False
                    for bname in self.bnames:
                        if band['name'] == bname:
                            found = True
                            break
                    if not found:
                        continue
            meta = {'_knet_url': url}
            self.logger.warning('Crawling {}'.format(url))
            yield response.follow(url, self.parse_band, meta=meta)

    def parse_band(self, response):
        band = {}
        for p in response.css('.entry-content > p'):
            # First paragraph contains band info.
            if not band:
                name_node = p.css('img + br + strong::text') \
                    or p.css('img + br + b::text') \
                    or p.css('img + strong::text') \
                    or p.css('img + b::text') \
                    or p.xpath('strong/img/following-sibling::text()') \
                    or p.xpath('b/img/following-sibling::text()') \
                    or p.xpath('img/following::p/strong/text()') \
                    or p.xpath('following::p/img/following::p/strong/text()')
                name = self.normalize_band_name(name_node.extract_first())
                assert name, 'No band name'
                # TODO(Kagami): Parse more info.
                band['id'] = self.uuid()
                band['name'] = name
                band['urls'] = [response.meta['_knet_url']]
            # Idol info paragraph.
            elif p.css('span::text').\
                    re_first(r'(?i)(stage|real|birth)\s+name'):
                self.parse_idol(response, band, p)
        assert band, 'No band'
        self.save_band(band)
        # Technically not a warning but scrapy dumps a lot with INFO
        # level so we use WARNING to avoid that.
        self.logger.warning('Updated {}'.format(band['name']))

    def parse_idol(self, response, band, p):
        idol = {}
        for span in p.css('span'):
            key = span.css('::text').extract_first()
            if not key:
                continue
            val = span.xpath('following::text()').extract_first()
            # Try more aggressive search only if node was really empty.
            if not val or not val.strip():
                val = span.xpath('following::text()/following::text()').\
                    extract_first()
                if not val:
                    continue
            key = key.strip()
            val = val.strip()
            # Can check only this late.
            if not key.endswith(':') and not val.startswith(':'):
                continue
            key = self.normalize_idol_key(key)
            val = self.normalize_idol_val(key, val)
            if key and val:
                idol[key] = val

        assert idol.get('name'), 'No idol name'

        idol = self.normalize_idol(idol)
        with suppress(KeyError):
            if idol['name'] not in self.keep_only[band['name']]:
                return

        idol['id'] = self.uuid()
        idol['band_id'] = band['id']
        self.save_idol(band, idol)

    def normalize_band_name(self, name):
        orig = name
        name = name.strip()
        name = re.sub(r'\s*\(.*', '', name)
        name = re.sub(r'\s+', ' ', name)
        name = re.sub(r'’', "'", name)

        # Fix profile bugs.
        # TODO(Kagami): Normalize more special chars in names.
        if orig == 'F(x)':
            name = 'f(x)'
        elif name == 'Dal\u2605Shabet':
            name = 'Dal Shabet'
        elif name == 'Cosmic Girls':
            name = 'WJSN'
        elif name == 'Pristin':
            name = 'PRISTIN'

        return name

    def normalize_idol_key(self, key):
        key = key.lower()
        key = re.sub(r'\s*:$', '', key)
        key = re.sub(r'\s+', '_', key)
        if key == 'stage_name' or key == 'sage_name':
            key = 'name'
        elif key == 'real_name':
            key = 'birth_name'
        elif key == 'birthday':
            key = 'birth_date'
        elif key == 'position':
            key = 'positions'
        elif key == 'specialty' or key == 'speciality' or key == 'instruments':
            key = 'specialties'
        elif key == 'twitter_account':
            key = 'twitter'
        return key

    def normalize_idol_val(self, key, val):
        val = re.sub(r'^:\s*', '', val)
        val = re.sub(r'’', "'", val)

        if key == 'birth_date':
            try:
                val = datetime.strptime(val, '%B %d, %Y')
            except ValueError:
                val = None
        elif key == 'height':
            try:
                val = re.search(r'([\d.]+)\s*cm', val).group(1)
                val = float(val) if '.' in val else int(val)
            except AttributeError:
                val = None
        elif key == 'weight':
            try:
                val = re.search(r'([\d.]+)\s*kg', val).group(1)
                val = float(val) if '.' in val else int(val)
            except AttributeError:
                val = None
        elif key == 'zodiac_sign' or key == 'zodiac':
            val = None
        elif key == 'nationality':
            val = val.lower()
        elif key == 'positions' or key == 'specialties':
            val = [s.lower() for s in re.split(r',\s*', val)] if val else []
        elif key == 'twitter' or key == 'instagram':
            # @ uieing
            val = re.sub(r'@\s*', '', val)
            # @so_yul22 (she deactivated her account)
            val = re.sub(r'\s*\(.*', '', val)
        elif key.endswith('_facts'):
            val = None

        # Fix profile bugs.
        if key == 'name' and val.startswith('Jiin went on'):
            val = 'Jisun'
        elif key == 'name' and val == 'ROSÉ':
            val = 'Rose'
        elif key == 'birth_name' and val.startswith('Nam Ji Hyun, but she'):
            val = 'Nam Jihyun'

        return val

    def normalize_idol(self, idol):
        idol = idol.copy()
        name = idol['name']

        # Alt name var 1.
        with suppress(AttributeError):
            name, alt_name = re.\
                match(r'(.*?)\s*\(.*?known\s+as\s+(.*?)\)', name).\
                groups()
            idol['name'] = name
            idol['alt_names'] = [alt_name]

        # Alt name var 2.
        with suppress(AttributeError):
            name, alt_name = re.\
                match(r'(.*?)\s*/\s*(.*)', name).\
                groups()
            idol['name'] = name
            idol['alt_names'] = [alt_name]

        # Alt name var 3.
        with suppress(AttributeError):
            name, alt_name = re.\
                match(r'(.*?)\s+or\s+(.*)', name).\
                groups()
            idol['name'] = name
            idol['alt_names'] = [alt_name]

        # Hangul name.
        with suppress(AttributeError):
            name, name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', name).\
                groups()

            # Fix profile bugs.
            if name == 'EXY':
                name = 'Exy'

            idol['name'] = name
            idol['name_hangul'] = name_hangul

        # Birth name var 1.
        with suppress(KeyError, AttributeError):
            birth_name = idol['birth_name']
            birth_name, korean_name = re.\
                match(r'(.*?)\s*\{(.*?)\}', birth_name).\
                groups()
            idol['birth_name'] = birth_name
            idol['korean_name'] = korean_name

        # Birth name var 2.
        with suppress(KeyError, AttributeError):
            birth_name = idol['birth_name']
            birth_name, birth_name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', birth_name).\
                groups()

            # Fix profile bugs.
            if birth_name == 'Kang Yaebin':
                birth_name = 'Kang Yebin'

            idol['birth_name'] = birth_name
            idol['birth_name_hangul'] = birth_name_hangul

        # Korean name.
        with suppress(KeyError, AttributeError):
            korean_name = idol['korean_name']
            korean_name, korean_name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', korean_name).\
                groups()
            idol['korean_name'] = korean_name
            idol['korean_name_hangul'] = korean_name_hangul

        return idol
