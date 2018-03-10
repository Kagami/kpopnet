from ._profile import ProfileSpider


class Kpopinfo114Spider(ProfileSpider):
    name = 'kpopinfo114'
    start_urls = ['https://kpopinfo114.wordpress.com/female_artist_profiles/']
