package pe

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"time"
)

type PE struct {
	buf bytes.Buffer

	LinkerVersion           [2]byte
	SizeOfInitializedData   int32
	SizeOfUninitializedData int32
	addressOfEntryPoint     int32
	baseOfCode              int32
	imageBase               int64
	sectionAlignment        int32
	fileAlignment           int32

	CheckSum           int32
	sizeOfStackReserve int64
	sizeOfStackCommit  int64
	sizeOfHeapReserve  int64
	sizeOfHeapCommit   int64

	sections []*section
}

func New(sectionAlignment int, fileAlignment int) *PE {
	if fileAlignment < 512 || fileAlignment > 65536 {
		panic("FileAlignment must be between 512 and 65536 (inclusive)")
	}
	if !isPow2(fileAlignment) {
		panic("FileAlignment must be power of 2")
	}
	if sectionAlignment < fileAlignment {
		panic("SectionAlignment must be greater than or equal to FileAlignment")
	}

	pe := &PE{
		fileAlignment:       int32(fileAlignment),
		sectionAlignment:    int32(sectionAlignment),
		imageBase:           0x400000,
		baseOfCode:          0x1000,
		addressOfEntryPoint: 0x1000,

		sizeOfStackReserve: 0x100000,
		sizeOfStackCommit:  0x1000,
		sizeOfHeapReserve:  0x100000,
		sizeOfHeapCommit:   0x1000,
	}
	pe.sections = []*section{
		newSection(pe, ".text", 0x60000020),
		newSection(pe, ".idata", 0x40000040),
	}

	return pe
}

func (p *PE) actualSizeOfHeaders() int32 {
	// DOS header: 0x40 bytes
	// PE header: 0x04 bytes
	// File header: 0x14 bytes
	// Optional header: 0xF0 bytes
	// Section size: 40 bytes
	return int32(0x40 + 0x04 + 0x14 + 0xF0 + 40*len(p.sections))
}

func (p *PE) sizeOfHeaders() int32 {
	return multipleOf(
		p.actualSizeOfHeaders(),
		p.fileAlignment,
	)
}
func (p *PE) sizeOfCode() (s int32) {
	return p.sections[0].header.SizeOfRawData
}
func (p *PE) sizeOfImage() int32 {
	return p.sections[len(p.sections)-1].header.VirtualAddress + p.sectionAlignment
}

func (p *PE) writeDOSHeader() {
	p.buf.Write([]byte{
		/*0x00*/ 0x4D, 0x5A, // signature
		/*0x02*/ 0x00, 0x00, // lastsize
		/*0x04*/ 0x00, 0x00, // nblocks
		/*0x06*/ 0x00, 0x00, // nreloc
		/*0x08*/ 0x00, 0x00, // hdrsize
		/*0x0A*/ 0x00, 0x00, // minalloc
		/*0x0C*/ 0x00, 0x00, // maxalloc
		/*0x0E*/ 0x00, 0x00, // ss
		/*0x10*/ 0x00, 0x00, // sp
		/*0x12*/ 0x00, 0x00, // checksum
		/*0x14*/ 0x00, 0x00, // ip
		/*0x16*/ 0x00, 0x00, // cs
		/*0x18*/ 0x00, 0x00, // relocpos
		/*0x1A*/ 0x00, 0x00, // noverlay
	})
	p.buf.Write(
		/*0x1C*/ make([]byte, 8), // reserved
	)
	p.buf.Write([]byte{
		/*0x20*/ 0x00, 0x00, // oem_id
		/*0x22*/ 0x00, 0x00, // oem_info
	})
	p.buf.Write(
		/*0x24*/ make([]byte, 20), // reserved
	)
	p.buf.Write([]byte{
		/*0x38*/ 0x40, 0x00, 0x00, 0x00, // e_lfanew
		/*0x40*/
	})
}

func (p *PE) writePEHeader() {
	p.buf.Write([]byte{
		/*0x00*/ 0x50, 0x45, 0x00, 0x00, // Signature
		/*0x04*/
	})
	p.writeFileHeader()
	p.writeOptionalHeader()
}

func (p *PE) writeFileHeader() {
	TimeDateStamp := int32(time.Now().Unix())

	p.buf.Write([]byte{
		/*0x00*/ 0x64, 0x86, // Machine
	})
	/*0x02*/ binary.Write(&p.buf, binary.LittleEndian, int16(len(p.sections)))
	/*0x04*/ binary.Write(&p.buf, binary.LittleEndian, TimeDateStamp)
	p.buf.Write([]byte{
		/*0x08*/ 0x00, 0x00, 0x00, 0x00, // PointerToSymbolTable
		/*0x0C*/ 0x00, 0x00, 0x00, 0x00, // NumberOfSymbols
	})
	p.buf.Write([]byte{
		/*0x10*/ 0xF0, 0x00, // SizeOfOptionalHeader
		/*0x12*/ 0x22, 0x00, // Characteristics
		/*0x14*/
	})
}

