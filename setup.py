from setuptools import setup, find_packages


setup(
    name='kpopnet',
    version='0.0.0',
    author='Kagami Hiiragi',
    author_email='kagami@genshiken.org',
    url='https://github.com/Kagami/kpopnet',
    description='K-pop face recognition neural network',
    license='CC0',
    package_dir={'': 'py'},
    packages=find_packages('py', exclude=['tests']),
    entry_points={
        'console_scripts': ['kpopnet = kpopnet.cli:main'],
    },
    install_requires=[
        'docopt>=0.6.2',
        'scrapy>=1.5.0',
        'Pillow>=5.0.0',
    ],
    extras_require={
      'tests': ['flake8'],
    },
    classifiers=[
        'Development Status :: 3 - Alpha',
        'License :: CC0 1.0 Universal (CC0 1.0) Public Domain Dedication',
        'Operating System :: OS Independent',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
    ],
)
