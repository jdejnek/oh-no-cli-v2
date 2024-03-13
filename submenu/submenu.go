package submenu

type model struct {
	options  []string
	cursor   int
	selected int
}

func submenu() model {
	var selected int
	return model{
		options:  []string{"View", "Create", "Delete"},
		selected: selected,
	}
}
