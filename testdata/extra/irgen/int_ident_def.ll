define void @f() {
0:
	%a = alloca i32
	store i32 5, i32* %a
	ret void
}
