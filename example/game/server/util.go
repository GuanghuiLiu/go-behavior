package server

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(n int) int {
	return rand.Intn(n)
}
func RandUint64() uint64 {
	return rand.Uint64()
}
func RandUint8(n uint8) uint8 {
	return uint8(rand.Intn(int(n)))
}
func RandUint64n(n int) uint64 {
	return uint64(rand.Intn(n))
}

func isAtString(s string, sa []string) bool {
	for _, v := range sa {
		if s == v {
			return true
		}
	}
	return false
}
func isAtInt(s int, sa []int) bool {
	for _, v := range sa {
		if s == v {
			return true
		}
	}
	return false
}
func maxSame(i, j int) int {
	if i > j {
		return i
	}
	return j
}

// 不做任意数量洗牌，是为了公平，游戏中不洗牌
func shuffleCards(cards *[52]uint8) {
	for i, _ := range cards {
		target := RandUint8(52)
		cards[i], cards[target] = cards[target], cards[i]
	}
}

func judgeCardLv(cards []uint8) (uint8, []uint8) {
	// fmt.Println("ccc10111", cards)
	if len(cards) != 7 {
		return CardTypeErr, nil
	}
	tongHua, tCards := isTongHua(cards)
	// fmt.Println("ccc101112", tongHua, tCards)
	if tongHua {
		sHunZi, max := isSHunZi(tCards)
		if sHunZi {
			if max[len(max)-1] == 14 {
				return CardTypeHJTHS, max
			}
			return CardTypeTHS, max
		}
		return CardTypeTH, max
	}
	// fmt.Println("ccc1011123")
	return commonCharge(cards)
}

func isTongHua(cards []uint8) (bool, []uint8) {

	h, ho, f, m := []uint8{}, []uint8{}, []uint8{}, []uint8{}
	for _, card := range cards {
		switch {
		case card < CardTypeGapLV1:
			h = append(h, card)
		case card > CardTypeGapLV1 && card < CardTypeGapLV2:
			ho = append(ho, card)
		case card > CardTypeGapLV2 && card < CardTypeGapLV3:
			f = append(f, card)
		case card > CardTypeGapLV3:
			m = append(m, card)
		}
	}
	switch {
	case len(h) > 4:
		return true, h
	case len(ho) > 4:
		return true, ho
	case len(f) > 4:
		return true, f
	case len(m) > 4:
		return true, m
	}
	return false, cards
}

func isSHunZi(cards []uint8) (bool, []uint8) {
	// fmt.Println("0ccc10111234,before", cards)
	sortCardH(&cards)
	// fmt.Println("0ccc10111234,after", cards)
	list := []uint8{cards[0]}
	for i := 1; i < len(cards); i++ {
		switch list[len(list)-1]%CardTypeGapLV1 - cards[i]%CardTypeGapLV1 {
		case 0:
		case 1:
			if len(list) == 5 {
				return true, append(list, cards[i])
			}
			// List最后一个是2，而且前面有3个，而且第一个是A，这个是最小的地顺
			if cards[i]%CardTypeGapLV1 == 2 && len(list) == 3 && list[0]%CardTypeGapLV1 == 14 {
				return true, append(list, cards[i], list[0])
			}
		default:
			list = []uint8{cards[i]}
		}
	}
	return false, cards
}

func commonCharge(cards []uint8) (uint8, []uint8) {
	// fmt.Println("ccc10111234")
	shuZi, res := isSHunZi(cards)
	// fmt.Println("ccc10111235", shuZi, res)
	if shuZi {
		return CardTypeSZ, res
	}
	sameArrayList := make([][]uint8, 1)
	sameArrayList[0] = make([]uint8, 1)

	sameArrayList[0][0] = cards[0]
	// sameArrayList[0]= append(sameArrayList[0],cards[0])
	// 折叠牌，成二维切片
	for i := 1; i < len(cards); i++ {
		// 二维切片最后一个元素
		if sameArrayList[len(sameArrayList)-1][len(sameArrayList[len(sameArrayList)-1])-1]%CardTypeGapLV1 == cards[i]%CardTypeGapLV1 {
			sameArrayList[len(sameArrayList)-1] = append(sameArrayList[len(sameArrayList)-1], cards[i])
		} else {
			sameArrayList = append(sameArrayList, []uint8{cards[i]})
		}
	}
	sortTwoDimensional(&sameArrayList)
	// fmt.Println("ccc101112356", sameArrayList)
	switch len(sameArrayList[0]) {
	case 4:
		return CardTypeSIT, append(sameArrayList[0], sameArrayList[1][0])
	case 3:
		if len(sameArrayList[1]) > 1 {
			return CardTypeHL, append(sameArrayList[0], sameArrayList[1][:2]...)
		}
		return CardTypeSANT, append(sameArrayList[0], sameArrayList[1][0], sameArrayList[2][0])
	case 2:
		if len(sameArrayList[1]) == 2 {
			return CardTypeLD, append(sameArrayList[0], sameArrayList[1][0], sameArrayList[1][1], sameArrayList[2][0])
		}
		return CardTypeYD, append(sameArrayList[0], sameArrayList[1][0], sameArrayList[2][0], sameArrayList[3][0])
	case 1:
		return CardTypeGP, append(sameArrayList[0], sameArrayList[1][0], sameArrayList[2][0], sameArrayList[3][0], sameArrayList[4][0])
	}

	return CardTypeErr, nil
}

func isST(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}
func isHL(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}

func isSZ(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}
func isSANT(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}
func isLD(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}
func isYD(cards []uint8) (bool, uint8) {
	return false, maxCard(cards)
}

func sortRole() {

}

func maxCard(cards []uint8) uint8 {
	max := cards[0]
	for _, card := range cards {
		if card%CardTypeGapLV1 > max {
			max = card % CardTypeGapLV1
		}
	}
	return max
}

func contain(n uint8, list []uint8) bool {
	for _, v := range list {
		if v == n {
			return true
		}
	}
	return false
}
func containCount(n uint8, list []uint8) (count uint8) {
	for _, v := range list {
		if v == n {
			count++
		}
	}
	return
}

func sortCardT(cards *[]uint8) {
	for i := 0; i < len(*cards)-1; i++ {
		for j := 0; j < len(*cards)-1-i; j++ {
			if (*cards)[j] > (*cards)[j+1] {
				(*cards)[j], (*cards)[j+1] = (*cards)[j+1], (*cards)[j]
			}
		}
	}
}
func sortCardH(cards *[]uint8) {
	for i := 0; i < len(*cards)-1; i++ {
		for j := 0; j < len(*cards)-1-i; j++ {
			if (*cards)[j]%CardTypeGapLV1 < (*cards)[j+1]%CardTypeGapLV1 {
				(*cards)[j], (*cards)[j+1] = (*cards)[j+1], (*cards)[j]
			}
		}
	}
}

func sortTwoDimensional(cards *[][]uint8) {
	for i := 0; i < len(*cards)-1; i++ {
		for j := 0; j < len(*cards)-1-i; j++ {
			if len((*cards)[j]) < len((*cards)[j+1]) {
				(*cards)[j], (*cards)[j+1] = (*cards)[j+1], (*cards)[j]
			}
		}
	}
}

func uint8toUint32List(old []uint8) (res []uint32) {
	for _, v := range old {
		res = append(res, uint32(v))
	}
	return
}

func comList(a, b []uint32) []uint32 {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return b
		}
	}
	return a
}
