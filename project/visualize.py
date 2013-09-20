""" visualize.py: make a plot of the scaling of a function
Author: James Fairbanks
Date: 2013-09-19
"""

import pandas as pd
import matplotlib
matplotlib.use("pdf")
import matplotlib.pyplot as plt
import sys


def main(logy):
    """we just plot the data off stdin"""
    filep = sys.stdin
    dataf = pd.read_csv(filep, sep=" ", index_col=0, header=None)
    dataf.plot(logy=logy)
    print(dataf)
    plt.savefig("scaling.pdf")

if __name__ == '__main__':
    Logscale = False
    if len(sys.argv) >=2:
        Logscale = sys.argv[1]=="-log"
    main(logy=Logscale)
