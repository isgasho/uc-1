define void @f() {
0:
	%a = alloca [5 x i32]
	%1 = getelementptr [5 x i32], [5 x i32]* %a, i64 0, i64 2
	%2 = load i32, i32* %1
	ret void
}
