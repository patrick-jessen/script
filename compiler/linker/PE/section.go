package pe

import "encoding/binary"

type sectionHeader struct {
	Name                 [8]byte /*0x00*/
	VirtualSize          int32   /*0x08*/
	VirtualAddress       int32   /*0x0C*/
	SizeOfRawData        int32   /*0x10*/
	PointerToRawData     int32   /*0x14*/
	PointerToRelocations int32   /*0x18*/
	PointerToLinenumbers int32   /*0x1C*/
	NumberOfRelocations  int16   /*0x20*/
	NumberOfLinenumbers  int16   /*0x22*/
	Characteristics      uint32  /*0x24*/
	/*0x28*/
}

type section struct {
	header sectionHeader
	pe     *PE
	data   []byte
}

func newSection(pe *PE, name string, characteristics uint32) *section {
	if len(name) > 8 {
		panic("name must be at most 8 characters")
	}
	s := &section{
		pe: pe,
		header: sectionHeader{
			Characteristics: characteristics,
			SizeOfRawData:   pe.fileAlignment,
		},
	}
	copy(s.header.Name[:], name)
	return s
}

func (s *section) SetData(d []byte) {
	s.data = d
	s.header.SizeOfRawData = multipleOf(int32(len(s.data)), s.pe.fileAlignment)
	s.header.VirtualSize = int32(len(s.data))
}

func (s *section) update() {
	s.calcRawAddress()
	s.calcVirtualAddress()
}

func (s *section) writeHeader() {
	binary.Write(&s.pe.buf, binary.LittleEndian, s.header)
}

func (s *section) writeData() {
	padding := int(s.header.SizeOfRawData) - len(s.data)

	binary.Write(&s.pe.buf, binary.LittleEndian, s.data)
	s.pe.buf.Write(make([]byte, padding))
}

func (s *section) calcVirtualAddress() {
	va := s.pe.baseOfCode
	for _, currS := range s.pe.sections {
		if currS == s {
			break
		}
		va += multipleOf(currS.header.VirtualSize, s.pe.sectionAlignment)
	}
	s.header.VirtualAddress = va
}

func (s *section) calcRawAddress() {
	offset := s.pe.sizeOfHeaders()

	for _, currS := range s.pe.sections {
		if currS == s {
			break
		}
		offset += currS.header.SizeOfRawData
	}
	s.header.PointerToRawData = offset
}
