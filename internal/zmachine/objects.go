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
	Num     byte               `json:"num"`
	Desc    string             `json:"desc"`
	Attrs   [32]bool           `json:"attrs"`
	Parent  byte               `json:"parent"`
	Sibling byte               `json:"sibling"`
	Child   byte               `json:"child"`
	Props   []*property        `json:"props"`
	PropMap map[byte]*property `json:"prop_map"`
}

type property struct {
	Num  byte   `json:"num"`
	Size byte   `json:"size"`
	Data []byte `json:"data"`
	Addr uint16 `json:"addr"` // address in memory where this property is stored
}

// Parses the object table and initializes the objects in the machine
// This is called once during machine initialization
func (m *Machine) initObjects() {
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
			Num:     byte(objCount),
			Props:   make([]*property, 0), // max 64 properties
			PropMap: make(map[byte]*property),
		}

		// Read attributes bit at a time (4 bytes, 32 bits total)
		for i := 0; i < 4; i++ {
			attrByte := m.mem[objEntryAddr+uint16(i)]
			for bit := 0; bit < 8; bit++ {
				attrIndex := i*8 + bit
				obj.Attrs[attrIndex] = (attrByte & (1 << (7 - bit))) != 0
			}
		}

		// Read parent, sibling, child (1 byte each)
		obj.Parent = m.mem[objEntryAddr+4]
		obj.Sibling = m.mem[objEntryAddr+5]
		obj.Child = m.mem[objEntryAddr+6]

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
		obj.Desc = decode.String(descWords, m.abbr) // decode string literal
		currPropAddr += 1 + uint16(descSize)*2

		// Now read properties until we hit a 0 size byte in the header
		for {
			propHeader := m.mem[currPropAddr]

			// End of properties for this object
			if propHeader == 0 {
				break
			}

			// Property number is in the lower 5 bits, size in upper 3 bits
			propNum, propSize := decode.PropSizeNumber(propHeader) // Alternative way to get propNum and propSize
			propData := make([]byte, propSize)
			for i := byte(0); i < propSize; i++ {
				propData[i] = m.mem[currPropAddr+1+uint16(i)]
			}

			prop := &property{
				Num:  propNum,
				Size: propSize,
				Data: propData,
				Addr: currPropAddr + 1, // Point to data not header
			}
			obj.Props = append(obj.Props, prop)
			obj.PropMap[propNum] = prop

			currPropAddr += 1 + uint16(propSize)
		}

		// Debug output
		m.trace("Parsed object '%s' (%d): parent=%d, sibling=%d, child=%d, props=%04x, attr=%v\n",
			obj.Desc, obj.Num, obj.Parent, obj.Sibling, obj.Child, propAddr, obj.Attrs)
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

	return o.Attrs[attrNum]
}

func (o *zObject) setAttribute(attrNum byte, value bool) {
	if attrNum > 31 {
		return
	}

	o.Attrs[attrNum] = value
}

func (o *zObject) removeObjectFromParent(m *Machine) {
	if o.Parent == NULL_OBJECT {
		return // No parent to remove from
	}

	parentObj := m.getObject(o.Parent)
	if parentObj == nil {
		return // Parent object not found
	}

	// If this object is the first child of the parent
	if parentObj.Child == o.Num {
		parentObj.Child = o.Sibling
	} else {
		// If this object is not the first child, find the previous sibling
		siblingNum := parentObj.Child
		for siblingNum != NULL_OBJECT {
			siblingObj := m.getObject(siblingNum)
			if siblingObj == nil {
				break // Should not happen
			}
			if siblingObj.Sibling == o.Num {
				// Found the previous sibling, update its sibling pointer
				siblingObj.Sibling = o.Sibling
				break
			}
			siblingNum = siblingObj.Sibling
		}
	}

	// Clear this object's parent and sibling pointers
	o.Parent = NULL_OBJECT
	o.Sibling = NULL_OBJECT
}

func (o *zObject) insertIntoParent(m *Machine, newParentNum byte) {
	// First remove from current parent if any
	o.removeObjectFromParent(m)

	// Set new parent
	o.Parent = newParentNum
	newParentObj := m.getObject(newParentNum)
	if newParentObj == nil {
		return // New parent object not found
	}

	// Insert as first child of new parent
	o.Sibling = newParentObj.Child
	newParentObj.Child = o.Num
}

func (o *zObject) getPropertyValue(propNum byte, defaults []uint16) uint16 {
	prop, exists := o.PropMap[propNum]
	if !exists {
		// Return default value if property not found
		if int(propNum) < len(defaults) {
			return defaults[propNum-1]
		}
		return 0
	}

	// Return property value as uint16 (assuming properties are at most 2 bytes)
	if prop.Size == 1 {
		return uint16(prop.Data[0])
	} else if prop.Size == 2 {
		return decode.GetWord(prop.Data, 0)
	}

	return 0 // Unsupported property size
}

func (o *zObject) setPropertyValue(propNum byte, value uint16) {
	prop, exists := o.PropMap[propNum]
	if !exists {
		// Property does not exist, cannot set
		return
	}

	// Set property value (assuming properties are at most 2 bytes)
	if prop.Size == 1 {
		prop.Data[0] = byte(value & 0xFF)
	} else if prop.Size == 2 {
		prop.Data[0] = byte((value >> 8) & 0xFF)
		prop.Data[1] = byte(value & 0xFF)
	}
}
