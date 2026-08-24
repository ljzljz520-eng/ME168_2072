package catalog

import (
	"fmt"
	"gymrecommend/internal/model"
	"sort"
	"strings"
)

type Catalog struct {
	classes map[string]model.Class
	tags    map[string][]string
}

func New() *Catalog { return &Catalog{classes: map[string]model.Class{}, tags: map[string][]string{}} }
func (c *Catalog) Add(class model.Class, tags ...string) error {
	if e := model.ValidateClass(class); e != nil {
		return e
	}
	if _, ok := c.classes[class.ID]; ok {
		return fmt.Errorf("class already exists")
	}
	c.classes[class.ID] = class
	c.tags[class.ID] = unique(tags)
	return nil
}
func (c *Catalog) Update(class model.Class) error {
	if e := model.ValidateClass(class); e != nil {
		return e
	}
	if _, ok := c.classes[class.ID]; !ok {
		return fmt.Errorf("class not found")
	}
	c.classes[class.ID] = class
	return nil
}
func (c *Catalog) Remove(id string) error {
	if _, ok := c.classes[id]; !ok {
		return fmt.Errorf("class not found")
	}
	delete(c.classes, id)
	delete(c.tags, id)
	return nil
}
func (c *Catalog) Get(id string) (model.Class, bool) { v, ok := c.classes[id]; return v, ok }
func (c *Catalog) All() []model.Class {
	out := make([]model.Class, 0, len(c.classes))
	for _, v := range c.classes {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
func (c *Catalog) Search(query string) []model.Class {
	q := strings.ToLower(strings.TrimSpace(query))
	out := []model.Class{}
	for _, v := range c.classes {
		if q == "" || strings.Contains(strings.ToLower(v.Title), q) || strings.Contains(strings.ToLower(v.Coach), q) || contains(c.tags[v.ID], q) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
func (c *Catalog) BySlot(slot string) []model.Class {
	out := []model.Class{}
	for _, v := range c.classes {
		if v.Slot == slot {
			out = append(out, v)
		}
	}
	return out
}
func (c *Catalog) ByKind(kind string) []model.Class {
	out := []model.Class{}
	for _, v := range c.classes {
		if kind == "" || v.Kind == kind {
			out = append(out, v)
		}
	}
	return out
}
func (c *Catalog) Tag(id string, tag string) error {
	if _, ok := c.classes[id]; !ok {
		return fmt.Errorf("class not found")
	}
	if tag == "" {
		return fmt.Errorf("tag empty")
	}
	c.tags[id] = unique(append(c.tags[id], tag))
	return nil
}
func (c *Catalog) Tags(id string) []string { return append([]string(nil), c.tags[id]...) }
func unique(in []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, v := range in {
		v = strings.TrimSpace(strings.ToLower(v))
		if v != "" && !seen[v] {
			out = append(out, v)
			seen[v] = true
		}
	}
	return out
}
func contains(values []string, q string) bool {
	for _, v := range values {
		if strings.Contains(v, q) {
			return true
		}
	}
	return false
}
