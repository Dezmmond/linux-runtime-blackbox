#!/usr/bin/env python3
import tempfile
import time

def main():
    files = []

    for _ in range(5):
        f = tempfile.NamedTemporaryFile(mode="w+", delete=True)
        f.write("blackbox test\n")
        f.flush()
        files.append(f)

    time.sleep(5)

if __name__ == "__main__":
    main()
