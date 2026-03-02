package common

import (
	"errors"
	"strconv"
	"strings"
)

// ParseRangeHeader 解析Range请求头，返回起始/结束字节
func ParseRangeHeader(rangeHeader string, fileSize int64) (start, end int64, err error) {
	if rangeHeader == "" {
		return 0, 0, nil // 无Range请求，返回整个文件
	}

	// Range格式：bytes=0-1023
	parts := strings.SplitN(rangeHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "bytes" {
		return 0, 0, errors.New("invalid Range header")
	}

	rangeParts := strings.SplitN(parts[1], "-", 2)
	if len(rangeParts) != 2 {
		return 0, 0, errors.New("invalid Range range")
	}

	// 解析起始字节
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil || start < 0 {
			return 0, 0, errors.New("invalid start offset")
		}
	}

	// 解析结束字节
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil || end < start {
			return 0, 0, errors.New("invalid end offset")
		}
	} else {
		end = fileSize - 1 // 无结束字节，到文件末尾
	}

	// 边界检查
	if start >= fileSize {
		return 0, 0, errors.New("start offset beyond file size")
	}
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}
