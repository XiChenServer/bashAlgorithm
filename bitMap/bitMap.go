package bitMap

type BitMap struct {
	bitmap []byte
	length int
}

func NewBitMap() *BitMap {
	return &BitMap{}
}

// Add 添加数据到位图中
func (bitmap *BitMap) Add(num int) {
	// 如果输入为零，直接返回，不做任何处理
	if num == 0 {
		return
	}

	// 计算目标位置的字节索引和偏移量
	// byteIndex 用于确定 num 应该在位图中的哪个字节
	// offsetIndex 用于确定 num 在该字节中的哪个位
	byteIndex := num / 8
	offsetIndex := num % 8

	// 如果计算出的字节索引超出当前位图长度，则扩展位图容量
	// 如果 offsetIndex > 0，说明该数值映射到位图中某个字节的部分位，需要完整地添加一个字节
	// 如果 offsetIndex == 0，说明该数值正好是一个字节的起始位置，不需要额外添加字节
	if byteIndex+1 >= bitmap.length {
		if offsetIndex > 0 {
			// 扩展位图大小，以容纳新的字节
			bitmap.bitmap = append(bitmap.bitmap, make([]byte, byteIndex-len(bitmap.bitmap)+1)...)
		} else {
			// 如果正好对齐在字节边界上，不需要额外的字节
			bitmap.bitmap = append(bitmap.bitmap, make([]byte, byteIndex-len(bitmap.bitmap))...)
		}
	}

	// 根据 offsetIndex 的值，设置对应位置的位
	if offsetIndex == 0 {
		// 如果 offsetIndex 为 0，说明这个数值正好位于上一个字节的最后一位
		// 因此，我们需要将上一个字节的最后一位 置为 1
		if byteIndex >= 1 {
			bitmap.bitmap[byteIndex-1] |= 1
		}
	} else {
		// 否则，使用左移操作将对应的位设置为 1
		// 这里使用 8-offsetIndex 来确定要设置的位的位置
		bitmap.bitmap[byteIndex] |= 1 << (8 - offsetIndex)
	}

	// 更新位图的总长度，以反映扩展后的大小
	bitmap.length = len(bitmap.bitmap)
}

// Del 从位图中删除指定的数值
func (bitmap *BitMap) Del(num int) {
	// 计算目标位置的字节索引和偏移量
	// byteIndex 用于确定 num 在位图中的哪个字节
	// offsetIndex 用于确定 num 在该字节中的哪个位
	byteIndex := num / 8
	offsetIndex := num % 8

	// 检查该数值是否存在于位图中，如果存在则继续
	if bitmap.Exist(num) {
		// 如果 offsetIndex 为 0，说明这个数值正好位于上一个字节的最后一位
		if offsetIndex == 0 {
			// 确保当前字节索引大于或等于 1，避免越界
			if byteIndex >= 1 {
				// 检查上一个字节的最后一位是否为 1（即该数值是否存在）
				if (bitmap.bitmap[byteIndex-1] & 1) != 0 {
					// 使用按位与操作将上一个字节的最后一位清零，删除该数值
					bitmap.bitmap[byteIndex-1] &= ^uint8(1)
				}
			}
			return // 提前返回，因为该数值已经被删除
		}

		// 如果 offsetIndex 不为 0，则计算该数值在位图中的位置
		if bitmap.bitmap[byteIndex]&(1<<(8-offsetIndex)) != 0 {
			// 使用按位与操作将对应位置的位清零，删除该数值
			bitmap.bitmap[byteIndex] &= ^uint8(1 << (8 - offsetIndex))
		}
	}

	return // 函数结束
}

// Exist 检查指定的数值是否存在于位图中
func (bitmap *BitMap) Exist(num int) bool {
	// 如果输入为零，直接返回 false，因为零不在位图中表示
	if num == 0 {
		return false
	}

	// 计算目标位置的字节索引和偏移量
	byteIndex := num / 8
	offsetIndex := num % 8

	// 检查索引是否超出位图的当前长度
	// 如果 byteIndex 大于位图长度，或者刚好位于边界且有偏移量，则说明数值不存在
	if (byteIndex > bitmap.length) || (byteIndex == bitmap.length && offsetIndex > 0) {
		return false
	}

	// 如果 offsetIndex 为 0，说明数值可能位于上一个字节的最后一位
	if offsetIndex == 0 {
		// 检查上一个字节的最后一位是否为 1（即该数值是否存在）
		if byteIndex >= 1 {
			return bitmap.bitmap[byteIndex-1]&1 != 0
		}
	}

	// 对应的位是否为 1，确定该数值是否存在
	return bitmap.bitmap[byteIndex]&(1<<(8-offsetIndex)) != 0
}

// Len 返回位图的当前长度（字节数）
func (bitmap *BitMap) Len() int {
	return bitmap.length
}

func (bitmap *BitMap) ToString() string {
	bits := ""
	for _, b := range bitmap.bitmap {
		for i := 7; i >= 0; i-- {
			bits += string('0' + (b >> i & 1))
		}
	}
	return bits
}
