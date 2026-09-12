package rank

type Rank struct {
	ID          int
	Chat        string
	Display     string
	Priority    int
	Permissions map[string]bool
}

func (r *Rank) HasPermission(perm string) bool {
	if r == nil {
		return false
	}

	return r.Permissions[perm]
}
