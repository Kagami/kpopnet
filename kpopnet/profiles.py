from .spiders import profile_spiders, run_spider


def update(spider_name, **kwargs):
    """
    Collect/update profiles info using given spider and write them to
    data directory. Safe to use multiple times, previously collected
    data will be preserved.
    """
    spider = profile_spiders[spider_name]
    return run_spider(spider, **kwargs)
