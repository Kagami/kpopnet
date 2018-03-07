import json
from urllib.parse import urlencode

from scrapy.http import Request

from ._image import ImageSpider


class GoogleimagesSpider(ImageSpider):
    name = 'googleimages'
    custom_settings = {
        # NOTE(Kagami): Connected to search query parameters in some
        # way, do not change.
        'USER_AGENT': ('Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 '
                       '(KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36'),
    }

    def build_search_url(self, bname, mname):
        return 'https://www.google.com/search?' + urlencode([
            ('q', '{} {}'.format(bname, mname)),
            ('espv', '2'),
            ('biw', '1366'),
            ('bih', '667'),
            ('site', 'webhp'),
            ('source', 'lnms'),
            ('tbm', 'isch'),
            ('tbs', 'itp:face,ift:jpg'),
            ('sa', 'X'),
            ('ei', 'XosDVaCXD8TasATItgE'),
            ('ved', '0CAcQ_AUoAg'),
        ])

    def start_requests(self):
        for bname, mname in self.get_all_member_names():
            if not self.update_all and self.has_images_by_name(bname, mname):
                continue
            url = self.build_search_url(bname, mname)
            meta = {'_knet_bname': bname, '_knet_mname': mname}
            yield Request(url, dont_filter=True, meta=meta)

    def parse(self, response):
        response.meta['_knet_num_saved'] = 0
        response.meta['_knet_items'] = response.css(".rg_meta")
        # TODO(Kagami): Speedup download by emitting MAX_N requests at
        # the same time and fetching remaining after that?
        return self.parse_item(response.meta)

    def parse_item(self, meta):
        if meta['_knet_num_saved'] >= self.MAX_IMAGES_PER_MEMBER:
            return
        try:
            item = meta['_knet_items'].pop(0)
        except IndexError:
            self.logger.warning('Saved only {} images for {} - {}'.format(
                meta['_knet_num_saved'], meta['_knet_bname'],
                meta['_knet_mname']))
            return
        try:
            item_text = item.css("::text").extract_first()
            item_info = json.loads(item_text)
            image_url = item_info['ou']
        except Exception:
            return
        else:
            yield Request(image_url, self.save_item,
                          errback=self.errback_item,
                          dont_filter=True,
                          meta=meta)

    def save_item(self, response):
        if self.save_image(response):
            response.meta['_knet_num_saved'] += 1
        return self.parse_item(response.meta)

    def errback_item(self, failure):
        return self.parse_item(failure.request.meta)
