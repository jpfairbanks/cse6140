import pandas as pd
import matplotlib
matplotlib.use("pdf")
import matplotlib.pyplot as plt
import sys
import os.path

figdir = "figures"

df = pd.read_csv(sys.stdin, " ", header=None, index_col=0,
                 names=[2**(i) for i in range(6)]+["rand"])
print(df)
#print(df["X2"]/2**df.index)
df.plot(logy=True)
plt.title("run time for array access")
plt.xlabel("scale")
plt.ylabel("seconds")
plt.savefig(os.path.join([figdir,"graph.pdf"]))

plt.figure()
sizes = 2**df.index
print(sizes)
petf = (df.T/sizes).T
print( petf )
petf.plot(logy=True)
plt.title("normalized running time")
plt.xlabel("scale")
plt.ylabel("nanoseconds per element")
plt.savefig(os.path.join([figdir,"perelement.pdf"]))
