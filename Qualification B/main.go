package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type lawn [][]int

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
		s := fmt.Sprintf(format, ex...)
		fmt.Print(s)
		of.Write([]byte(s))
	}

	// Parse input file
	N_, _ := strconv.ParseInt(lines[0], 0, 0)
	N := int(N_)
	lawns := make([]lawn, N)
	refs := make([]lawn, N)

	var words []string
	var nums []int
	var cont bool
	var rn, cn, ri, ci int
	r := 1

	getnums := func(in string) {
		words = strings.Split(in, " ")
		nums = make([]int, len(words))
		for i, v := range words {
			n_, _ := strconv.ParseInt(v, 0, 0)
			nums[i] = int(n_)
		}
	}

	mow := func(r, c int, ref, l lawn) bool {
		// Check row
		target := ref[r][c]
		for _, v := range ref[r] {
			if v > target {
				// Check column
				for ri := 0; ri < rn; ri++ {
					if ref[ri][c] > target {
						return false
					}
				}
				// No problem, mow column
				for ri := 0; ri < rn; ri++ {
					if l[ri][c] > target {
						l[ri][c] = target
					}
				}
				return true
			}
		}
		// No problem, mow row
		for ci, _ := range l[r] {
			if l[r][ci] > target {
				l[r][ci] = target
			}
		}
		return true
	}

	for i := 0; i < N; i++ {
		// Get row, col nums
		getnums(lines[r])
		rn, cn = int(nums[0]), int(nums[1])
		r++

		// Fill lawn table
		refs[i] = make(lawn, rn)
		lawns[i] = make(lawn, rn)
		for ri, _ := range refs[i] {
			getnums(lines[r])
			refs[i][ri] = nums
			lawns[i][ri] = make([]int, cn)
			for ci, _ = range lawns[i][ri] {
				lawns[i][ri][ci] = 100
			}
			r++
		}

		cont = true
		// Mow
		for ri = 0; ri < rn; ri++ {
			for ci = 0; ci < cn; ci++ {
				mow(ri, ci, refs[i], lawns[i])
			}
		}
		// Check if mowed lawn correct
		for ri = 0; ri < rn; ri++ {
			for ci = 0; ci < cn; ci++ {
				if refs[i][ri][ci] != lawns[i][ri][ci] {
					cont = false
					break
				}
			}
			if !cont {
				out("Case #%d: %s\n", i+1, "NO")
				break
			}
		}
		if cont {
			out("Case #%d: %s\n", i+1, "YES")
		}
	}
}
