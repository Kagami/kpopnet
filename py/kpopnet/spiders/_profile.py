from uuid import uuid4

import scrapy

from ..io import has_band_by_url, save_band, save_member


class ProfileSpider(scrapy.Spider):
    """
    Just a collection of useful wrappers. We may pass settings object in
    the future to support e.g. custom locations of profile data.
    """

    def uuid(self):
        return str(uuid4())

    def has_band_by_url(self, band):
        return has_band_by_url(band)

    def save_band(self, band):
        return save_band(band)

    def save_member(self, band, member):
        return save_member(band, member)
