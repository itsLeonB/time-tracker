package mapper

import "github.com/itsLeonB/catfeinated-time-tracker/internal/model"

func YoutrackToExternal(youtrackTask model.YoutrackTask) model.ExternalTask {
	return model.ExternalTask{
		Provider: "YouTrack",
		Number:   youtrackTask.IdReadable,
		Name:     youtrackTask.Summary,
		Project:  youtrackTask.GetEpicName(),
	}
}
