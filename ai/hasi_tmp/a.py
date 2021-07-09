import json
import sys

hole_x = []
hole_y = []
figure_x = []
figure_y = []
for filepath in sys.argv[1:]:
    with open(filepath) as f:
        a = json.load(f)
        for x, y in a["hole"]:
            hole_x.append(x)
            hole_y.append(y)
        for x, y in a["figure"]["vertices"]:
            figure_x.append(x)
            figure_y.append(y)

print("min hole_x", min(hole_x))
print("max hole_x", max(hole_x))
print("min hole_y", min(hole_y))
print("max hole_y", max(hole_y))
print("min figure_x", min(figure_x))
print("max figure_x", max(figure_x))
print("min figure_y", min(figure_y))
print("max figure_y", max(figure_y))
