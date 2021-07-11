import subprocess
import time

subprocess.run(["cargo", "build", "--release", "--target", "x86_64-pc-windows-msvc"])
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "a", "1", "30"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "b", "30", "60"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "c", "60", "85"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "d", "85", "100"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "e", "100", "115"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "f", "110", "120"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "g", "120", "126"])
time.sleep(2)
subprocess.Popen(["./target/x86_64-pc-windows-msvc/release/application_a", "h", "126", "133"])
