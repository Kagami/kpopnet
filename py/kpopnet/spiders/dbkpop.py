from ._profile import ProfileSpider


class DbkpopSpider(ProfileSpider):
    name = 'dbkpop'
    start_urls = ['http://dbkpop.com/db/female-k-pop-idols/']
