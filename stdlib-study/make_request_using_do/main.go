package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir := flag.String("dir", ".", "directory to save file")
	timeout := flag.Duration("timeout", 30*time.Second, "timeout for download")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatal("usage: download [-timeout duration] url filename")
	}

	url, filename := args[0], args[1]

	c := http.Client{
		Timeout: *timeout,
	}

	if err := downloadAndSave(context.TODO(), &c, url, filename, dir); err != nil {
		log.Fatal(err)
	}
}

func downloadAndSave(ctx context.Context, c *http.Client, url, dst string, dir *string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: GET %q: %v", url, err)
	}
	// `Do` serializes a http.Request, sends it to the server
	// and then deserializes the response to a http.Response
	resp, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status)
	}

	dstPath := filepath.Join(*dir, dst)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, resp.Body); err != nil {
		return fmt.Errorf("copying response to file: %v", err)
	}
	return err
}
