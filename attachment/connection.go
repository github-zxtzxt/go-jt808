package attachment

import (
	"errors"
	"io"
	"net"
)

type connection struct {
	conn        net.Conn
	streamFunc  func() StreamDataHandler
	handle      JT808DataHandler
	fileEventer FileEventer
}

func newConnection(c net.Conn, streamFunc func() StreamDataHandler,
	handle JT808DataHandler, fileEventer FileEventer) *connection {
	return &connection{
		conn:        c,
		streamFunc:  streamFunc,
		handle:      handle,
		fileEventer: fileEventer,
	}
}

func (c *connection) run() {
	var (
		curData  = make([]byte, 100*1024)
		progress = &PackageProgress{
			ProgressStage: ProgressStageInit,
			Record:        map[string]*Package{},
			historyData:   make([]byte, 0),
			handle:        c.handle,
			streamFunc:    c.streamFunc,
		}
	)
	defer func() {
		progress.ProgressStage = ProgressStageSuccessQuit
		if progress.ExtensionFields.Err != nil {
			progress.ProgressStage = ProgressStageFailQuit
		}
		c.fileEventer.OnEvent(progress)
		clear(curData)
	}()

	for {
		if n, err := c.conn.Read(curData); err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			progress.ExtensionFields.Err = err
			return
		} else if n > 0 {
			progress.historyData = append(progress.historyData, curData[:n]...)
			for err := range progress.iter() {
				if err == nil {
					if progress.hasJT808Reply() {
						data, err := progress.handle.ReplyData()
						progress.ExtensionFields.RecentPlatformData = data
						progress.ExtensionFields.Err = err
						if _, err := c.conn.Write(data); err != nil {
							progress.ExtensionFields.Err = errors.Join(err, progress.ExtensionFields.Err)
							return
						}
					}
					c.fileEventer.OnEvent(progress)
				} else if errors.Is(err, ErrInsufficientDataLen) {
					// 是数据长度不够的错误情况 就结束
					break
				} else {
					progress.ExtensionFields.Err = err
					return
				}
			}
		}

	}
}