package binchunk

//➜  ch02 xxd -u -g 1 hello.out
//00000000: 1B 4C 75 61 53 00 19 93 0D 0A 1A 0A 04 08 04 08  .LuaS...........
//00000010: 08 78 56 00 00 00 00 00 00 00 00 00 00 00 28 77  .xV...........(w
//00000020: 40 01 11 40 68 65 6C 6C 6F 5F 77 6F 72 6C 64 2E  @..@hello_world.
//00000030: 6C 75 61 00 00 00 00 00 00 00 00 00 01 02 04 00  lua.............
//00000040: 00 00 06 00 40 00 41 40 00 00 24 40 00 01 26 00  ....@.A@..$@..&.
//00000050: 80 00 02 00 00 00 04 06 70 72 69 6E 74 04 0E 48  ........print..H
//00000060: 65 6C 6C 6F 2C 20 57 6F 72 6C 64 21 01 00 00 00  ello, World!....
//00000070: 01 00 00 00 00 00 04 00 00 00 01 00 00 00 01 00  ................
//00000080: 00 00 01 00 00 00 01 00 00 00 00 00 00 00 01 00  ................
//00000090: 00 00 05 5F 45 4E 56                             ..._ENV

//header 中的常量
const (
	LUA_SINATURE     = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

//tag 表类型常量
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header                 // 头部
	sizeUpvalue byte       // 主函数 upvalue 数量
	mainFunc    *prototype // 主函数原型
}

type header struct {
	signature       [4]byte // 0x1b4c7561 '1B 4C 75 61' Magic Number ESC L u a 的ASCII
	version         byte    // 0x53 '53' 5.3.4 5*16+3=83
	format          byte    // 0x00 '00'
	luacData        [6]byte // 0x1993 '19 93 0D 0A 1A 0A' "\x19\x93\r\n\xla\n"
	cintSize        byte    // 0x04 '04' cint 类型占用字节数
	sizezSize       byte    // 0x08 '08'
	instructionSize byte    // 0x04 '04'
	luaIntegerSize  byte    // 0x08 '08'
	luaNumberSize   byte    // 0x08 '08'
	luacInt         int64   //0x5678 '78 56 00 00 00 00 00 00' 判断机器大小端
	luacNum         float64 //370.5 '00 00 00 00 00 28 77 40' 所使用浮点数格式 IEEE754
}

//prototype 函数原型
type prototype struct {
	Source          string        // @hello_world.lua '11 40 68 65 6C 6C 6F 5F 77 6F 72 6C 64 2E' 源文件名 只有主函数里有 0x11 代表字符串长度+1
	LineDefined     uint32        // '00 00 00 00' 函数起止行号 主函数都为0 非主函数都大于0
	LastLineDefined uint32        // '00 00 00 00'
	NumParams       byte          // '00' 函数固定参数个数 生成的主函数没有 所以为0
	IsVararg        byte          // '01' 是否为变长参数函数 主函数是变长参数函数 所以为1
	MaxStackSize    byte          // '02' 寄存器数量
	Code            []uint32      // '04 00 00 00  06 00 40 00  41 40 00 00  24 40 00 01  26 00 80 00'指令表 每条指令四个字节
	Constants       []interface{} // '02 00 00 00   04 06 70 72 69 6E 74  04 0E 48 65 6C 6C 6F 2C 20 57 6F 72 6C 64 21' 常量表
	Upvalues        []Upvalue     // '01 00 00 00    01 00' upvalue表 每个值占用两字节
	Protos          []*prototype  // '00 00 00 00' 子函数原型表
	LineInfo        []uint32      // '04 00 00 00  01 00 00 00  01 00 00 00  01 00 00 00  01 00 00 00 ' 行号表 指令对应的行号
	LocVars         []LocVar      // '00 00 00 00' 局部变量长度为0
	UpvalueNames    []string      // '01 00 00 00     05 5F 45 4E 56' _ENV upvalue name表
}

//Upvalue 闭包引用
type Upvalue struct {
	Instack byte
	Idx     byte
}

//LocVar 局部变量
type LocVar struct {
	VarName string
	StartPC uint32 //起止指令索引
	EndPC   uint32
}
