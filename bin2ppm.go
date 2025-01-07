package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Println(`Usage: bin2ppm [-d] [x z]

Options:
  -d        Decode PPM to binary

Arguments:
  x z       Dimensions of the image (default: 32 32) for encoding

Examples:
  Encode: cat input.bin | bin2ppm > output.ppm
  Decode: cat input.ppm | bin2ppm -d > output.bin`)
}

func main() {
	decode := flag.Bool("d", false, "Decode PPM to binary")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	width, height := 32, 32 // Default dimensions

	if len(args) == 2 && !*decode {
		var err error
		width, err = strconv.Atoi(args[0])
		if err != nil || width <= 0 {
			fmt.Fprintln(os.Stderr, "Invalid width:", args[0])
			os.Exit(1)
		}
		height, err = strconv.Atoi(args[1])
		if err != nil || height <= 0 {
			fmt.Fprintln(os.Stderr, "Invalid height:", args[1])
			os.Exit(1)
		}
	} else if len(args) > 0 {
		usage()
		os.Exit(1)
	}

	if *decode {
		if err := decodePPM(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding PPM:", err)
			os.Exit(1)
		}
	} else {
		if err := encodeBinaryToPPM(os.Stdin, os.Stdout, width, height); err != nil {
			fmt.Fprintln(os.Stderr, "Error encoding binary to PPM:", err)
			os.Exit(1)
		}
	}
}

func encodeBinaryToPPM(input io.Reader, output io.Writer, width, height int) error {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	_, err := fmt.Fprintf(writer, "P3\n# Created by bin2ppm.go - https://github.com/706f6c6c7578/bin2ppm\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		b := scanner.Bytes()[0]
		if _, err := fmt.Fprintf(writer, "%d\n", b); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func decodePPM(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)

	// Read PPM header
	if !scanner.Scan() || scanner.Text() != "P3" {
		return errors.New("invalid PPM format")
	}

	// Skip comments
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		break
	}

	// Read dimensions
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return errors.New("invalid PPM dimensions")
	}

	// Skip max color value
	if !scanner.Scan() || scanner.Text() != "255" {
		return errors.New("invalid PPM max color value")
	}

	// Decode pixel values
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return errors.New("invalid pixel value")
		}
		if value < 0 || value > 255 {
			return errors.New("pixel value out of range")
		}
		if _, writeErr := writer.Write([]byte{byte(value)}); writeErr != nil {
			return writeErr
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
