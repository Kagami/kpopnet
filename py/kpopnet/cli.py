"""
K-pop neural network spiders and utilities.

Usage:
  kpopnet profile update [BAND...] [options]
  kpopnet image update [options]
  kpopnet [-h | --help]
  kpopnet [-V | --version]

Options:
  -h --help     Show this screen.
  -V --version  Show version.
  --debug       Enable debug logging.
  -s SPIDER     Select spider.
  --all         Update already collected data.
  --bail        Exit on first error.
"""

import sys
import pkg_resources

from docopt import docopt


def main():
    version = pkg_resources.require('kpopnet')[0].version
    args = docopt(__doc__, version=version)

    if args['profile']:
        from . import profiles
        if args['update']:
            return profiles.update(
                args['-s'] or 'kprofiles',
                update_all=args['--all'],
                bnames=args['BAND'],
                bail=True,
                debug=args['--debug'])
    elif args['image']:
        from . import images
        if args['update']:
            return images.update(
                args['-s'] or 'googleimages',
                update_all=args['--all'],
                bail=args['--bail'],
                debug=args['--debug'])

    print('No command selected, try --help.', file=sys.stderr)
    sys.exit(1)
