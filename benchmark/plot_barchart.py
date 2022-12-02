#!/usr/bin/env python

import argparse
import json
import matplotlib.pyplot as plt

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("file", help="JSON file with benchmark results")
parser.add_argument("--title", help="Plot title")
parser.add_argument("--out", help="Outfile file")
args = parser.parse_args()

with open(args.file) as f:
    results = json.load(f)["results"]

x = ["daselv2", "dasel", "jq", "yq"]
x_pos = [i for i, _ in enumerate(x)]

# mean is in seconds. convert to ms.
mean_times = [b["mean"] * 1000 for b in results]

plt.bar(x_pos, mean_times, color='green', align='center')

plt.ylabel("Execution Time (ms)")
if args.title is not None:
    plt.title(args.title)

plt.xticks(x_pos, x, horizontalalignment='center')

plt.savefig(args.out, dpi=None, facecolor='w', edgecolor='w',
            orientation='portrait', format=None,
            transparent=False, bbox_inches=None, pad_inches=0.1)
