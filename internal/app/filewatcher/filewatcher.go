package filewatcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pramod-janardhana/dir-watcher/config"
	"github.com/pramod-janardhana/dir-watcher/dto/filewatcher"
	serror "github.com/pramod-janardhana/dir-watcher/internal/server/error"
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

type Controller struct {
}

func NewController() Controller {
	return Controller{}
}

func (c Controller) UpdateConfig(req *filewatcher.UpdateConfigReq) (*filewatcher.UpdateConfigRes, error) {
	// verifying cron expression
	{
		_, err := os.Stat(req.DirOrFileToWatch)
		if err != nil {
			zlog.Errorf("failed to verify path(%s): %s", req.DirOrFileToWatch, err.Error())
			return nil, serror.NewError(http.StatusBadRequest, "could find the path, please verify if the path exists and try again")
		}
	}

	// updating the config file
	{
		path, err := config.FileWatcherConfigPath()
		if err != nil {
			zlog.Errorf("could not get config.filewatcher.windows.json file path: %v", err)
			return nil, serror.NewError(http.StatusInternalServerError, "could not find the config file, please contact support")
		}

		// read the config file
		data, err := ioutil.ReadFile(path)
		if err != nil {
			zlog.Errorf("error reading config file: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// unmarshal the JSON data into a struct
		var conf = config.FileWatcher{}
		if err := json.Unmarshal(data, &conf); err != nil {
			zlog.Errorf("error unmarshalling JSON data: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// updating the fields in the struct
		{
			if req.DirOrFileToWatch != "" {
				conf.DirOrFileToWatch = req.DirOrFileToWatch
			}
		}

		// marshal the struct back to json
		updatedJSON, err := json.MarshalIndent(conf, "", "	")
		if err != nil {
			zlog.Errorf("error marshalling json data: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}

		// writing the updated JSON data back to the file
		if err := ioutil.WriteFile(path, updatedJSON, 0644); err != nil {
			zlog.Errorf("error writing updated JSON to config file: %s", err.Error())
			return nil, serror.NewError(http.StatusInternalServerError, "something went wrong, please try again later")
		}
	}

	toReturn := filewatcher.UpdateConfigRes("restarting service with updated config")
	return &toReturn, nil
}
