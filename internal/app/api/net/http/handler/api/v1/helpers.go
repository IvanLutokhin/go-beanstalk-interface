package v1

import (
	"encoding/base64"
	"fmt"
)

func MarshalJobData(data JobData) ([]byte, error) {
	switch data.Format {
	case JobDataFormatRAW:
		return []byte(data.Payload), nil

	case JobDataFormatBase64:
		return base64.StdEncoding.DecodeString(data.Payload)

	default:
		return nil, fmt.Errorf("job data format '%s' is not supported", data.Format)
	}
}

func UnmarshalJobData(format string, data []byte) (JobData, error) {
	switch format {
	case JobDataFormatRAW:
		return JobData{Format: format, Payload: string(data)}, nil

	case JobDataFormatBase64:
		return JobData{Format: format, Payload: base64.StdEncoding.EncodeToString(data)}, nil

	default:
		return JobData{}, fmt.Errorf("job data format '%s' is not supported", format)
	}
}
