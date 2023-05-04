package processor

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/damirm/links-warehouse/internal/warehouse"
)

func PopulateQueueFromDirectory(directory string, s *warehouse.WarehouseService) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	queueSize := 0

	for _, fileInfo := range files {
		filePath := filepath.Join(directory, fileInfo.Name())
		log.Printf("reading file: %s", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			link := scanner.Text()

			// TODO: Only habr links is supported for now.
			if !strings.Contains(link, "habr") {
				continue
			}

			queueSize++
			u, err := url.Parse(link)
			if err != nil {
				return err
			}
			err = s.QueueLink(u)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("added %d links to queue", queueSize)

	return nil
}
