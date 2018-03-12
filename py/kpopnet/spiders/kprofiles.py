import re
from contextlib import suppress
from datetime import datetime

from ._profile import ProfileSpider


class KprofilesSpider(ProfileSpider):
    name = 'kprofiles'
    start_urls = ['http://kprofiles.com/k-pop-girl-groups/']

    def parse(self, response):
        urls = response.css('.entry-content > p > a::attr(href)').extract()
        for url in sorted(set(urls)):
            if not self.update_all and self.has_band_by_url(url):
                continue
            if url.endswith('-profile/'):
                meta = {'_knet_url': url}
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
                    or p.xpath('img/following::p/strong/text()') \
                    or p.xpath('following::p/img/following::p/strong/text()')
                name = self.normalize_band_name(name_node.extract_first())
                assert name, 'No band name'
                # TODO(Kagami): Parse more info.
                band['id'] = self.uuid()
                band['name'] = name
                band['urls'] = [response.meta['_knet_url']]
            # Member info paragraph.
            elif p.css('span::text').\
                    re_first(r'(?i)(stage|real|birth)\s+name'):
                self.parse_member(response, band, p)
        assert band, 'No band'
        self.save_band(band)

    def parse_member(self, response, band, p):
        member = {}
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
            key = self.normalize_member_key(key)
            val = self.normalize_member_val(key, val)
            if key and val:
                member[key] = val
        # if not member.get('name'):
        #     from IPython import embed; embed()
        assert member.get('name'), 'No member name'
        member = self.normalize_member(member)
        member['id'] = self.uuid()
        member['band_id'] = band['id']
        self.save_member(band, member)

    def normalize_band_name(self, name):
        orig = name
        name = name.strip()
        name = re.sub(r'\s*\(.*', '', name)
        name = re.sub(r'\s+', ' ', name)
        name = re.sub(r'’', "'", name)

        # Fix profile bugs.
        # TODO(Kagami): Better normalization.
        if orig == 'F(x)':
            name = 'f(x)'
        elif name == 'Dal\u2605Shabet':
            name = 'Dal Shabet'

        return name

    def normalize_member_key(self, key):
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
        elif key == 'zodiac_sign':
            key = 'zodiac'
        elif key == 'twitter_account':
            key = 'twitter'
        return key

    def normalize_member_val(self, key, val):
        val = re.sub(r'^:\s*', '', val)
        val = re.sub(r'’', "'", val)

        if key == 'name':
            # Ara / Yooara
            val = re.sub(r'\s*/.*', '', val)
        elif key == 'birth_date':
            try:
                val = datetime.strptime(val, '%B %d, %Y')
            except ValueError:
                val = None
        elif key == 'height':
            try:
                val = int(re.search(r'(\d+)\s*cm', val).group(1))
            except AttributeError:
                val = None
        elif key == 'weight':
            try:
                val = int(re.search(r'(\d+)\s*kg', val).group(1))
            except AttributeError:
                val = None
        elif key == 'zodiac':
            val = val.lower()
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

        return val

    def normalize_member(self, member):
        member = member.copy()

        # Name field must always present.
        with suppress(AttributeError):
            name = member['name']
            name, name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', name).\
                groups()

            # Fix profile bugs.
            if name == 'EXY':
                name = 'Exy'

            member['name'] = name
            member['name_hangul'] = name_hangul

        # Normalize birth name if any.
        with suppress(KeyError, AttributeError):
            birth_name = member['birth_name']
            birth_name, birth_name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', birth_name).\
                groups()
            member['birth_name'] = birth_name
            member['birth_name_hangul'] = birth_name_hangul

        # Normalize korean name if any.
        with suppress(KeyError, AttributeError):
            korean_name = member['korean_name']
            korean_name, korean_name_hangul = re.\
                match(r'(.*?)\s*\((.*?)\)', korean_name).\
                groups()
            member['korean_name'] = korean_name
            member['korean_name_hangul'] = korean_name_hangul

        return member
