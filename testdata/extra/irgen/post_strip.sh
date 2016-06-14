#!/bin/bash
for f in "function.ll" "local.ll"; do
	sar -i "\t%1 = alloca i32\n" "" "${f}"
	sar -i "\tstore i32 0, i32[*] %1\n" "" "${f}"
done
