import pandas as pd
import matplotlib
matplotlib.use("pdf")
import matplotlib.pyplot as plt
import sys

df = pd.read_csv(sys.stdin, " ", header=None, index_col=0)
print(df)
print(df["X2"]/2**df.index)
df.plot(logy=True)
plt.savefig("graph.pdf")
