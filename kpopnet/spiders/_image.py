import io

import scrapy
from PIL import Image

from ..io import has_images_by_name, get_all_member_names, save_image_by_name


def is_valid_jpeg(data):
    try:
        assert data[6:10] in (b'JFIF', b'Exif')
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

    MAX_IMAGES_PER_MEMBER = 10

    def has_images_by_name(self, bname, mname):
        return has_images_by_name(bname, mname)

    def get_all_member_names(self):
        return get_all_member_names()

    def save_image(self, response):
        # TODO(Kagami): Ensure there is only single face in the image.
        if not is_valid_jpeg(response.body):
            return False
        bname = response.meta['_knet_bname']
        mname = response.meta['_knet_mname']
        return save_image_by_name(bname, mname, response.body)
