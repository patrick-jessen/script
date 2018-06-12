package pe

import (
	"bytes"
	"encoding/binary"
)

type Importer struct {
	descriptors   []importDescriptor
	lookupTables  []importTable
	names         []importName
	addressTables []importTable
	dllNames      []string

	dlls       map[string]int
	nameOffset int
}

func (i *Importer) update(offset int) {
	for idx := 0; idx < len(i.descriptors); idx++ {
		i.descriptors[idx].update(i, idx, offset)
	}
	for idx := 0; idx < len(i.lookupTables); idx++ {
		i.lookupTables[idx].update(i)
	}
	for idx := 0; idx < len(i.addressTables); idx++ {
		i.addressTables[idx].update(i)
	}
}

func (i *Importer) Import(symbol string, dll string) {
	var lookupTable *importTable
	var addressTable *importTable

	descIdx, ok := i.dlls[dll]
	if !ok {
		descIdx = len(i.descriptors)
		i.dlls[dll] = descIdx
		i.descriptors = append(i.descriptors, importDescriptor{})
		i.lookupTables = append(i.lookupTables, importTable{})
		i.addressTables = append(i.addressTables, importTable{})
		i.dllNames = append(i.dllNames, dll)
	}
	lookupTable = &i.lookupTables[descIdx]
	addressTable = &i.addressTables[descIdx]

	nameIdx := int64(len(i.names))
	i.names = append(i.names, importName{name: symbol})

	lookupTable.entries = append(lookupTable.entries, nameIdx)
	addressTable.entries = append(addressTable.entries, nameIdx)
}

func (i *Importer) Write(buf *bytes.Buffer, offset int) {
	i.update(offset)

	for _, descr := range i.descriptors {
		binary.Write(buf, binary.LittleEndian, descr)
	}
	binary.Write(buf, binary.LittleEndian, importDescriptor{}) // Last entry must be 0

	for _, lookup := range i.lookupTables {
		lookup.write(buf)
	}

	for _, name := range i.names {
		name.write(buf)
	}

	for _, address := range i.addressTables {
		address.write(buf)
	}

	for _, s := range i.dllNames {
		buf.Write([]byte(s))
		buf.WriteByte(0x00)
	}
}

func newImporter() *Importer {
	return &Importer{
		dlls: make(map[string]int),
	}
}

type importDescriptor struct {
	ImportLookupTableRVA  int32 // RVA to ImportLookupTable
	TimeDateStamp         int32 // reserved
	ForwarderChain        int32 // Not used
	NameRVA               int32 // RVA to null-terminated DLL string
	ImportAddressTableRVA int32
}

func (id *importDescriptor) update(i *Importer, idx int, offset int) {
	offset += (len(i.descriptors) + 1) * 20 // each importDescriptor is 20 bytes + one empty

	var lookUpOffset int
	for l := 0; l < len(i.lookupTables); l++ {
		if l == idx {
			lookUpOffset = offset
		}
		offset += (len(i.lookupTables[l].entries) + 1) * 8 // each importTable is 8 bytes + one empty
	}
	id.ImportLookupTableRVA = int32(lookUpOffset)
	i.nameOffset = offset

	for n := 0; n < len(i.names); n++ {
		offset += 2 + len(i.names[n].name) + 1 // hint is 2, \0 is 1
		if offset%2 == 1 {
			offset++
		}
	}

	var addressOffset int
	for l := 0; l < len(i.addressTables); l++ {
		if l == idx {
			addressOffset = offset
		}
		offset += (len(i.addressTables[l].entries) + 1) * 8 // each importTable is 8 bytes + one empty
	}
	id.ImportAddressTableRVA = int32(addressOffset)

	for n := 0; n < idx; n++ {
		offset += len(i.dllNames[n]) + 1 // \0 is 1
	}
	id.NameRVA = int32(offset)
}

type importTable struct {
	// only bits [0-30] are set
	// they are RVAs to a name table entry
	// last entry is 0
	entries []int64
}

func (ilt *importTable) update(i *Importer) {
	for idx := 0; idx < len(ilt.entries); idx++ {
		nameIdx := int(ilt.entries[idx])
		offset := i.nameOffset
		for ni := 0; ni < nameIdx; ni++ {
			offset += 2 + len(i.names[ni].name) + 1 // hint is 2, \0 is 1
			if offset%2 == 1 {
				offset++
			}
		}

		ilt.entries[idx] = int64(offset)
	}
}

func (ilt *importTable) write(buf *bytes.Buffer) {
	for _, entry := range ilt.entries {
		binary.Write(buf, binary.LittleEndian, entry)
	}
	binary.Write(buf, binary.LittleEndian, int64(0)) // Last entry must be 0
}

type importName struct {
	hint int16
	name string
}

func (in *importName) write(buf *bytes.Buffer) {
	name := []byte(in.name)

	binary.Write(buf, binary.LittleEndian, in.hint)
	buf.Write(name)
	buf.WriteByte(0x00)

	if len(name)%2 == 0 {
		// Pad to align the next entry on an even boundary
		buf.WriteByte(0x00)
	}
}
