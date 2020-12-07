package checker

import (
    "github.com/cstuartroe/teko/src/lexparse"
)

type TypeModule struct {
  codeblock lexparse.Codeblock
  checker Checker
}

var LoadedFiles map[string]*TypeModule = map[string]*TypeModule{}


func LoadFile(filename string) *TypeModule {
  if _, ok := LoadedFiles[filename]; !ok {
    main_codeblock := lexparse.ParseFile(filename)
    main_checker := NewChecker(&BaseChecker)

    for _, stmt := range main_codeblock.GetStatements() {
      lexparse.PrintNode(stmt)
      main_checker.checkStatement(stmt)
    }

    main_module := TypeModule{
      codeblock: main_codeblock,
      checker: main_checker,
    }

    LoadedFiles[filename] = &main_module
  }

  return LoadedFiles[filename]
}
