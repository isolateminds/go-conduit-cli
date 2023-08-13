package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/isolateminds/go-conduit-cli/internal/docker/response"
)

type formater struct{ writer io.Writer }

func (f *formater) Write(p []byte) (n int, err error) {
	var data response.ContainerStats
	decoder := json.NewDecoder(bytes.NewReader(p))
	err = decoder.Decode(&data)
	if err != nil {
		return 0, fmt.Errorf("StatsFormatterError:  %s", err)
	}
	b, err := json.Marshal(&response.FormatedContainerStats{
		ID:          data.ID,
		Name:        data.Name,
		CpuUsage:    data.FormatCpuUsagePercentage(),
		MemoryUsage: data.FormatMemoryUsage(),
		NetworkIO:   data.FormatNetworkIO(),
		DiskIO:      data.FormatDiskIO(),
	})
	if err != nil {
		return 0, fmt.Errorf("StatsFormatterError:  %s", err)
	}
	if n, err = f.writer.Write(b); err != nil {
		return n, err
	}
	return len(p), nil
}

// Formats the incoming stats and passes it to the supplied writer
func StatsFormatter(writer io.Writer) *formater {
	return &formater{writer: writer}
}
