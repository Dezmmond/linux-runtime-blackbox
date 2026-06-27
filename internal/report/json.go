package report

import (
	"encoding/json"
	"io"

	"github.com/yourname/linux-runtime-blackbox/internal/model"
)

func WriteJSON(w io.Writer, r model.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func ReadJSON(r io.Reader) (model.Report, error) {
	var out model.Report
	dec := json.NewDecoder(r)
	if err := dec.Decode(&out); err != nil {
		return model.Report{}, err
	}
	return out, nil
}
