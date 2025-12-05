// =======================================================================
// Package: zmachine - Core Z-machine interpreter
// objects.go - Object table parsing and object management
//
// Copyright (c) 2025 Ben Coleman. Licensed under the MIT License
// =======================================================================

package zmachine

import (
	"fmt"

	"github.com/benc-uk/gozm/internal/decode"
)

const (
	NULL_OBJECT = 0
)

type zObject struct {
	num         byte
	description string
	attr        [32]bool
	parent      byte
	sibling     byte
	child       byte
	properties  []*property
	propMap     map[byte]*property
}

type property struct {
	num  byte
	size byte
	data []byte
	addr uint16 // address in memory where this property is stored
}

// Parses the object table and initializes the objects in the machine
// This is called once during machine initialization
func (m *Machine) initObjects() {
	fmt.Printf("Initializing object table at %04x\n", m.objectsAddr)
	if len(m.objects) > 0 || m.objectsAddr == 0 {
		return // Already initialized or we have no object table
	}

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
		objEntryAddr := objTableAddr + uint16(objCount*9)

		// Stop if we've reached the property table
		if objEntryAddr >= propTableAddr {
			break
		}

		// Parse object entry
		objCount++
		obj := &zObject{
			num:        byte(objCount),
			properties: make([]*property, 0), // max 64 properties
			propMap:    make(map[byte]*property),
		}

		// Read attributes bit at a time (4 bytes, 32 bits total)
		for i := 0; i < 4; i++ {
			attrByte := m.mem[objEntryAddr+uint16(i)]
			for bit := 0; bit < 8; bit++ {
				attrIndex := i*8 + bit
				obj.attr[attrIndex] = (attrByte & (1 << (7 - bit))) != 0
			}
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

		// Parse properties table for this object

		// Property table header starts with a description: [ size byte | string literal (size * 2 bytes) ]
		currPropAddr := propAddr
		descSize := m.mem[currPropAddr]
		descWords := make([]uint16, descSize)
		for i := uint16(0); i < uint16(descSize); i++ {
			word := decode.GetWord(m.mem, currPropAddr+1+i*2)
			descWords[i] = word
		}
		obj.description = decode.String(descWords, m.abbr) // decode string literal
		currPropAddr += 1 + uint16(descSize)*2

		// Now read properties until we hit a 0 size byte in the header
		for {
			propHeader := m.mem[currPropAddr]
			if propHeader == 0 {
				// End of properties for this object
				currPropAddr++
				break
			}

			// Property number is in the lower 5 bits, size in upper 3 bits
			propNum, propSize := decode.PropSizeNumber(propHeader) // Alternative way to get propNum and propSize
			propData := make([]byte, propSize)
			for i := byte(0); i < propSize; i++ {
				propData[i] = m.mem[currPropAddr+1+uint16(i)]
			}

			prop := &property{
				num:  propNum,
				size: propSize,
				data: propData,
				addr: currPropAddr + 1, // Point to data not header
			}
			obj.properties = append(obj.properties, prop)
			obj.propMap[propNum] = prop

			currPropAddr += 1 + uint16(propSize)
		}

		// Debug output
		m.trace("Parsed object '%s' (%d): parent=%d, sibling=%d, child=%d, props=%04x, attr=%v\n",
			obj.description, obj.num, obj.parent, obj.sibling, obj.child, propAddr, obj.attr)
		m.trace("%s\n\n", obj.propDebugDump())

		// We assume objects are in order by number
		m.objects = append(m.objects, obj)

		// Sanity check to prevent infinite loop
		if objCount > 255 {
			break
		}
	}
}

// Helper to get an object by its number
func (m *Machine) getObject(objNum byte) *zObject {
	if objNum == NULL_OBJECT || int(objNum) > len(m.objects) {
		panic(fmt.Sprintf("FATAL: Attempt to access object %d", objNum))
	}

	return m.objects[objNum-1]
}

func (o *zObject) hasAttribute(attrNum byte) bool {
	if attrNum > 31 {
		return false
	}

	return o.attr[attrNum]
}

func (o *zObject) setAttribute(attrNum byte, value bool) {
	if attrNum > 31 {
		return
	}

	o.attr[attrNum] = value
}

func (o *zObject) removeObjectFromParent(m *Machine) {
	if o.parent == NULL_OBJECT {
		return // No parent to remove from
	}

	parentObj := m.getObject(o.parent)
	if parentObj == nil {
		return // Parent object not found
	}

	// If this object is the first child of the parent
	if parentObj.child == o.num {
		parentObj.child = o.sibling
	} else {
		// If this object is not the first child, find the previous sibling
		siblingNum := parentObj.child
		for siblingNum != NULL_OBJECT {
			siblingObj := m.getObject(siblingNum)
			if siblingObj == nil {
				break // Should not happen
			}
			if siblingObj.sibling == o.num {
				// Found the previous sibling, update its sibling pointer
				siblingObj.sibling = o.sibling
				break
			}
			siblingNum = siblingObj.sibling
		}
	}

	// Clear this object's parent and sibling pointers
	o.parent = NULL_OBJECT
	o.sibling = NULL_OBJECT
}

func (o *zObject) insertIntoParent(m *Machine, newParentNum byte) {
	// First remove from current parent if any
	o.removeObjectFromParent(m)

	// Set new parent
	o.parent = newParentNum
	newParentObj := m.getObject(newParentNum)
	if newParentObj == nil {
		return // New parent object not found
	}

	// Insert as first child of new parent
	o.sibling = newParentObj.child
	newParentObj.child = o.num
}

func (o *zObject) getPropertyValue(propNum byte, defaults []uint16) uint16 {
	prop, exists := o.propMap[propNum]
	if !exists {
		// Return default value if property not found
		if int(propNum) < len(defaults) {
			return defaults[propNum-1]
		}
		return 0
	}

	// Return property value as uint16 (assuming properties are at most 2 bytes)
	if prop.size == 1 {
		return uint16(prop.data[0])
	} else if prop.size == 2 {
		return decode.GetWord(prop.data, 0)
	}

	return 0 // Unsupported property size
}

func (o *zObject) setPropertyValue(propNum byte, value uint16) {
	prop, exists := o.propMap[propNum]
	if !exists {
		// Property does not exist, cannot set
		return
	}

	// Set property value (assuming properties are at most 2 bytes)
	if prop.size == 1 {
		prop.data[0] = byte(value & 0xFF)
	} else if prop.size == 2 {
		prop.data[0] = byte((value >> 8) & 0xFF)
		prop.data[1] = byte(value & 0xFF)
	}
}
