package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	X = 88
	O = 79
	N = 46
	T = 84
)

type Table struct {
	grid [][]int
	id   int
}

func NewTable(id int, r1, r2, r3, r4 []byte) *Table {
	table := new(Table)
	table.id = id
	table.grid = make([][]int, 4, 4)

	var s []byte
	var c int
	for i, _ := range table.grid {
		table.grid[i] = make([]int, 4, 4)
		switch i {
		case 0:
			s = r1
		case 1:
			s = r2
		case 2:
			s = r3
		case 3:
			s = r4
		}
		for ii, v := range s {
			c = int(v)
			if c == 46 {
				c = 0
			}
			table.grid[i][ii] = c
		}
	}
	return table
}

func (t *Table) String() string {
	return fmt.Sprintf("%v\n", t.grid)
}

func main() {
	// Parse args for file name
	flag.Parse()
	fn := flag.Arg(0)

	// Open file
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	// Create reader for file
	rd := bufio.NewReader(f)

	var l, r1, r2, r3 []byte
	i, n, sr := 0, 0, 1
	var N int
	var tables []*Table

	// Parse file
	for {
		l, _, err = rd.ReadLine()
		if err == io.EOF {
			break
		}

		switch {
		case i == 0:
			// First row
			N_, _ := strconv.ParseInt(string(l), 0, 0)
			N = int(N_)
			tables = make([]*Table, N, N)
		case i == sr:
			// row 1
			r1 = l
		case i == sr+1:
			// row 2
			r2 = l
		case i == sr+2:
			// row 3
			r3 = l
		case i == sr+3:
			// row 4
			tables[n] = NewTable(n+1, r1, r2, r3, l)
			n++
			sr += 5
		}

		i++
	}

	// Output
	of, err := os.OpenFile(strings.Replace(fn, ".in", ".out", -1), os.O_WRONLY|os.O_CREATE, 0666)
	defer of.Close()
	if err != nil {
		panic(err)
	}
	done := func(tab *Table, res string) {
		of.Write([]byte(fmt.Sprintf("Case #%d: %s\n", tab.id, res)))
	}

	var result string
	var pX, pO, tot, ne int
	for _, tab := range tables {
		pX, pO, ne = 0, 0, 0
		// Check rows
		for i, _ = range tab.grid {
			tot = 0
			for _, v := range tab.grid[i] {
				tot += v
				if v == 0 {
					ne++
				}
			}
			switch tot {
			case X * 4, X*3 + T:
				pX++
			case O * 4, O*3 + T:
				pO++
			}
		}

		// Check columns
		for ci, _ := range tab.grid {
			tot = 0
			for ri, _ := range tab.grid {
				tot += tab.grid[ri][ci]
			}
			switch tot {
			case X * 4, X*3 + T:
				pX++
			case O * 4, O*3 + T:
				pO++
			}
		}

		// Check diagonals
		tot = tab.grid[0][0] + tab.grid[1][1] + tab.grid[2][2] + tab.grid[3][3]
		switch tot {
		case X * 4, X*3 + T:
			pX++
		case O * 4, O*3 + T:
			pO++
		}

		tot = tab.grid[3][0] + tab.grid[2][1] + tab.grid[1][2] + tab.grid[0][3]
		switch tot {
		case X * 4, X*3 + T:
			pX++
		case O * 4, O*3 + T:
			pO++
		}

		// Declare result
		switch {
		case pX == pO && ne == 0:
			result = "Draw"
		case pX > pO:
			result = fmt.Sprintf("%s won", string(X))
		case pO > pX:
			result = fmt.Sprintf("%s won", string(O))
		case ne > 0:
			result = "Game has not completed"
		}
		done(tab, result)
	}
	of.Write([]byte("\n"))
}
