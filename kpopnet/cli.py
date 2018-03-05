"""
K-pop face recognition neural network and utilities.

Usage:
  kpopnet profiles update [-s <spider>] [--all] [--bail]
  kpopnet [-h | --help]
  kpopnet [-V | --version]

Options:
  -h --help     Show this screen.
  -V --version  Show version.
  -s <spider>   Select spider [default: kprofiles].
  --all         Update already collected profiles.
  --bail        Exit on first error.
"""

import pkg_resources

from docopt import docopt

from . import profiles


def main():
    version = pkg_resources.require('kpopnet')[0].version
    args = docopt(__doc__, version=version)

    if args['profiles']:
        if args['update']:
            return profiles.update(
                args['-s'],
                update_all=args['--all'],
                bail=args['--bail'])
