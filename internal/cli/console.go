package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Console struct {
	reader *bufio.Reader
	writer io.Writer
}

func NewConsole(reader io.Reader, writer io.Writer) *Console {
	return &Console{
		reader: bufio.NewReader(reader),
		writer: writer,
	}
}

func (c *Console) ReadLine(prompt string) (string, error) {
	if prompt != "" {
		fmt.Fprint(c.writer, prompt)
	}

	input, err := c.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(input) > 0 {
			return strings.TrimSpace(input), nil
		}

		return "", err
	}

	return strings.TrimSpace(input), nil
}
