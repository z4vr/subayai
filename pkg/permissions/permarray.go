package permissions

// PermArray is a set of permissions in form of strings.
// Example:
// ["+admin.*", "-mod.kick", "+mod.ban"]
type PermArray []string

// Update updates the permissions of the user. It returns a new PermArray and a bool indicating if permissions were changed.
func (p PermArray) Update(newPerm string, override bool) (newPermsArray PermArray, changed bool) {
	newPermsArray = make(PermArray, len(p)+1)

	i := 0
	add := true
	for _, perm := range p {
		if len(perm) > 0 && perm[1:] == newPerm[1:] {
			add = false

			if override {
				newPermsArray[i] = newPerm
				i++
				continue
			}

			if perm[0] != newPerm[0] {
				continue
			}
		}

		newPermsArray[i] = perm
		i++
	}

	if add {
		newPermsArray[i] = newPerm
		i++
	}

	newPermsArray = newPermsArray[:i]

	changed = !p.Equals(newPermsArray)

	return
}

// Merge updates all entries of p using Update one by one with all entries of newPerms.
// Parameter override is passed to the Update function.
func (p PermArray) Merge(newPerms PermArray, override bool) PermArray {
	for _, cp := range newPerms {
		p, _ = p.Update(cp, override)
	}
	return p
}

// Equals returns true when p2 has the same elements in the same order as p.
func (p PermArray) Equals(p2 PermArray) bool {
	if len(p) != len(p2) {
		return false
	}

	for i, v := range p {
		if v != p2[i] {
			return false
		}
	}

	return true
}

// Check returns true if the passed domainName matches positively on the permission array p.
func (p PermArray) Check(domainName string) bool {
	lvl := -1
	allow := false

	for _, perm := range p {
		m, a := permissionCheckDNs(domainName, perm)
		if m > lvl {
			allow = a
			lvl = m
		}
	}

	return allow
}
