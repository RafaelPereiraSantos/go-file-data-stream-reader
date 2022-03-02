package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type (
	// dataChunkHandler, function that will handle the data as soon as it is determinated by the dataChunkDelimiter
	// function.
	dataChunkHandler func([]byte) error

	// dataChunkDelimiter, function that determinates the size of the chunk that is going to be processed, it receives a
	// byte array and should return "false", the original "byte array" paramenter and "nil" in case the chunk is not
	// enought or return "true" following by the chunk to be processed and the left over bytes that should not be
	// processed at least for now.
	// NOTE: the boolean returned is in case that the byte array that was send is enough to be processed and there is no
	// left overs to return.
	dataChunkDelimiter func([]byte) (bool, []byte, []byte)

	errorHandler func(error)
)

const newLineByte = byte('\n')

func main() {
	sizeOfTheChunkToBeFetched := 2
	dataSource, _ := os.Open("data_input_example.txt")

	chunkHandler := func(b []byte) error {
		msg := fmt.Sprintf("Text: [%s], size: [%d]", string(b), len(b))
		fmt.Println(msg)
		return nil
	}

	err := processDataSourceInChunks(dataSource, sizeOfTheChunkToBeFetched, chunkHandler, delimiteByNewLine)

	if err != nil {
		log.Fatalf("Exit due to [%v]", err)
	}

}

// processDataSourceInChunks, it is a function that will split a byte array in chunks of data to process each part at a
// time allowing large files to be processed in small parts avoiding large ammounts of memory to be allocation. This
// method is primarily focused on dealing with files containing JSON data splited in lines.
func processDataSourceInChunks(
	dataSource io.Reader,
	chunkSize int,
	chunkHandler dataChunkHandler,
	chunkDelimiter dataChunkDelimiter) error {
	leftOver := make([]byte, 0)
	eof := false

	for {
		var err error
		enoughDataInChunkToBeProcessed := false
		chunkToBeProcessed := make([]byte, 0, chunkSize+1)

		// This loop is used to retrieve small parts of the data from the io.Reader then check if all the data fetched
		// so far is enough to be considered a "chunk" by applying the dataChunkDelimiter function of the data so far
		// collected every time a new part is retrieved.
		for {
			tempChunk := make([]byte, chunkSize, chunkSize+1)

			checkLeftOverFirst := len(leftOver) > 0

			if checkLeftOverFirst {
				tempChunk = leftOver
				leftOver = make([]byte, 0)
			} else {
				_, err = dataSource.Read(tempChunk)
			}

			if err != nil {

				eof = err == io.EOF

				if eof {
					break
				}

				return err
			}

			chunkToBeProcessed = append(chunkToBeProcessed, tempChunk...)

			enoughDataInChunkToBeProcessed, chunkToBeProcessed, leftOver = chunkDelimiter(chunkToBeProcessed)

			if enoughDataInChunkToBeProcessed {
				break
			}
		}

		chunkWithoutNewLine := removeNewLine(chunkToBeProcessed)

		err = chunkHandler(chunkWithoutNewLine)

		if err != nil {
			return err
		}

		if eof {
			break
		}
	}

	return nil
}

func removeNewLine(b []byte) []byte {
	bParts := bytes.Split(b, []byte("\n"))

	if len(bParts) == 1 {
		return b
	}

	if len(bParts) == 0 {
		return []byte{}
	}

	bWithNoNewLine := make([]byte, 0)

	for _, part := range bParts {
		bWithNoNewLine = append(bWithNoNewLine, part...)
	}

	return bWithNoNewLine
}

// delimiteByNewLine, one implementaiton of dataChunkDelimiter, this function will receive a byte array as parameter and
// will try to determinete whether or not this chunk of data is enough to be processed by checking by a new line "\n"
// character at any point of the array, all data before the new line will be considered an complete chunk, part after
// the new line will be considered as left overs.
func delimiteByNewLine(chunk []byte) (bool, []byte, []byte) {
	chunkCopy := make([]byte, len(chunk), len(chunk)+1)
	copy(chunkCopy, chunk)

	// by splitting the chunk using a reparator as new line, we could define the chunk and the left over by choosing
	// the first index as the chunk and all the other elements as left overs.
	parts := bytes.Split(chunkCopy, []byte{newLineByte})

	thereIsLeftOver := len(parts) > 1

	if thereIsLeftOver {
		leftOver := make([]byte, len(parts[1]))

		chunkToBeProcessed := parts[0]
		leftOverParts := parts[1:]

		// all leftover must be concatenated and a new line should be add at the end each part in order to return it
		// as it was given to the method as parameter so it could be iterated again in future usages.
		for i, part := range leftOverParts {
			partLen := len(part)

			if partLen == 0 {
				continue
			}

			leftOver = append(leftOver, part...)

			// no new line should be add at the last index to prevent adding new lines at parts of text that does not
			// contain them.
			if i < partLen-1 {
				leftOver = append(leftOver, newLineByte)
			}
		}

		return true, chunkToBeProcessed, leftOver
	}

	return false, chunk, nil
}
