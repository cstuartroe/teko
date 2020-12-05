package checker

var LoadedFiles map[string]*Codeblock = map[string]*Codeblock{}


func LoadFile(filename string) *Codeblock {
  if codeblock, ok := LoadedFiles[filename]; ok {
    return codeblock
  }

  main_codeblock := NewCodeblock(&BaseCodeblock)
  main_codeblock.startFile(filename)

  for main_codeblock.hasMore() {
    main_codeblock.grabStatement()
  }

  return &main_codeblock
}
