import subprocess
import time

subprocess.run(["cargo", "build", "--release", "--target", "x86_64-pc-windows-msvc"])
for n in range(0, 5):
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "a" + str(n), "1", "30"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "b" + str(n), "30", "60"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "c" + str(n), "60", "85"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "d" + str(n), "85", "100"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "e" + str(n), "100", "115"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "f" + str(n), "110", "120"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "g" + str(n), "120", "126"])
    time.sleep(2)
    subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "h" + str(n), "126", "133"])
