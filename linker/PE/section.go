package pe

import "encoding/binary"

type SectionHeader struct {
	Name                 [8]byte /*0x00*/
	VirtualSize          int32   /*0x08*/
	VirtualAddress       int32   /*0x0C*/
	SizeOfRawData        int32   /*0x10*/
	PointerToRawData     int32   /*0x14*/
	PointerToRelocations int32   /*0x18*/
	PointerToLinenumbers int32   /*0x1C*/
	NumberOfRelocations  int16   /*0x20*/
	NumberOfLinenumbers  int16   /*0x22*/
	Characteristics      int32   /*0x24*/
	/*0x28*/
}

type Section struct {
	header SectionHeader
	pe     *PE
	data   []byte
}

func newSection(pe *PE, name string) *Section {
	if len(name) > 8 {
		panic("name must be at most 8 characters")
	}
	s := &Section{
		pe: pe,
		header: SectionHeader{
			Characteristics: 0x60000020,
			VirtualAddress:  pe.baseOfCode,
		},
	}
	copy(s.header.Name[:], name)
	return s
}

func (s *Section) SetData(d []byte) {
	s.data = d
	s.header.SizeOfRawData = int32(multipleOf(len(s.data), int(s.pe.fileAlignment)))
	s.header.VirtualSize = int32(len(s.data))
}

func (s *Section) writeHeader() {
	s.calcRawAddress()

	binary.Write(&s.pe.buf, binary.LittleEndian, s.header)
}

func (s *Section) writeData() {
	padding := int(s.header.SizeOfRawData) - len(s.data)

	binary.Write(&s.pe.buf, binary.LittleEndian, s.data)
	s.pe.buf.Write(make([]byte, padding))
}

func (s *Section) calcRawAddress() {
	offset := s.pe.sizeOfHeaders()

	for _, currS := range s.pe.sections {
		if currS == s {
			break
		}
		offset += currS.header.SizeOfRawData
	}
	s.header.PointerToRawData = offset
}
