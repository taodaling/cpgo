# cpgo

It's a server for competitive companion, a test framework for C++, and implement simple inline mechanism for competitive programming.

# How to use it 

Compile the source at first, it's easy:

```sh
# go build
```

Run the binary file in shell 

```sh
# cpgo 2>/dev/null
```

Compile your source file `main.cpp` with output file `main` or `main.exe`.

Then the server will automatically detect your compiled file changed and run whole tests for it. At the same time, all file included by `#include "your header"` will be inlined
into a new file named `inline.cpp`. You can submit this file with all library included.

# How to parse task from web page

Install competitive companion and add a new port `50823` to it.

# Where can I get library

There are lots of cool libraries, I'll mention some:

- [ac-library](https://github.com/atcoder/ac-library)
- [kactl](https://github.com/kth-competitive-programming/kactl)
