package catalog

import (
	"encoding/json"
	"fmt"
	"gymrecommend/internal/model"
	"io"
)

func Export(w io.Writer, c *Catalog) error { return json.NewEncoder(w).Encode(c.All()) }
func Import(data []byte, c *Catalog) error {
	var classes []model.Class
	if e := json.Unmarshal(data, &classes); e != nil {
		return e
	}
	for _, class := range classes {
		if e := c.Add(class); e != nil {
			return fmt.Errorf("import %s: %w", class.ID, e)
		}
	}
	return nil
}
func Clone(c *Catalog) *Catalog {
	out := New()
	for _, class := range c.All() {
		_ = out.Add(class, c.Tags(class.ID)...)
	}
	return out
}
func Validate(c *Catalog) []error {
	errs := []error{}
	for _, class := range c.All() {
		if e := model.ValidateClass(class); e != nil {
			errs = append(errs, e)
		}
		if class.Enrolled < 0 {
			errs = append(errs, fmt.Errorf("negative enrollment %s", class.ID))
		}
	}
	return errs
}
