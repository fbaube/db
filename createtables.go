package db

// TODO: Change "DESC" to "DESCR" (cos DESC is SQL reserved word)
// TODO: Use fieldnames HED and BOD

// collate nocase

// TODO: ADD FILE ATTRIBUTES and ROOT ELM and OTHER MCFILE STUFF

var schemasALL = []string{
	// schemaINB,
	// schemaFILE,
	schemaTREF,
}

/*
// schemaINB is an Input Batch.
// `creatime` as a text might not implement timezones.
var schemaINB string = `INB(
 idx_inb integer not null primary key, -- NOTE: "integer", not "int"
relfilepath text not null,
absfilepath text not null,
 	 creatime text not null default (datetime('now')), -- UTC ISO-8601
       desc text not null default "No load message", -- from CLI commit msg
     filect int  not null check (filect >= 0) default 0
  )`

// schemaFILE is used for any discrete chunk of content. We don't want to
// allow nulls, cos we get extra error checking that way, so for filepaths
// (`*FP`), allow the distinguished strings `stdin` and `other`.
var schemaFILE string = `FILE
(idx_file integer not null primary key,
  idx_inb integer references INB, -- can be null if stdin/other
   creatime text not null default (datetime('now')), -- UTC ISO-8601
    contype text not null, -- "MAP", "TPX", etc.
    rootelm text not null, -- can be ""
   mimetype text not null,
    doctype text not null,
      mtype text not null,
relfilepath text not null, -- rel.to INB.fullFP, so rel.to all others in InBatch
absfilepath text not null,
        hed text not null, -- Header: Metadata; TODO: K-V store in JSON ?
				bod text not null, -- Body: the Content
FOREIGN KEY(idx_inb) REFERENCES INB(idx_inb)
 )`
*/

// schemaTREF is a reference from a Map to a Topic (or subclass thereof).
// Note that this does NOT (necessarily) refer to the DITA `topictref` element!
//
// The relationship is N-to-N btwn Maps and Topics, but `schemaTREF` might
// not be unique because a topic might be explicitly referenced used more
// than once in a map. So for simplicity, let's create only one `schemaTREF`
// per topic per map file, and see if it creates problems elsewhere later on.
var schemaTREF string = `TREF
(i_MAP integer not null references FILE,
 i_TPX integer not null references FILE,
 creatime text not null -- ISO-8601
  )`

/*

A LinkInfos has
xmlIDs
xmlIDrefs
Conrefs
Keyrefs
Datarefs

An XmlItems has
- IDs & IDREFs
- ENTITY & ELEMENT directives

*/
