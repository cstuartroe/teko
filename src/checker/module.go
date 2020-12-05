package checker

var LoadedFiles map[string]*Codeblock = map[string]*Codeblock{}


func LoadFile(filename string) *Codeblock {
  if _, ok := LoadedFiles[filename]; !ok {
    main_codeblock := NewCodeblock(&BaseCodeblock)
    main_codeblock.startFile(filename)

    for main_codeblock.hasMore() {
      main_codeblock.checkNextStatement()
    }

    LoadedFiles[filename] = &main_codeblock
  }

  return LoadedFiles[filename]
}
