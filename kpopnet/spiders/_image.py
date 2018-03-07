import scrapy

from ..io import has_images_by_name, get_all_member_names, save_image_by_name


def has_jpeg_header(data):
    return data[6:10] in (b'JFIF', b'Exif')


class ImageSpider(scrapy.Spider):
    """
    Just a collection of useful wrappers. We may pass settings object in
    the future to support e.g. custom locations of image data.
    """

    MAX_IMAGES_PER_MEMBER = 10

    def has_images_by_name(self, bname, mname):
        return has_images_by_name(bname, mname)

    def get_all_member_names(self):
        return get_all_member_names()

    def save_image(self, response):
        # TODO(Kagami): File might be corrupted so probably better to
        # fully decode?
        if not has_jpeg_header(response.body):
            return False
        bname = response.meta['_knet_bname']
        mname = response.meta['_knet_mname']
        save_image_by_name(bname, mname, response.body)
        return True
