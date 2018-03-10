from ._profile import ProfileSpider


class NowkpopSpider(ProfileSpider):
    name = 'nowkpop'
    start_urls = ['https://www.nowkpop.com/category/k-popprofiles/']
