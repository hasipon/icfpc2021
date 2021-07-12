import subprocess
import time

subprocess.run(["cargo", "build", "--release", "--target", "x86_64-pc-windows-msvc"])

for n in range(0, 11):
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "d" + str(n), str(n * 8 + 46), "133"])


#subprocess.run(["./target/x86_64-pc-windows-msvc/release/application_a", "a" + str(99), str(1), "133"])