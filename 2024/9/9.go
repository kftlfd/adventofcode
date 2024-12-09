package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	inp_path := os.Args[1]
	inp_file, err := os.ReadFile(inp_path)
	if err != nil {
		panic(err)
	}

	disk_map := string(inp_file)
	disk_layout := []string{}
	file_id := 0

	for i, c := range strings.Split(disk_map, "") {
		if c == "\n" {
			break
		}
		cnt, err := strconv.Atoi(c)
		if err != nil {
			panic(fmt.Errorf("not a number: %v", c))
		}
		if i%2 == 0 {
			cur_file_id := strconv.Itoa(file_id)
			file_id += 1
			for j := 0; j < cnt; j++ {
				disk_layout = append(disk_layout, cur_file_id)
			}
		} else {
			for j := 0; j < cnt; j++ {
				disk_layout = append(disk_layout, ".")
			}
		}
	}

	d1 := make([]string, len(disk_layout))
	copy(d1, disk_layout)
	// fmt.Println(d1)

	compact(d1)
	// fmt.Println(disk_layout)

	ans := checksum(d1)
	fmt.Println("checksum:", ans)

	//
	// part 2
	//

	d2 := make([]string, len(disk_layout))
	copy(d2, disk_layout)
	// fmt.Println("d2:", d2)

	chunks := to_chunks(d2)
	// fmt.Println("chunks:", chunks)

	chunks_compacted := compact_chunks(chunks)
	// fmt.Println("comp chunks:", chunks_compacted)

	d2_compacted := from_chunks(chunks_compacted)
	// fmt.Println("d2 comp:", d2_compacted)

	d2_checksum := checksum(d2_compacted)
	fmt.Println("checksum 2:", d2_checksum)
}

func compact(disk []string) {
	left := 0
	right := len(disk) - 1
	for left < right {
		if disk[left] != "." {
			left += 1
			continue
		}
		if disk[right] == "." {
			right -= 1
			continue
		}
		disk[left], disk[right] = disk[right], disk[left]
	}
}

func checksum(disk []string) int {
	sum := 0
	for position, c := range disk {
		if c == "." {
			continue
		}
		id, err := strconv.Atoi(c)
		if err != nil {
			panic(fmt.Errorf("not a number: %v", c))
		}
		sum += position * id
	}
	return sum
}

type Chunk struct {
	id  string
	len int
}

func to_chunks(disk []string) []Chunk {
	out := []Chunk{
		{id: disk[0], len: 0},
	}
	for _, c := range disk {
		if c == out[len(out)-1].id {
			out[len(out)-1].len += 1
		} else {
			out = append(out, Chunk{id: c, len: 1})
		}
	}
	return out
}

func from_chunks(chunks []Chunk) []string {
	disk := []string{}
	for _, chunk := range chunks {
		for i := 0; i < chunk.len; i++ {
			disk = append(disk, chunk.id)
		}
	}
	return disk
}

func compact_chunks(chunks []Chunk) []Chunk {
	out := make([]Chunk, len(chunks))
	copy(out, chunks)

	right := len(chunks) - 1
	for right > 0 {
		cur := out[right]
		if cur.id == "." {
			right -= 1
			continue
		}

		left := 0
		free_space_found := false
		for left < right {
			if out[left].id == "." && out[left].len >= cur.len {
				free_space_found = true
				break
			}
			left += 1
		}

		if !free_space_found {
			right -= 1
			continue
		}

		extra_free_space := out[left].len - cur.len

		// fmt.Printf("-----\n%v\n", out)
		if extra_free_space == 0 {
			out[left], out[right] = out[right], out[left]
			right -= 1
		} else {
			new_out := make([]Chunk, 0)
			new_out = append(new_out, out[:left]...)
			new_out = append(new_out, cur)
			new_out = append(new_out, Chunk{id: ".", len: extra_free_space})
			new_out = append(new_out, out[left+1:right]...)
			new_out = append(new_out, Chunk{id: ".", len: cur.len})
			new_out = append(new_out, out[right+1:]...)
			out = new_out
		}
		// fmt.Printf("%v\n", out)
	}

	return out
}
