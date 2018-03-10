from ._profile import ProfileSpider


class WikipediaSpider(ProfileSpider):
    name = 'wikipedia'
    start_urls = [
        'https://en.wikipedia.org/wiki/List_of_South_Korean_idol_groups'
    ]
