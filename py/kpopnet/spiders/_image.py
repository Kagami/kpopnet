import io

import scrapy
from PIL import Image

from ..io import has_images_by_name, get_all_idol_names, save_image_by_name


def is_valid_jpeg(data):
    try:
        assert data.startswith(b'\xff\xd8')
        im = Image.open(io.BytesIO(data))
        assert im.format == 'JPEG'
        assert im.width < 5000
        assert im.height < 5000
        im.load()
    except Exception:
        return False
    else:
        return True


class ImageSpider(scrapy.Spider):
    """
    Useful wrappers and common image spider methods.
    """

    MAX_IMAGES_PER_IDOL = 10

    def has_images_by_name(self, bname, iname):
        return has_images_by_name(bname, iname)

    def get_all_idol_names(self):
        return get_all_idol_names()

    def save_image(self, response):
        # TODO(Kagami): Ensure there is only single face in the image.
        # Also make sure its resolution is good enough.
        if not is_valid_jpeg(response.body):
            return False
        bname = response.meta['_knet_bname']
        iname = response.meta['_knet_iname']
        return save_image_by_name(bname, iname, response.body)
