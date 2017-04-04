package getRegionBoundaries

import (
	"image"
	_ "image/jpeg" // load jpeg image decoder
	_ "image/png"  // load png image decoder
	"math/rand"

	"log"

	"io"
)

func init() {
}

func insertVal(val int, node *ColorNode) {
	if node.value == val {
		return
	}
	// should be at right sub-tree
	if node.value < val {
		if node.right != nil {
			insertVal(val, node.right)
		} else {
			right := new(ColorNode)
			right.value = val
			node.right = right
		}
	} else {
		if node.left != nil {
			insertVal(val, node.left)
		} else {
			left := new(ColorNode)
			left.value = val
			node.left = left
		}
	}
}
func generateColorNodeTree(segs []int) *ColorNode {
	rootIndex := rand.Intn(len(segs))
	root := new(ColorNode)
	root.value = segs[rootIndex]
	for _, val := range segs {
		insertVal(val, root)
	}
	calcRank(-1, root)
	return root
}
func calcRank(last int, node *ColorNode) int {
	if node.left == nil {
		node.rank = last + 1
		return node.rank
	}
	node.rank = calcRank(last, node.left) + 1
	if node.right != nil {
		return calcRank(node.rank, node.right)
	}
	return node.rank
}

func findRank(node *ColorNode, colorVal int) int {
	if node.value == colorVal {
		return node.value
	}
	if node.left != nil && node.value > colorVal {
		return findRank(node.left, colorVal)
	}
	if node.right != nil && node.value < colorVal {
		return findRank(node.right, colorVal)
	}
	return node.value
}
func getRegionBoundaries(f io.Reader, redSegs, greenSegs, blueSegs []int, sampleEach int) (boundaries [][]int, e error) {
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("unable to decode image", err)
		e = err
		return
	}

	redColorNodeTree := generateColorNodeTree(redSegs)
	greenColorNodeTree := generateColorNodeTree(greenSegs)
	blueColorNodeTree := generateColorNodeTree(blueSegs)

	dividers := make([][][]int, len(redSegs))
	for green := range dividers {
		dividers[green] = make([][]int, len(greenSegs))
		for blue := range dividers[green] {
			dividers[green][blue] = make([]int, len(blueSegs))
		}
	}
	boundaries = make([][]int, img.Bounds().Dy()/sampleEach*2-1)
	for y := range boundaries {
		boundaries[y] = make([]int, img.Bounds().Dx()/sampleEach*2-1)
	}
	for y := 0; y < len(boundaries); y += 2 {
		rowData := boundaries[y]
		for x := 0; x < len(rowData); x += 2 {
			r, g, b, a := img.At(x*sampleEach/2, y*sampleEach/2).RGBA()
			// fmt.Printf("[x:%d, y:%d, r:%d, g:%d, b:%d,a %d]", x*sampleEach/2, y*sampleEach/2, r*255/a, g*255/a, b*255/a, a)
			rowData[x] = findRank(redColorNodeTree, 1<<24+int(r*255/a))<<16 + findRank(greenColorNodeTree, int(g*255/a))<<8 + findRank(blueColorNodeTree, int(b*255/a))
		}
		// fmt.Println("")
	}

	//find boundaries
	for y := 1; y < len(boundaries); y++ {
		rowData := boundaries[y]
		for x := 1; x < len(rowData); x++ {
			if x%2 == 1 && y%2 == 0 && x+1 < len(rowData) && rowData[x-1] != rowData[x+1] {
				rowData[x] = 1
			}
			if x%2 == 0 && y%2 == 1 && y+1 < len(boundaries) && boundaries[y-1][x] != boundaries[y+1][x] {
				rowData[x] = 1
			}
		}
	}

	//to fill the gaps between edge pixels
	for y := 1; y < len(boundaries); y++ {
		rowData := boundaries[y]
		for x := 1; x < len(rowData); x++ {
			if x+1 < len(rowData) && rowData[x-1] == 1 && rowData[x+1] == 1 {
				rowData[x] = 1
			}
			if y+1 < len(boundaries) && boundaries[y-1][x] == 1 && boundaries[y+1][x] == 1 {
				rowData[x] = 1
			}
		}
	}

	//find junction points
	for y := 1; y < len(boundaries); y += 2 {
		rowData := boundaries[y]
		for x := 1; x < len(rowData); x += 2 {
			if x+1 < len(rowData) && y+1 < len(boundaries) && rowData[x-1]+rowData[x+1]+boundaries[y-1][x]+boundaries[y+1][x] > 2 {
				rowData[x] = 1
			}
		}
	}

	// clear
	// for y := 0; y < len(boundaries); y += 2 {
	// 	rowData := boundaries[y]
	// 	for x := 0; x < len(rowData); x += 2 {
	// 		rowData[x] = 0
	// 	}
	// }
	return
}
