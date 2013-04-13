package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

var isPalindromeList = make(map[int64]bool, math.MaxInt16)
var isFairNSquareList = make(map[int64]bool, math.MaxInt16)

var FairNSquareList = make([]int64, 0, math.MaxInt16)

func main() {
	// Parse args for file name
	flag.Parse()
	fn := flag.Arg(0)

	// Input file
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	input, _ := ioutil.ReadAll(f)
	lines := strings.Split(string(input), "\n")

	// Output file
	of, err := os.OpenFile(strings.Replace(fn, ".in", ".out", -1), os.O_WRONLY|os.O_CREATE, 0666)
	defer of.Close()
	if err != nil {
		panic(err)
	}

	// Output function
	out := func(format string, ex ...interface{}) {
		s := fmt.Sprintf(format+"\n", ex...)
		fmt.Print(s)
		of.Write([]byte(s))
	}

	// Parse input file
	N_, _ := strconv.ParseInt(lines[0], 0, 0)
	N := int(N_)

	nFnS := 0

	var line []string
	var FnSmax, low, upp, cnt int64
	for i := 0; i < N; i++ {
		line = strings.Split(lines[i+1], " ")
		low, _ = strconv.ParseInt(line[0], 0, 0)
		upp, _ = strconv.ParseInt(line[1], 0, 0)
		if nFnS == 0 || FnSmax < upp {
			for FnSmax < upp {
				FnSmax = nextFairNSquare(FnSmax)
				FairNSquareList = append(FairNSquareList, FnSmax)
			}
		}
		cnt = 0
		for _, v := range FairNSquareList {
			if v >= low && v <= upp {
				cnt++
			} else if v > upp {
				break
			}
		}
		out("Case #%d: %d", i+1, cnt)
	}
}

func nDigits(dig int64) int64 {
	return int64(math.Floor(math.Log10(float64(dig)) + 1))
}

func digitAt(dig, pos int64) int64 {
	d := float64(dig)
	p := int(nDigits(dig) - pos)
	for i := 0; i < p; i++ {
		d = d / 10.0
	}
	return int64(d) % 10
}

func getDigits(dig int64) []int {
	d := float64(dig)
	n := int(nDigits(dig))
	digs := make([]int, n)
	for i := 0; i < n; i++ {
		digs[n-i-1] = int(int64(d) % 10)
		d = d / 10.0
	}
	return digs
}

func digsToInt(digs []int) (dig int64) {
	n := len(digs)
	for i := 0; i < n; i++ {
		dig += int64(float64(digs[i]) * math.Pow10(n-i-1))
	}
	return
}

func isPalindrome(dig int64) bool {
	chk, ok := isPalindromeList[dig]
	if ok {
		return chk
	}

	str := strconv.FormatInt(dig, 10)
	l := len(str)
	mid := int(math.Floor(float64(l) / 2.0))
	for i := 0; i < mid; i++ {
		if str[i] != str[l-i-1] {
			isPalindromeList[dig] = false
			return false
		}
	}
	isPalindromeList[dig] = true
	return true
}

func nextPalindrome(dig int64) int64 {
	digs := getDigits(dig)
	n := len(digs)
	m := float64(n) / 2.0
	piv := int(math.Ceil(m))
	mid := n - piv
	pal := make([]int, n)
	for i := 0; i < piv; i++ {
		pal[i] = digs[i]
		pal[n-i-1] = digs[i]
	}
	npal := digsToInt(pal)
	if npal > dig {
		return npal
	} else if npal == dig {
		nFronts := getDigits(digsToInt(digs[:piv]) + 1)
		nBacks := make([]int, mid)
		for i := 0; i < mid; i++ {
			nBacks[mid-i-1] = nFronts[i]
		}
		nFronts = append(nFronts, nBacks...)
		return digsToInt(nFronts)
	} else {
		return nextPalindrome(npal)
	}
	return 0
}

func isFairNSquare(dig int64) bool {
	chk, ok := isFairNSquareList[dig]
	if ok {
		return chk
	}

	sqr := math.Sqrt(float64(dig))
	if sqr == math.Trunc(sqr) && isPalindrome(int64(sqr)) {
		isFairNSquareList[dig] = true
		return true
	}

	isFairNSquareList[dig] = false
	return false
}

func nextFairNSquare(dig int64) int64 {
	next := dig
	for {
		next = nextPalindrome(next)
		if isFairNSquare(next) {
			return next
		}
	}
	return next
}
