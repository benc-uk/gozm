package zmachine

import "gozm/internal/decode"

const (
	OBJECT_ENTRY_SIZE    = 9
	OBJECT_PARENT_INDEX  = 4
	OBJECT_SIBLING_INDEX = 5
	OBJECT_CHILD_INDEX   = 6
	NULL_OBJECT_INDEX    = 0
)

type zObject struct {
	num     byte
	attr    [32]bool
	parent  byte
	sibling byte
	child   byte
	//properties map[byte][]byte
}

// Parses the object table and initializes the objects in the machine
func (m *Machine) initObjects() {
	// Parse property defaults table, don't think it's used much
	for i := 0; i < 31; i++ {
		propAddr := m.objectsAddr + uint16(i*2)
		val := decode.GetWord(m.mem, propAddr)
		m.propDefaults[i] = val
	}

	objTableAddr := m.objectsAddr + 62 // skip 31 properties * 2 bytes each
	objCount := 0
	propTableAddr := uint16(0xffff) // will be set when first object is read
	for {
		objEntryAddr := objTableAddr + uint16(objCount*OBJECT_ENTRY_SIZE)

		// Stop if we've reached the property table
		if objEntryAddr >= propTableAddr {
			break
		}

		// Parse object entry
		objCount++
		obj := &zObject{
			num: byte(objCount),
		}

		// Read attributes bit at a time (4 bytes, 32 bits total)
		for i := 0; i < 32; i++ {
			byteIndex := i / 8
			bitIndex := 7 - (i % 8)
			attrByte := m.mem[objEntryAddr+uint16(byteIndex)]
			obj.attr[i] = (attrByte&(1<<bitIndex) != 0)
		}

		// Read parent, sibling, child (1 byte each)
		obj.parent = m.mem[objEntryAddr+4]
		obj.sibling = m.mem[objEntryAddr+5]
		obj.child = m.mem[objEntryAddr+6]

		// Read properties pointer (2 bytes)
		propAddr := decode.GetWord(m.mem, objEntryAddr+7)
		if objCount == 1 {
			propTableAddr = propAddr
		}

		m.trace("Parsed object %d: parent=%d, sibling=%d, child=%d, props=%04x, attr=%v\n",
			obj.num, obj.parent, obj.sibling, obj.child, propAddr, obj.attr)

		m.objects[obj.num] = obj

		// Sanity check to prevent infinite loop
		if objCount > 255 {
			break
		}
	}
}
