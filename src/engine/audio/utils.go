package audio

import (
	"bytes"
	"io"
)

// ループ可能なストリームを作成するためのラッパー
type loopStream struct {
	*bytes.Reader
	originalData []byte
}

func newLoopStream(data []byte) *loopStream {
	return &loopStream{
		Reader:       bytes.NewReader(data),
		originalData: data,
	}
}

func (s *loopStream) Read(p []byte) (n int, err error) {
	n, err = s.Reader.Read(p)
	if err == io.EOF {
		s.Reader.Reset(s.originalData)
		remaining := len(p) - n
		if remaining > 0 {
			m, err := s.Reader.Read(p[n:])
			n += m
			if err != nil && err != io.EOF {
				return n, err
			}
		}
		return n, nil
	}
	return n, err
} 