"""
K-pop face recognition neural network and utilities.

Usage:
  kpopnet profiles update [options]
  kpopnet images update [options]
  kpopnet [-h | --help]
  kpopnet [-V | --version]

Options:
  -h --help     Show this screen.
  -V --version  Show version.
  -s <spider>   Select spider.
  --all         Update already collected data.
  --bail        Exit on first error.
"""

import pkg_resources

from docopt import docopt

from . import profiles
from . import images


def main():
    version = pkg_resources.require('kpopnet')[0].version
    args = docopt(__doc__, version=version)

    if args['profiles']:
        if args['update']:
            return profiles.update(
                args['-s'] or 'kprofiles',
                update_all=args['--all'],
                bail=args['--bail'])
    elif args['images']:
        if args['update']:
            return images.update(
                args['-s'] or 'googleimages',
                update_all=args['--all'],
                bail=args['--bail'])
