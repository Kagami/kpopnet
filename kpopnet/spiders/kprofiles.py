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
                meta = {'_kpopnet_url': url}
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
                band['name'] = name
                band['urls'] = [response.meta['_kpopnet_url']]
            # Member info paragraph.
            elif p.css('span::text').re_first(r'(?i)stage\s+name:'):
                self.parse_member(response, band, p)
        assert band, 'No band'
        self.save_band(band)

    def parse_member(self, response, band, p):
        member = {}
        for span in p.css('span'):
            key = span.css('::text').extract_first()
            if not key or not key.strip().endswith(':'):
                continue
            key = self.normalize_member_key(key)
            # TODO(Kagami): Parse twitter/instagram.
            val = span.xpath('following::text()').extract_first()
            if not val or not val.strip():
                val = span.xpath('following::text()/following::text()').\
                    extract_first()
            if not val or not val.strip():
                continue
            val = self.normalize_member_val(key, val)
            member[key] = val
        member = self.normalize_member(member)
        assert member.get('name'), 'No member name'
        self.save_member(band, member)

    def normalize_band_name(self, name):
        name = name.strip()
        with suppress(AttributeError):
            name = re.match(r'(.*)\s+\(', name).group(1)
        name = re.sub(r'’', "'", name)

        # Fix profile bugs.
        if name == 'F(x)':
            name = 'f(x)'

        return name

    def normalize_member_key(self, key):
        key = key.strip()
        key = key.lower()
        key = re.sub(r':$', '', key)
        key = re.sub(r'\s+', '_', key)
        if key == 'zodiac_sign':
            key = 'zodiac'
        elif key == 'position':
            key = 'positions'
        elif key == 'specialty':
            key = 'specialties'
        elif key == 'birthday':
            key = 'birth_date'
        return key

    def normalize_member_val(self, key, val):
        val = val.strip()
        val = re.sub(r'’', "'", val)

        if key == 'birth_date':
            with suppress(ValueError):
                val = datetime.strptime(val, '%B %d, %Y')
        elif key == 'zodiac':
            val = val.lower()
        elif key == 'nationality':
            val = val.lower()
        elif key == 'height':
            with suppress(AttributeError):
                val = int(re.search(r'(\d+)\s+cm', val).group(1))
        elif key == 'weight':
            with suppress(AttributeError):
                val = int(re.search(r'(\d+)\s+kg', val).group(1))
        elif key == 'positions' or key == 'specialties':
            val = [s.lower() for s in re.split(r',\s+', val)] if val else []

        # Fix profile bugs.
        if key == 'stage_name' and val.startswith('Jiin went on'):
            val = 'Jisun'

        return val

    def normalize_member(self, member):
        member = member.copy()

        # Name field must always present.
        name = member.pop('stage_name')
        with suppress(AttributeError):
            name, name_hangul = re.\
                match(r'(.*)\s+\((.*)\)', name).\
                groups()
            member['name_hangul'] = name_hangul
        member['name'] = name

        # Normalize birth name if any.
        try:
            birth_name = member['birth_name']
        except KeyError:
            pass
        else:
            with suppress(AttributeError):
                birth_name, birth_name_hangul = re.\
                    match(r'(.*)\s+\((.*)\)', birth_name).\
                    groups()
                member['birth_name'] = birth_name
                member['birth_name_hangul'] = birth_name_hangul

        return member
