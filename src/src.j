func main() {
    test = "Hello"
    print(test, "world")
}

var test = 123

// start:
//   PUSH rsb
//   MOV rbp, rsp
//   SUB rsp, 0x20
//   CALL main
//   MOV rax, 0
//   CALL ExitProcess
//   LEAVE
//   RET
//
// main:
//   PUSH rbp
//   MOV rbp, rsp
//   SUB rsp, 0x20
//   MOV QWORD [rbp+0x8], data0
//   MOV rcx, [rbp+0x8]
//   MOV rdx, data1
//   CALL print
//   LEAVE
//   RET
//
// print:
//   PUSH rbp
//   MOV rbp, rsp
//   SUB rsp, 0x20
//
//   MOV r9, 0
//   MOV r8, rdx
//   MOV rdx, rcx
//   MOV rcx, 0
//   CALL MessageBoxA
//   LEAVE
//   RET