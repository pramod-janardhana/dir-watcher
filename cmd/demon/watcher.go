package demon

import (
	"github.com/knadh/koanf/providers/file"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

func Watch(file *file.File, callback func()) {
	file.Watch(func(event interface{}, err error) {
		if err != nil {
			zlog.Errorf("error while watching config file: %s", err.Error())
			return
		}

		zlog.Debugf("config has changed")
		callback()
	})
}
