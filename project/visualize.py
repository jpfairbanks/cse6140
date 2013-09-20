""" visualize.py: make a plot of the scaling of a function
Author: James Fairbanks
Date: 2013-09-19
"""

import pandas as pd
import matplotlib
matplotlib.use("pdf")
import matplotlib.pyplot as plt
import numpy as np
import sys

fp = sys.stdin
df = pd.read_csv(fp, sep=" ", index_col=0, header=None)
fig = df.plot()
print(df)
plt.savefig("scaling.pdf")