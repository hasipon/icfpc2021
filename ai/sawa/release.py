import subprocess
import time

subprocess.run(["cargo", "build", "--release", "--target", "x86_64-pc-windows-msvc"])
for n in range(0, 11):
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "a" + str(n), str(n * 7 + 56), "133"])