func (p *PE) writeOptionalHeader() {
	p.buf.Write([]byte{
		/*0x00*/ 0x0B, 0x02, // Magic
		/*0x02*/ p.LinkerVersion[0], // MajorLinkerVersion
		/*0x03*/ p.LinkerVersion[1], // MinorLinkerVersion
	})
	/*0x04*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfCode())
	/*0x08*/ binary.Write(&p.buf, binary.LittleEndian, p.SizeOfInitializedData)
	/*0x0C*/ binary.Write(&p.buf, binary.LittleEndian, p.SizeOfUninitializedData)
	/*0x10*/ binary.Write(&p.buf, binary.LittleEndian, p.addressOfEntryPoint)
	/*0x14*/ binary.Write(&p.buf, binary.LittleEndian, p.baseOfCode)
	/*0x18*/ binary.Write(&p.buf, binary.LittleEndian, p.imageBase)
	/*0x20*/ binary.Write(&p.buf, binary.LittleEndian, p.sectionAlignment)
	/*0x24*/ binary.Write(&p.buf, binary.LittleEndian, p.fileAlignment)
	p.buf.Write([]byte{
		/*0x28*/ 0x06, 0x00, // MajorOperatingSystemVersion
		/*0x2A*/ 0x00, 0x00, // MinorOperatingSystemVersion
		/*0x2C*/ 0x00, 0x00, // MajorImageVersion
		/*0x2E*/ 0x00, 0x00, // MinorImageVersion
		/*0x30*/ 0x06, 0x00, // MajorSubsystemVersion
		/*0x32*/ 0x00, 0x00, // MinorSubsystemVersion
		/*0x34*/ 0x00, 0x00, 0x00, 0x00, // Win32VersionValue
	})
	/*0x38*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfImage())
	/*0x3C*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfHeaders())
	/*0x40*/ binary.Write(&p.buf, binary.LittleEndian, p.CheckSum)
	p.buf.Write([]byte{
		/*0x44*/ 0x03, 0x00, // Subsystem
		/*0x46*/ 0x00, 0x00, // DllCharacteristics
	})
	/*0x48*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfStackReserve)
	/*0x50*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfStackCommit)
	/*0x58*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfHeapReserve)
	/*0x60*/ binary.Write(&p.buf, binary.LittleEndian, p.sizeOfHeapCommit)
	p.buf.Write([]byte{
		/*0x68*/ 0x00, 0x00, 0x00, 0x00, // LoaderFlags
		/*0x6C*/ 0x10, 0x00, 0x00, 0x00, // NumberOfRvaAndSizes
		/*0x70*/
	})
	p.writeDataDirectories()
}

func (p *PE) writeDataDirectories() {
	type dataDirectory struct {
		Address int32
		Size    int32
	}
	exportTable := dataDirectory{}
	importTable := dataDirectory{
		Address: p.sections[1].header.VirtualAddress,
		Size:    p.sections[1].header.VirtualSize,
	}
	resourceTable := dataDirectory{}
	exceptionTable := dataDirectory{}
	certificateTable := dataDirectory{}
	baseRelocationTable := dataDirectory{}
	debug := dataDirectory{}
	architecture := dataDirectory{}
	globalPtr := dataDirectory{}
	TLSTable := dataDirectory{}
	loadConfigTable := dataDirectory{}
	boundImport := dataDirectory{}
	IAT := dataDirectory{}
	delayImportDescriptor := dataDirectory{}
	CLRRuntimeHeader := dataDirectory{}
	reserved := dataDirectory{}
	/*0x70*/ binary.Write(&p.buf, binary.LittleEndian, exportTable)
	/*0x78*/ binary.Write(&p.buf, binary.LittleEndian, importTable)
	/*0x80*/ binary.Write(&p.buf, binary.LittleEndian, resourceTable)
	/*0x88*/ binary.Write(&p.buf, binary.LittleEndian, exceptionTable)
	/*0x90*/ binary.Write(&p.buf, binary.LittleEndian, certificateTable)
	/*0x98*/ binary.Write(&p.buf, binary.LittleEndian, baseRelocationTable)
	/*0xA0*/ binary.Write(&p.buf, binary.LittleEndian, debug)
	/*0xA8*/ binary.Write(&p.buf, binary.LittleEndian, architecture)
	/*0xB0*/ binary.Write(&p.buf, binary.LittleEndian, globalPtr)
	/*0xB8*/ binary.Write(&p.buf, binary.LittleEndian, TLSTable)
	/*0xC0*/ binary.Write(&p.buf, binary.LittleEndian, loadConfigTable)
	/*0xC8*/ binary.Write(&p.buf, binary.LittleEndian, boundImport)
	/*0xD0*/ binary.Write(&p.buf, binary.LittleEndian, IAT)
	/*0xD8*/ binary.Write(&p.buf, binary.LittleEndian, delayImportDescriptor)
	/*0xE0*/ binary.Write(&p.buf, binary.LittleEndian, CLRRuntimeHeader)
	/*0xE8*/ binary.Write(&p.buf, binary.LittleEndian, reserved)
	/*0xF0*/
}

func (p *PE) writeSectionTable() {
	for _, s := range p.sections {
		s.writeHeader()
	}
}

func (p *PE) writeSectionData() {
	for _, s := range p.sections {
		s.writeData()
	}
}

func (p *PE) WriteFile(path string) {
	for _, s := range p.sections {
		s.update()
	}

	p.writeDOSHeader()
	p.writePEHeader()
	p.writeSectionTable()

	// pad headers
	p.buf.Write(make([]byte, p.sizeOfHeaders()-p.actualSizeOfHeaders()))

	p.writeSectionData()

	err := ioutil.WriteFile(path, p.buf.Bytes(), 0)
	if err != nil {
		panic(err)
	}
}

func (p *PE) SetCode(data []byte) {
	p.sections[0].SetData(data)
}
