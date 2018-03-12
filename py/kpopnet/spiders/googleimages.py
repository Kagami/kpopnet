import json
from contextlib import suppress
from urllib.parse import urlencode

from scrapy.http import Request

from ._image import ImageSpider


class GoogleimagesSpider(ImageSpider):
    name = 'googleimages'
    custom_settings = {
        # NOTE(Kagami): Connected to query parameters in some way, do
        # not change.
        'USER_AGENT': ('Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 '
                       '(KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36'),
        'DNS_TIMEOUT': 20,
        'DOWNLOAD_TIMEOUT': 30,
        'DOWNLOAD_MAXSIZE': 2 * 1024 * 1024,
        'REFERRER_POLICY': 'no-referrer',
    }

    def build_search_url(self, bname, iname):
        return 'https://www.google.com/search?' + urlencode([
            ('q', '{} {}'.format(bname, iname)),  # Search query
            ('espv', '2'),
            ('biw', '1366'),
            ('bih', '667'),
            ('site', 'webhp'),
            ('source', 'lnms'),
            ('tbm', 'isch'),
            ('tbs', 'itp:face,ift:jpg'),  # Face type, JPEG format
            ('sa', 'X'),
            ('ei', 'XosDVaCXD8TasATItgE'),
            ('ved', '0CAcQ_AUoAg'),
        ])

    def start_requests(self):
        self.search_reqs = []
        for bname, iname in self.get_all_idol_names():
            if not self.update_all and self.has_images_by_name(bname, iname):
                continue
            url = self.build_search_url(bname, iname)
            meta = {'_knet_bname': bname, '_knet_iname': iname}
            req = Request(url, dont_filter=True, meta=meta, priority=-1)
            self.search_reqs.append(req)
        with suppress(IndexError):
            yield self.search_reqs.pop(0)

    def parse(self, response):
        # NOTE(Kagami): Has to be mutable because errbacks receive
        # shallow copy of metadata.
        response.meta['_knet_reqs'] = {'all': 0, 'saved': 0}
        response.meta['_knet_items'] = response.css('.rg_meta')
        # Do first N requests in parallel.
        response.meta['_knet_parallel'] = True
        for _ in range(self.MAX_IMAGES_PER_IDOL):
            yield self.parse_item(response.meta)
        with suppress(IndexError):
            yield self.search_reqs.pop(0)

    def parse_item(self, meta):
        if meta['_knet_reqs']['saved'] >= self.MAX_IMAGES_PER_IDOL:
            return
        try:
            item = meta['_knet_items'].pop(0)
        except IndexError:
            self.logger.warning('Saved only {} images for {} - {}'.format(
                meta['_knet_reqs']['saved'],
                meta['_knet_bname'],
                meta['_knet_iname']))
            return
        try:
            item_text = item.css('::text').extract_first()
            item_info = json.loads(item_text)
            image_url = item_info['ou']
        except Exception:
            return self.parse_item(meta)
        else:
            return Request(image_url, self.callback_item,
                           errback=self.errback_item,
                           dont_filter=True,
                           meta=meta)

    def callback_item(self, response):
        if self.save_image(response):
            response.meta['_knet_reqs']['saved'] += 1
        else:
            self.logger.warning('Skipped image {}'.format(response))
        return self.finally_item(response.meta)

    def errback_item(self, failure):
        return self.finally_item(failure.request.meta)

    def finally_item(self, meta):
        meta['_knet_reqs']['all'] += 1
        if meta['_knet_parallel']:
            if meta['_knet_reqs']['all'] < self.MAX_IMAGES_PER_IDOL:
                return
            # Fetch remaining.
            meta['_knet_parallel'] = False
        return self.parse_item(meta)
