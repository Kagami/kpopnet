from scrapy import signals
from scrapy.crawler import CrawlerProcess

from .kprofiles import KprofilesSpider
from .nowkpop import NowkpopSpider
from .kpopinfo114 import Kpopinfo114Spider
from .wikipedia import WikipediaSpider
from .googleimages import GoogleimagesSpider


USER_AGENT = (
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) '
    'AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36'
)


profile_spiders = {}
for spider in [
    KprofilesSpider,
    NowkpopSpider,
    Kpopinfo114Spider,
    WikipediaSpider,
]:
    profile_spiders[spider.name] = spider


image_spiders = {}
for spider in [
    GoogleimagesSpider,
]:
    image_spiders[spider.name] = spider


class HttpErrorMiddleware:
    def process_spider_input(self, response, spider):
        if not 200 <= response.status < 300:
            raise Exception('HTTP error {}'.format(response))


def run_spider(spider, bail=False, **kwargs):

    def process_spider_error(failure, response, spider):
        nonlocal had_error
        had_error = True

    had_error = False
    spider_middlewares = {}
    if bail:
        spider_middlewares['kpopnet.spiders.HttpErrorMiddleware'] = 1
    process = CrawlerProcess({
        'USER_AGENT': USER_AGENT,
        'CLOSESPIDER_ERRORCOUNT': 1 if bail else 0,
        'CONCURRENT_REQUESTS_PER_DOMAIN': 1,
        'SPIDER_MIDDLEWARES': spider_middlewares,
    })
    crawler = process.create_crawler(spider)
    crawler.signals.connect(process_spider_error, signals.spider_error)
    process.crawl(crawler, **kwargs)
    process.start()
    return 1 if had_error else 0
