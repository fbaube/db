package db

var schemasALL = []string{
	schemaTREF,
}

// schemaTREF is a reference from a Map to a Topic (or subclass thereof).
// Note that this does NOT (necessarily) refer to the DITA `topictref` element!
//
// The relationship is N-to-N btwn Maps and Topics, but `schemaTREF` might
// not be unique because a topic might be explicitly referenced used more
// than once in a map. So for simplicity, let's create only one `schemaTREF`
// per topic per map file, and see if it creates problems elsewhere later on.
var schemaTREF = `TREF
(i_MAP integer not null references FILE,
 i_TPX integer not null references FILE,
 creatime text not null -- ISO-8601
  )`
