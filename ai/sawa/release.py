import subprocess

subprocess.run(["cargo", "build", "--release", "--target", "x86_64-pc-windows-msvc"])
subprocess.run(["./target/x86_64-pc-windows-msvc/release/application_a"])
