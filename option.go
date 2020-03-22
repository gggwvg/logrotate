package logrotate

import (
	"errors"
	"os"
	"path/filepath"
)

type period string

const (
	// PeriodHourly rotates log every hour
	PeriodHourly period = "hourly"
	// PeriodDaily rotates log by every day
	PeriodDaily period = "daily"
	// PeriodWeekly rotates log by every week
	PeriodWeekly period = "weekly"
	// PeriodMonthly rotates log by every month
	PeriodMonthly period = "monthly"
)

var (
	DefaultArchiveTimeFormat = "2006-01-02_15:04:05"
	DefaultMaxArchives       = 100
	DefaultMaxArchiveDays    = 14
)

type Options struct {
	// Filename is the file to write logs to.
	// It uses <process name>.log in os.TempDir() if empty.
	Filename string `json:"filename" toml:"filename" yaml:"filename"`

	// RotatePeriod is time period for rotate log.
	// It supports hourly, daily, weekly, monthly.
	RotatePeriod string `json:"rotate_period" toml:"rotate_period" yaml:"rotate_period"`

	// RotateSize is the maximum size of the log file before it gets rotated
	RotateSize string `json:"rotate_size" toml:"rotate_size" yaml:"rotate_size"`

	// MaxArchives is the maximum number of old log files to retain
	MaxArchives int `json:"max_archives" toml:"max_archives" yaml:"max_archives"`

	// MaxArchiveDays is the maximum number of days to archived files
	MaxArchiveDays int `json:"max_archive_days" toml:"max_archive_days" yaml:"max_archive_days"`

	// ArchiveTimeFormat is the format of the archived files
	ArchiveTimeFormat string `json:"archive_time_format" toml:"archive_time_format" yaml:"archive_time_format"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" toml:"compress" yaml:"compress"`

	rotateSize int64
	cron       string
}

func (o *Options) Apply() error {
	if o.RotateSize != "" {
		size, err := stringToBytes(o.RotateSize)
		if err != nil {
			return err
		}
		o.rotateSize = size
	}
	if o.RotatePeriod != "" {
		switch period(o.RotatePeriod) {
		case PeriodHourly:
			o.cron = "0 * * * *"
		case PeriodDaily:
			o.cron = "0 0 * * *"
		case PeriodWeekly:
			o.cron = "0 0 * * 0"
		case PeriodMonthly:
			o.cron = "0 0 1 * *"
		default:
			return errors.New("invalid rotate period")
		}
	}
	if o.RotateSize != "" && o.RotatePeriod != "" {
		if o.ArchiveTimeFormat == "" {
			o.ArchiveTimeFormat = DefaultArchiveTimeFormat
		}
		if o.MaxArchives <= 0 {
			o.MaxArchives = DefaultMaxArchives
		}
		if o.MaxArchiveDays <= 0 {
			o.MaxArchiveDays = DefaultMaxArchiveDays
		}
	}
	if o.Filename == "" {
		name := filepath.Base(os.Args[0]) + ".log"
		o.Filename = filepath.Join(os.TempDir(), name)
	}
	return nil
}

type OptionFunc func(*Options)

func Filename(name string) OptionFunc {
	return func(opts *Options) {
		opts.Filename = name
	}
}

func RotatePeriod(p period) OptionFunc {
	return func(opts *Options) {
		opts.RotatePeriod = string(p)
	}
}

func RotateSize(size string) OptionFunc {
	return func(opts *Options) {
		opts.RotateSize = size
	}
}

func ArchiveTimeFormat(format string) OptionFunc {
	return func(opts *Options) {
		opts.ArchiveTimeFormat = format
	}
}

func MaxArchives(number int) OptionFunc {
	return func(opts *Options) {
		opts.MaxArchives = number
	}
}

func MaxArchiveDays(days int) OptionFunc {
	return func(opts *Options) {
		opts.MaxArchiveDays = days
	}
}

func Compress(compress bool) OptionFunc {
	return func(opts *Options) {
		opts.Compress = compress
	}
}
