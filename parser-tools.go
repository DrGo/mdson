package mdson

type heading struct {
	name  string
	level int
}

func getHeading(line string) heading {
	i := 0
	for ; i < len(line) && line[i] == '#'; i++ { //no utf8 needed, we are only looking for a byte #
	}
	if i == 0 { //# not found, not a heading
		return heading{name: "", level: 0}
	}
	name := trimLower(line[i:])
	if name == "" { //no name, heading but invalid
		return heading{name: "", level: -1}
	}
	return heading{name: name, level: i}
}
