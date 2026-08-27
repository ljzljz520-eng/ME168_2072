package catalog

import (
	"fmt"
	"gymrecommend/internal/model"
)

type Reservation struct {
	ClassID  string
	MemberID string
	Sessions int
}

func Reserve(c *Catalog, id string, member string) (model.Class, error) {
	class, ok := c.Get(id)
	if !ok {
		return class, fmt.Errorf("class not found")
	}
	if member == "" {
		return class, fmt.Errorf("member required")
	}
	if !class.IsAvailable() {
		return class, fmt.Errorf("class full")
	}
	class.Enrolled++
	if e := c.Update(class); e != nil {
		return class, e
	}
	return class, nil
}
func Cancel(c *Catalog, id string) error {
	class, ok := c.Get(id)
	if !ok {
		return fmt.Errorf("class not found")
	}
	if class.Enrolled > 0 {
		class.Enrolled--
	}
	return c.Update(class)
}
func CapacityReport(c *Catalog) map[string]float64 {
	out := map[string]float64{}
	for _, class := range c.All() {
		if class.Capacity > 0 {
			out[class.ID] = float64(class.Enrolled) / float64(class.Capacity)
		}
	}
	return out
}
func AvailableForMember(c *Catalog, m model.Member) []model.Class {
	out := []model.Class{}
	for _, class := range c.All() {
		if class.IsAvailable() && model.SlotCompatible(m.PreferredSlot, class.Slot) {
			out = append(out, class)
		}
	}
	return out
}
func NeedWaitlist(class model.Class) bool { return class.Enrolled >= class.Capacity }
