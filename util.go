package logrotate

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

const (
	compressSuffix = ".gz"
)

func stringToBytes(s string) (bytes int64, err error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	if s == "" {
		return 0, nil
	}
	var numStr, unit string
	i := strings.IndexFunc(s, unicode.IsLetter)
	if i == -1 {
		numStr = s
	} else {
		numStr, unit = s[:i], s[i:]
	}
	bytes, err = strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return
	}
	if bytes <= 0 {
		bytes = 0
		return
	}
	switch unit {
	case "E", "EB", "EIB":
		bytes = bytes * EXABYTE
	case "P", "PB", "PIB":
		bytes = bytes * PETABYTE
	case "T", "TB", "TIB":
		bytes = bytes * TERABYTE
	case "G", "GB", "GIB":
		bytes = bytes * GIGABYTE
	case "M", "MB", "MIB":
		bytes = bytes * MEGABYTE
	case "K", "KB", "KIB":
		bytes = bytes * KILOBYTE
	case "B":
	default:
	}
	return
}

func archiveName(name string, timeFormat string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	prefix, ext := splitFilename(filename)
	t := time.Now().Format(timeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, t, ext))
}

func splitFilename(name string) (prefix, ext string) {
	ext = filepath.Ext(name)
	prefix = name[:len(name)-len(ext)]
	return
}

func timeFromName(timeFormat, filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(timeFormat, ts)
}

func compressFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open log file, error(%v)", err)
	}
	defer f.Close()
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat log file, error(%v)", err)
	}
	dst := path + compressSuffix
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file, error(%v)", err)
	}
	defer gzf.Close()
	gz := gzip.NewWriter(gzf)
	defer func() {
		if err != nil {
			_ = os.Remove(dst)
			err = fmt.Errorf("failed to compress file, error(%v)", err)
		}
	}()
	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
