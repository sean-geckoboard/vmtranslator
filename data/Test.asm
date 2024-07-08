// push 7
@7
D=A
@SP
A=M
M=D
@SP
M=M+1

// pop
@1
D=A
// at argument
@2
A=D+M
D=A
@13
M=D

@SP
M=M-1
A=M
D=M
@13
A=M
M=D
