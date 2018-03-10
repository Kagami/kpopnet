from .spiders import image_spiders, run_spider


def update(spider_name, **kwargs):
    """
    Collect/update idol images using given spider and write them to
    data directory. Safe to use multiple times, previously collected
    data will be preserved.
    """
    spider = image_spiders[spider_name]
    return run_spider(spider, **kwargs)
