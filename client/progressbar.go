package client

import "fmt"

type progressOutput struct {
	total int
}

func (p *progressOutput) writeProgress(val int) {
	progress := int((float32(val) / 254) * 100)

	fmt.Printf("\r%d", progress)

	if progress == 100 {
		fmt.Print("\n")
	}
}
