package getRegionBoundaries

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestGetBoundries(t *testing.T) {
	file, err := func(url, filename string) (file *os.File, err error) {
		// don't worry about errors
		response, _ := http.Get(url)
		defer response.Body.Close()
		file, e := os.Create(path.Join(os.TempDir(), filename))
		if e != nil {
			err = e
			log.Fatal(e)
			return
		}
		io.Copy(file, response.Body)
		defer file.Close()
		return file, nil
	}("http://i.imgur.com/m1UIjW1.jpg", "m1UIjW1.jpg")
	if err != nil {
		t.Error("error on file operation")
		return
	}
	log.Println(file.Name())

	file, _ = os.Open(file.Name())
	boundaries, _ := getRegionBoundaries(file, []int{0, 64, 96, 128, 160, 192}, []int{0, 64, 96, 128, 160, 192}, []int{0, 128}, 16)
	for _, rowData := range boundaries {
		for _, colData := range rowData {
			if colData == 1 {
				fmt.Printf("%d ", colData)
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println("")
	}
}

func TestGenerateColorNodeTree(t *testing.T) {
	testArray := make([]int, 100)
	for index := range testArray {
		testArray[index] = rand.Intn(255)
	}
	// log.Println(testArray)
	root := generateColorNodeTree(testArray)
	ch := make(chan int)
	go func() {
		last := -1
		for val := range ch {
			// log.Printf("%d,", val)
			if last >= val {
				t.Error("wrong order")
			}
		}
		// log.Println("")
	}()
	getRank(ch, root)
}

func getRank(ch chan int, node *ColorNode) {
	if node.left != nil {
		getRank(ch, node.left)
	}
	ch <- node.value
	if node.right != nil {
		getRank(ch, node.right)
	}
}
