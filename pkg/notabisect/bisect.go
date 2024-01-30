package notabisect

import (
	"math/big"
	"sort"

	"golang.org/x/exp/constraints"

	"github.com/arnokay/consistent-hashing/pkg/tools"
)

type allowedTypes interface {
	constraints.Float | constraints.Integer | string
}

func Bisect[T allowedTypes](slice []T, target T) int {
	return sort.Search(len(slice), func(i int) bool {
		return slice[i] >= target
	})
}

func BisectRight[T allowedTypes](slice []T, target T) int {
	return sort.Search(len(slice), func(i int) bool {
		return slice[i] > target
	})
}

func Insert[T allowedTypes](slice []T, target T) []T {
	i := Bisect(slice, target)
	return tools.Insert(slice, target, i)
}

func InsertRight[T allowedTypes](slice []T, target T) []T {
	i := BisectRight(slice, target)
	return tools.Insert(slice, target, i)
}

func BisectBigInt(slice []*big.Int, target *big.Int) int {
	return sort.Search(len(slice), func(i int) bool {
		if slice[i] == nil {
			return false
		}
		switch slice[i].Cmp(target) {
		case -1:
			return false
		case 0:
			return true
		case 1:
			return true
		default:
			panic("you not supposed to be here")
		}
	})
}

func BisectRightBigInt(slice []*big.Int, target *big.Int) int {
	return sort.Search(len(slice), func(i int) bool {
		switch slice[i].Cmp(target) {
		case -1:
			return false
		case 0:
			return false
		case 1:
			return true
		default:
			panic("you not supposed to be here")
		}
	})
}
