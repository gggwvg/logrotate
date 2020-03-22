package logrotate

import (
	"compress/gzip"
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
		bytes *= 1 << 60
	case "P", "PB", "PIB":
		bytes *= 1 << 50
	case "T", "TB", "TIB":
		bytes *= 1 << 40
	case "G", "GB", "GIB":
		bytes *= 1 << 30
	case "M", "MB", "MIB":
		bytes *= 1 << 20
	case "K", "KB", "KIB":
		bytes *= 1 << 10
	case "B":
	default:
	}
	return
}

func archiveName(name string, timeFormat string) string {
	dir := filepath.Dir(name)
	prefix, ext := splitFilename(name)
	t := time.Now().Format(timeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, t, ext))
}

func splitFilename(name string) (prefix, ext string) {
	filename := filepath.Base(name)
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)]
	return
}

func timeFromName(timeFormat, filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, fmt.Errorf("mismatched prefix(%s) filename(%s)", prefix, filename)
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, fmt.Errorf("mismatched extension(%s) filename(%s)", ext, filename)
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
