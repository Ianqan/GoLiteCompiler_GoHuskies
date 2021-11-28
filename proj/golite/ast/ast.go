package ast

import (
	"bytes"
	"fmt"
	"proj/golite/ir"
	st "proj/golite/symboltable"
	"proj/golite/token"
	"proj/golite/types"
)

// Node The base Node interface that all ast nodes have to access
type Node interface {
	TokenLiteral() string
	String() string
	TypeCheck([]string, *st.SymbolTable) []string // TO-DO
	TranslateToILoc([]ir.Instruction, *st.SymbolTable) []ir.Instruction
}

// Expr All expression nodes implement this interface
type Expr interface {
	Node
	GetType(*st.SymbolTable) types.Type // TO-DO
}

// Stmt All statement nodes implement this interface
type Stmt interface {
	Node
	PerformSABuild([]string, *st.SymbolTable) []string // TO-DO
}

/******* Stmt : Statement *******/

type Program struct {
	Token *token.Token
	st    *st.SymbolTable

	Package      *Package
	Import       *Import
	Types        *Types
	Declarations *Declarations
	Functions    *Functions
}

func (p *Program) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literal for program statement.")
}
func (p *Program) String() string {
	out := bytes.Buffer{}
	out.WriteString(p.Package.String())
	out.WriteString(p.Import.String())
	out.WriteString(p.Types.String())
	out.WriteString(p.Declarations.String())
	out.WriteString(p.Functions.String())
	return out.String()
}
func (p *Program) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	p.st = symTable
	errors = p.Package.PerformSABuild(errors, symTable)
	errors = p.Import.PerformSABuild(errors, symTable)
	errors = p.Types.PerformSABuild(errors, symTable)
	errors = p.Declarations.PerformSABuild(errors, symTable)
	errors = p.Functions.PerformSABuild(errors, symTable)
	return errors
}
func (p *Program) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = p.Package.TypeCheck(errors, symTable)
	errors = p.Import.TypeCheck(errors, symTable)
	errors = p.Types.TypeCheck(errors, symTable)
	errors = p.Declarations.TypeCheck(errors, symTable)
	errors = p.Functions.TypeCheck(errors, symTable)
	return errors
}
func (p *Program) TranslateToILoc(symTable *st.SymbolTable) []ir.Instruction {
	//instr := nil
	//return instr
	return nil
}

type Package struct {
	Token *token.Token
	//st    *st.SymbolTable
	Ident IdentLiteral
}

func (pkg *Package) TokenLiteral() string {
	if pkg.Token != nil {
		return pkg.Token.Literal
	}
	panic("Could not determine token literal for package statement")
}
func (pkg *Package) String() string {
	out := bytes.Buffer{}
	out.WriteString("package")
	out.WriteString(" ")
	out.WriteString(pkg.Ident.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}
func (pkg *Package) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}
func (pkg *Package) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	if pkg.Ident.TokenLiteral() != "main" {
		errors = append(errors, fmt.Sprintf("Only package main is allowed"))
	}
	return errors
}

type Import struct {
	Token *token.Token
	Ident IdentLiteral
}

func (imp *Import) TokenLiteral() string {
	if imp.Token != nil {
		return imp.Token.Literal
	}
	panic("Could not determine token literal for import statement")
}
func (imp *Import) String() string {
	out := bytes.Buffer{}
	out.WriteString("import")
	out.WriteString(" ")
	out.WriteString("\"")
	out.WriteString("fmt") // imp.Ident.String() equivalent
	out.WriteString("\"")
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}
func (imp *Import) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	return errors
}
func (imp *Import) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type Types struct {
	Token *token.Token
	//st    *st.SymbolTable
	TypeDeclarations []TypeDeclaration
}

func (tys *Types) TokenLiteral() string {
	if tys.Token != nil {
		return tys.Token.Literal
	}
	panic("Could not determine token literals for the types declarations")
}
func (tys *Types) String() string {
	out := bytes.Buffer{}
	for _, typedec := range tys.TypeDeclarations {
		out.WriteString(typedec.String())
		out.WriteString("\n")
	}
	return out.String()
}
func (tys *Types) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	for _, typedec := range tys.TypeDeclarations {
		errors = typedec.PerformSABuild(errors, symTable)
	}
	return errors
}
func (tys *Types) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, typedec := range tys.TypeDeclarations {
		errors = typedec.TypeCheck(errors, symTable)
	}
	return errors
}

type TypeDeclaration struct {
	Token  *token.Token
	st     *st.SymbolTable
	Ident  IdentLiteral
	Fields *Fields
}

func (td *TypeDeclaration) TokenLiteral() string {
	if td.Token != nil {
		return td.Token.Literal
	}
	panic("Could not determine token literals for the type declaration")
}
func (td *TypeDeclaration) String() string {
	out := bytes.Buffer{}
	out.WriteString("type")
	out.WriteString(" ")
	out.WriteString(td.Ident.String())
	out.WriteString(" ")
	out.WriteString("struct")
	out.WriteString("{\n")
	out.WriteString(td.Fields.String())
	out.WriteString("\n}")
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}
func (td *TypeDeclaration) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: find duplicate structures
	structName := td.Ident.TokenLiteral()
	scopeSymTable := st.New(symTable, structName)
	td.st = scopeSymTable

	if entry := symTable.Contains(structName); entry != nil {
		errors = append(errors, fmt.Sprintf("[%v]: struct %v already declared", td.Token.LineNum, structName))
	} else {
		var entry st.Entry
		entry = st.NewStructEntry(td.st)
		symTable.Insert(structName, &entry)
		errors = td.Fields.PerformSABuild(errors, td.st)
	}
	return errors
}
func (td *TypeDeclaration) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors2 := td.Fields.TypeCheck(errors, td.st)
	errors = append(errors, errors2...)
	return errors
}

type Fields struct {
	Token *token.Token
	Decls []Decl
}

func (fields *Fields) TokenLiteral() string {
	if fields.Token != nil {
		return fields.Token.Literal
	}
	panic("Could not determine token literals for fields")
}
func (fields *Fields) String() string {
	out := bytes.Buffer{}
	out.WriteString(fields.Decls[0].String())
	out.WriteString(";\n")
	remaining := fields.Decls[1:]
	for _, decl := range remaining {
		out.WriteString(decl.String())
		out.WriteString(";\n")
	}
	return out.String()
}
func (fields *Fields) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, decl := range fields.Decls {
		errors = decl.PerformSABuild(errors, symTable)
	}
	return errors
}
func (fields *Fields) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, decl := range fields.Decls {
		errors = decl.TypeCheck(errors, symTable)
	}
	return errors
}

type Decl struct {
	Token *token.Token
	Ident IdentLiteral
	Ty    *Type
}

func (decl *Decl) TokenLiteral() string {
	if decl.Token != nil {
		return decl.Token.Literal
	}
	panic("Could not determine token literals for decl")
}
func (decl *Decl) String() string {
	out := bytes.Buffer{}
	out.WriteString(decl.Ident.String())
	out.WriteString(" ")
	out.WriteString(decl.Ty.String())
	return out.String()
}
func (decl *Decl) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: find duplicate declarations in functions / structures
	varName := decl.Ident.TokenLiteral()
	if entry := symTable.Contains(varName); entry != nil {
		errors = append(errors, fmt.Sprintf("[%v]: variable %v already declared", decl.Token.LineNum, varName))
	} else {
		var entry st.Entry
		entry = st.NewVarEntry()
		symTable.Insert(varName, &entry)
	}
	return errors
}
func (decl *Decl) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: set type for the variable id

	// Decl = 'id' Type
	// get the type from Type
	varType := decl.Ty.GetType(symTable)
	// update / set type of 'id' in the symbol table
	entry := symTable.Contains(decl.Ident.TokenLiteral())
	// entry must be valid when processing declaration-type statements, because it must have been added to symboltable by PerformSABuild
	entry.SetType(varType)
	// include id and its type as a parameter of the function
	// can be function or struct, but only useful in function
	symTable.ScopeParamTys = append(symTable.ScopeParamTys, varType)
	symTable.ScopeParamNames = append(symTable.ScopeParamNames, decl.Ident.TokenLiteral())
	return errors
}

type Declarations struct {
	Token        *token.Token
	Declarations []Declaration
}

func (ds *Declarations) TokenLiteral() string {
	if ds.Token != nil {
		return ds.Token.Literal
	}
	panic("Could not determine token literals for the declarations")
}
func (ds *Declarations) String() string {
	out := bytes.Buffer{}
	for _, dec := range ds.Declarations {
		out.WriteString(dec.String())
		out.WriteString("\n")
	}
	return out.String()
}
func (ds *Declarations) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, dec := range ds.Declarations {
		errors = dec.PerformSABuild(errors, symTable)
	}
	return errors
}
func (ds *Declarations) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, dec := range ds.Declarations {
		errors = dec.TypeCheck(errors, symTable)
	}
	return errors
}

type Declaration struct {
	Token *token.Token
	Ids   *Ids
	Ty    *Type
}

func (d *Declaration) TokenLiteral() string {
	if d.Token != nil {
		return d.Token.Literal
	}
	panic("Could not determine token literals for declaration")
}
func (d *Declaration) String() string {
	out := bytes.Buffer{}
	out.WriteString("var")
	out.WriteString(" ")
	out.WriteString(d.Ids.String())
	out.WriteString(" ")
	out.WriteString(d.Ty.String())
	out.WriteString(";")
	return out.String()
}
func (d *Declaration) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none, duplicate definitions are examined in d.Ids.PerformSABuild()
	errors = d.Ids.PerformSABuild(errors, symTable)
	return errors
}
func (d *Declaration) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: set type for ids, symbol table only
	decType := d.Ty.GetType(symTable)
	for _, id := range d.Ids.Idents {
		entry := symTable.Contains(id.TokenLiteral())
		entry.SetType(decType)
	}
	return errors
}

type Ids struct {
	Token  *token.Token
	Idents []IdentLiteral
}

func (ids *Ids) TokenLiteral() string {
	if ids.Token != nil {
		return ids.Token.Literal
	}
	panic("Could not determine token literals for ids")
}
func (ids *Ids) String() string {
	out := bytes.Buffer{}
	out.WriteString(ids.Idents[0].String())
	remaining := ids.Idents[1:]
	for _, id := range remaining {
		out.WriteString(",")
		out.WriteString(id.String())
	}
	return out.String()
}
func (ids *Ids) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// Objective: find duplicate declarations
	for _, id := range ids.Idents {
		varName := id.TokenLiteral()
		if entry := symTable.Contains(varName); entry != nil {
			errors = append(errors, fmt.Sprintf("[%v]: variable [%v] already declared", id.Token.LineNum, varName))
		} else {
			var entry st.Entry
			entry = st.NewVarEntry()
			symTable.Insert(varName, &entry)
		}
	}
	return errors
}
func (ids *Ids) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none, accomplished in Declaration
	return errors
}

type Functions struct {
	Token     *token.Token
	Functions []Function
}

func (fs *Functions) TokenLiteral() string {
	if fs.Token != nil {
		return fs.Token.Literal
	}
	panic("Could not determine token literals for the functions")
}
func (fs *Functions) String() string {
	out := bytes.Buffer{}
	for _, fun := range fs.Functions {
		out.WriteString(fun.String())
	}
	return out.String()
}
func (fs *Functions) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, fun := range fs.Functions {
		errors = fun.PerformSABuild(errors, symTable)
	}
	return errors
}
func (fs *Functions) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	for _, fun := range fs.Functions {
		errors = fun.TypeCheck(errors, symTable)
	}
	return errors
}

type Function struct {
	Token        *token.Token
	st           *st.SymbolTable
	Ident        IdentLiteral
	Parameters   *Parameters
	ReturnType   *ReturnType
	Declarations *Declarations
	Statements   *Statements
}

func (f *Function) TokenLiteral() string {
	if f.Token != nil {
		return f.Token.Literal
	}
	panic("Could not determine token literals for functions")
}
func (f *Function) String() string {
	out := bytes.Buffer{}
	out.WriteString("func")
	out.WriteString(" ")
	out.WriteString(f.Ident.String())
	out.WriteString(" ")
	out.WriteString(f.Parameters.String())
	out.WriteString(" ")
	out.WriteString(f.ReturnType.String())
	out.WriteString("{")
	out.WriteString("\n")
	out.WriteString(f.Declarations.String())
	out.WriteString(" ")
	out.WriteString(f.Statements.String())
	out.WriteString("}")
	out.WriteString("\n")
	return out.String()
}
func (f *Function) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// Objective: find duplicate function definitions
	funcName := f.Ident.TokenLiteral()
	scopeSymTable := st.New(symTable, funcName)
	f.st = scopeSymTable
	if entry := symTable.Contains(funcName); entry != nil {
		errors = append(errors, fmt.Sprintf("[%v]: function [%v] has been declared", f.Token.LineNum, funcName))
	} else {
		var entry st.Entry
		entry = st.NewFuncEntry(f.ReturnType.GetType(symTable), f.st)
		symTable.Insert(funcName, &entry)
		errors = f.Parameters.PerformSABuild(errors, f.st)
		errors = f.Declarations.PerformSABuild(errors, f.st)
		errors = f.Statements.PerformSABuild(errors, f.st)
	}
	return errors
}
func (f *Function) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// Objective: add parameters, return type to function symbol table and entry in the outer symbol table
	// parameters are added to both inner symbol table and function signature in the outer symbol table by Decl invoked next line
	currScopeSt := symTable.Contains(f.Ident.TokenLiteral()).GetScopeST()
	f.st = currScopeSt
	errors = f.Parameters.TypeCheck(errors, f.st)
	errors = f.Declarations.TypeCheck(errors, f.st)
	errors = f.Statements.TypeCheck(errors, f.st)
	return errors
}

type Parameters struct {
	Token *token.Token
	Decls []Decl
}

func (params *Parameters) TokenLiteral() string {
	if params.Token != nil {
		return params.Token.Literal
	}
	panic("Could not determine token literals for parameters")
}
func (params *Parameters) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	var remaining []Decl
	if len(params.Decls) > 0 {
		out.WriteString(params.Decls[0].String())
		remaining = params.Decls[1:]
	}
	for _, decl := range remaining {
		out.WriteString(",")
		out.WriteString(decl.String())
	}
	out.WriteString(")")
	return out.String()
}
func (params *Parameters) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, decl := range params.Decls {
		errors = decl.PerformSABuild(errors, symTable)
	}
	return errors
}
func (params *Parameters) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, decl := range params.Decls {
		errors = decl.TypeCheck(errors, symTable)
	}
	return errors
}

type Statements struct {
	Token      *token.Token
	Statements []Statement
}

func (stmts *Statements) TokenLiteral() string {
	if stmts.Token != nil {
		return stmts.Token.Literal
	}
	panic("Could not determine token literals for the statements")
}
func (stmts *Statements) String() string {
	out := bytes.Buffer{}
	for _, stmt := range stmts.Statements {
		out.WriteString("\t")
		out.WriteString(stmt.String())
	}
	out.WriteString("\n")
	return out.String()
}
func (stmts *Statements) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, stmt := range stmts.Statements {
		errors = stmt.PerformSABuild(errors, symTable)
	}
	return errors
}
func (stmts *Statements) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	for _, stmt := range stmts.Statements {
		errors = stmt.TypeCheck(errors, symTable)
	}
	return errors
}

type Statement struct {
	Token *token.Token
	Stmt  Stmt
}

func (s *Statement) TokenLiteral() string {
	if s.Token != nil {
		return s.Token.Literal
	}
	panic("Could not determine token literals for statement")
}
func (s *Statement) String() string {
	out := bytes.Buffer{}
	out.WriteString(s.Stmt.String())
	return out.String()
}
func (s *Statement) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = s.Stmt.PerformSABuild(errors, symTable)
	return errors
}
func (s *Statement) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = s.Stmt.TypeCheck(errors, symTable)
	return errors
}

type Block struct {
	Token      *token.Token
	Statements *Statements
}

func (b *Block) TokenLiteral() string {
	if b.Token != nil {
		return b.Token.Literal
	}
	panic("Could not determine token literals for block")
}
func (b *Block) String() string {
	out := bytes.Buffer{}
	out.WriteString("{")
	out.WriteString("\n")
	out.WriteString(b.Statements.String())
	out.WriteString("}")
	return out.String()
}
func (b *Block) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = b.Statements.PerformSABuild(errors, symTable)
	return errors
}
func (b *Block) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = b.Statements.TypeCheck(errors, symTable)
	return errors
}

type Assignment struct {
	Token  *token.Token
	Lvalue *LValue
	Expr   *Expression
}

func (a *Assignment) TokenLiteral() string {
	if a.Token != nil {
		return a.Token.Literal
	}
	panic("Could not determine token literals for assignment")
}
func (a *Assignment) String() string {
	out := bytes.Buffer{}
	out.WriteString(a.Lvalue.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(a.Expr.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}
func (a *Assignment) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	return errors
}
func (a *Assignment) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: matching of types on both sides of the assignment statement
	errors = a.Lvalue.TypeCheck(errors, symTable)
	errors = a.Expr.TypeCheck(errors, symTable)
	if len(errors) == 0 {
		leftType := a.Lvalue.GetType(symTable)
		rightType := a.Expr.GetType(symTable)
		if leftType != rightType {
			errors = append(errors, fmt.Sprintf("[%v]: type mismatch: Cannot assign %v (Type %v) to %v (Type %v)",
				a.Token.LineNum, a.Expr.String(), rightType.GetName(), a.Lvalue.String(), leftType.GetName()))
			return errors
		}
		if leftType != types.IntTySig && leftType != types.BoolTySig {
			errors = append(errors, fmt.Sprintf("[%v]: %v is not assignable", a.Token.LineNum, a.Lvalue.String()))
			return errors
		}
		if leftType == types.IntTySig && a.Lvalue.Token.Type == token.NUM {
			errors = append(errors, fmt.Sprintf("[%v]: %v is not assignable", a.Token.LineNum, a.Lvalue.String()))
			return errors
		}
		if leftType == types.BoolTySig && (a.Lvalue.Token.Type == token.TRUE || a.Lvalue.Token.Type == token.FALSE) {
			errors = append(errors, fmt.Sprintf("[%v]: %v is not assignable", a.Token.LineNum, a.Lvalue.String()))
			return errors
		}
	}
	return errors
}

type Read struct {
	Token *token.Token
	Ident IdentLiteral
}

func (r *Read) TokenLiteral() string {
	if r.Token != nil {
		return r.Token.Literal
	}
	panic("Could not determine token literals for read")
}
func (r *Read) String() string {
	out := bytes.Buffer{}
	out.WriteString("fmt")
	out.WriteString(".")
	out.WriteString("Scan")
	out.WriteString("(")
	out.WriteString("&")
	out.WriteString(r.Ident.String())
	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}
func (r *Read) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	return errors
}
func (r *Read) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: verify the variable is declared
	varName := r.Ident.TokenLiteral()
	entry := symTable.Contains(varName)
	if entry == nil {
		errors = append(errors, fmt.Sprintf("[%v]: variable %v has not been declared", r.Token.LineNum, varName))
	}
	return errors
}

type Print struct {
	Token       *token.Token
	printMethod string // "Print" | "Println"
	Ident       IdentLiteral
}

func (p *Print) TokenLiteral() string {
	if p.Token != nil {
		return p.Token.Literal
	}
	panic("Could not determine token literals for print")
}
func (p *Print) String() string {
	out := bytes.Buffer{}
	out.WriteString("fmt")
	out.WriteString(".")
	out.WriteString(p.printMethod)
	out.WriteString("(")
	out.WriteString(p.Ident.String())
	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}
func (p *Print) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	return errors
}
func (p *Print) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: verify the variable is declared
	varName := p.Ident.TokenLiteral()
	entry := symTable.Contains(varName)
	if entry == nil {
		errors = append(errors, fmt.Sprintf("[%v]: variable %v has not been declared", p.Token.LineNum, varName))
	}
	return errors
}

type Conditional struct {
	Token     *token.Token
	Expr      *Expression
	Block     *Block
	ElseBlock *Block
}

func (cond *Conditional) TokenLiteral() string {
	if cond.Token != nil {
		return cond.Token.Literal
	}
	panic("Could not determine token literals for conditional")
}
func (cond *Conditional) String() string {
	out := bytes.Buffer{}
	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString(cond.Expr.String())
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(cond.Block.String())
	if cond.ElseBlock != nil {
		out.WriteString("else")
		out.WriteString(cond.ElseBlock.String())
	}
	return out.String()
}
func (cond *Conditional) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = cond.Block.PerformSABuild(errors, symTable)
	errors = cond.ElseBlock.PerformSABuild(errors, symTable)
	return errors
}
func (cond *Conditional) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: boolean expression as the conditional expression surrounded by parenthesis
	condType := cond.Expr.GetType(symTable)
	errors = cond.Expr.TypeCheck(errors, symTable)
	errors = cond.Block.TypeCheck(errors, symTable)
	errors = cond.ElseBlock.TypeCheck(errors, symTable)
	if len(errors) == 0 {
		if condType != types.BoolTySig {
			errors = append(errors, fmt.Sprintf("[%v]: boolean expression is desired, received %v Type %v", cond.Token.LineNum, cond.Expr.String(), condType.GetName()))
		}
	}
	return errors
}

type Loop struct {
	Token *token.Token
	Expr  *Expression
	Block *Block
}

func (lp *Loop) TokenLiteral() string {
	if lp.Token != nil {
		return lp.Token.Literal
	}
	panic("Could not determine token literals for loop")
}
func (lp *Loop) String() string {
	out := bytes.Buffer{}
	out.WriteString("for")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString(lp.Expr.String())
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(lp.Block.String())
	return out.String()
}
func (lp *Loop) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	errors = lp.Block.PerformSABuild(errors, symTable)
	return errors
}
func (lp *Loop) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: boolean expression as the conditional expression surrounded by parenthesis
	condType := lp.Expr.GetType(symTable)
	errors = lp.Expr.TypeCheck(errors, symTable)
	errors = lp.Block.TypeCheck(errors, symTable)
	if len(errors) == 0 {
		if condType != types.BoolTySig {
			errors = append(errors, fmt.Sprintf("[%v]: boolean expression is desired, received %v Type %v",
				lp.Token.LineNum, lp.Expr.String(), condType.GetName()))
		}
	}
	return errors
}

type Return struct {
	Token *token.Token // "RETURN"
	Expr  *Expression  // the return type, nil if not exists
}

func (ret *Return) TokenLiteral() string {
	if ret.Token != nil {
		return ret.Token.Literal
	}
	panic("Could not determine token literals for return")
}
func (ret *Return) String() string {
	out := bytes.Buffer{}
	out.WriteString("return")
	if ret.Expr != nil {
		out.WriteString(" ")
		out.WriteString(ret.Expr.String())
	}
	out.WriteString(";")
	return out.String()
}
func (ret *Return) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	return errors
}
func (ret *Return) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// objective: match return type with signature
	errors = ret.Expr.TypeCheck(errors, symTable)
	// go to outer symbol table and retrieve the entry
	funcEntry := symTable.Parent.Contains(symTable.ScopeName) // must exist
	decRetType := funcEntry.GetReturnTy()                     // must exist
	actRetType := ret.Expr.GetType(symTable)
	if len(errors) == 0 {
		if actRetType != decRetType {
			errors = append(errors, fmt.Sprintf("[%v]: return type expected %v, found %v", ret.Token.LineNum, decRetType, actRetType))
		}
	}
	return errors
}

// Invocation Statement, compared with InvocExpr
type Invocation struct {
	Token *token.Token
	Ident IdentLiteral
	Args  *Arguments
}

func (invoc *Invocation) TokenLiteral() string {
	if invoc.Token != nil {
		return invoc.Token.Literal
	}
	panic("Could not determine token literals for invocation statement")
}
func (invoc *Invocation) String() string {
	out := bytes.Buffer{}
	out.WriteString(invoc.Ident.String())
	out.WriteString(" ")
	out.WriteString(invoc.Args.String())
	out.WriteString(";")
	return out.String()
}
func (invoc *Invocation) PerformSABuild(errors []string, symTable *st.SymbolTable) []string {
	// objective: none
	return errors
}
func (invoc *Invocation) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// check whether function is declared
	funcName := invoc.Ident.TokenLiteral()
	entry := symTable.Contains(funcName)
	if entry == nil {
		errors = append(errors, fmt.Sprintf("[%v]: function %v has not been defined", invoc.Token.LineNum, funcName))
	} else {
		symTable = entry.GetScopeST()
		errors = invoc.Args.TypeCheck(errors, symTable)
	}
	return errors
}

func NewProgram(pac *Package, imp *Import, typ *Types, decs *Declarations, funs *Functions) *Program {
	return &Program{nil, nil, pac, imp, typ, decs, funs}
}
func NewPackage(ident IdentLiteral) *Package {
	return &Package{nil, ident}
}
func NewImport(ident IdentLiteral) *Import      { return &Import{nil, ident} }
func NewTypes(typdecs []TypeDeclaration) *Types { return &Types{nil, typdecs} }
func NewTypeDeclaration(ident IdentLiteral, fields *Fields) *TypeDeclaration {
	return &TypeDeclaration{nil, nil, ident, fields}
}
func NewFields(decls []Decl) *Fields                   { return &Fields{nil, decls} }
func NewDecl(ident IdentLiteral, ty *Type) *Decl       { return &Decl{nil, ident, ty} }
func NewDeclarations(decs []Declaration) *Declarations { return &Declarations{nil, decs} }
func NewDeclaration(ids *Ids, Type *Type) *Declaration { return &Declaration{nil, ids, Type} }
func NewIds(idents []IdentLiteral) *Ids                { return &Ids{nil, idents} }
func NewFunctions(funs []Function) *Functions          { return &Functions{nil, funs} }
func NewFunction(ident IdentLiteral, params *Parameters, returnType *ReturnType,
	declarations *Declarations, statements *Statements) *Function {
	return &Function{nil, nil, ident, params, returnType, declarations, statements}
}
func NewParameters(decls []Decl) *Parameters      { return &Parameters{nil, decls} }
func NewReturnType(str string) *ReturnType        { return &ReturnType{nil, NewType(str)} }
func NewStatements(stmts []Statement) *Statements { return &Statements{nil, stmts} }
func NewStatement(stmt Stmt) *Statement           { return &Statement{nil, stmt} }
func NewBlock(statement *Statements) *Block       { return &Block{nil, statement} }
func NewAssignment(lvalue *LValue, expr *Expression) *Assignment {
	return &Assignment{nil, lvalue, expr}
}
func NewRead(ident IdentLiteral) *Read { return &Read{nil, ident} }
func NewPrint(printMethod string, ident IdentLiteral) *Print {
	return &Print{nil, printMethod, ident}
}
func NewConditional(expr *Expression, block *Block, elseBlock *Block) *Conditional {
	return &Conditional{nil, expr, block, elseBlock}
}
func NewLoop(expr *Expression, block *Block) *Loop { return &Loop{nil, expr, block} }
func NewReturn(expr *Expression) *Return           { return &Return{nil, expr} }
func NewInvocation(ident IdentLiteral, args *Arguments) *Invocation {
	return &Invocation{nil, ident, args}
}

/***************** Expr : Expression *******************/

type Type struct {
	Token *token.Token
	// either "int"/"bool"/"*id", where id will actually be the literal for the struct name being defined.
	TypeLiteral string
}

func (t *Type) TokenLiteral() string {
	if t.Token != nil {
		return t.Token.Literal
	}
	panic("Could not determine token literals for type")
}
func (t *Type) String() string {
	return t.TypeLiteral
}
func (t *Type) GetType(symTable *st.SymbolTable) types.Type {
	if t.TypeLiteral == "int" {
		return types.IntTySig
	} else if t.TypeLiteral == "bool" {
		return types.BoolTySig
	} else {
		return types.StructTySig
	}
}
func (t *Type) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

type ReturnType struct {
	Token *token.Token
	Ty    *Type
}

func (rt *ReturnType) TokenLiteral() string {
	if rt.Token != nil {
		return rt.Token.Literal
	}
	panic("Could not determine token literals for returnType")
}
func (rt *ReturnType) String() string {
	out := bytes.Buffer{}
	out.WriteString(rt.Ty.String())
	return out.String()
}
func (rt *ReturnType) GetType(symTable *st.SymbolTable) types.Type {
	if rt.Ty.TypeLiteral == "" {
		return types.VoidTySig
	} else {
		return rt.Ty.GetType(symTable)
	}
}
func (rt *ReturnType) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = rt.Ty.TypeCheck(errors, symTable)
	return errors
}

type Arguments struct {
	Token *token.Token
	Exprs []Expression // MARKING
}

func (args *Arguments) TokenLiteral() string {
	if args.Token != nil {
		return args.Token.Literal
	}
	panic("Could not determine token literals for arguments")
}
func (args *Arguments) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	if len(args.Exprs) > 0 {
		out.WriteString(args.Exprs[0].String())
		remaining := args.Exprs[1:]
		for _, exp := range remaining {
			out.WriteString(",")
			out.WriteString(exp.String())
		}
	}
	out.WriteString(")")
	return out.String()
}
func (args *Arguments) GetType(symTable *st.SymbolTable) types.Type {
	return types.VoidTySig
}
func (args *Arguments) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// used as parameters for calling a function
	expectedTys := symTable.ScopeParamTys
	paramNames := symTable.ScopeParamNames
	for idx, expr := range args.Exprs {
		errors = expr.TypeCheck(errors, symTable)
		givenParamTy := expr.GetType(symTable)
		if givenParamTy != expectedTys[idx] {
			errors = append(errors, fmt.Sprintf("[%v]: Expected paramter %v type %v; given parameter %v type %v #{args.Exprs[idx]}",
				args.Token.LineNum, paramNames[idx], expectedTys[idx], expr.String(), givenParamTy))
		}
	}
	return errors
}

type LValue struct {
	Token  *token.Token
	Ident  IdentLiteral
	Idents []IdentLiteral
}

func (lv *LValue) TokenLiteral() string {
	if lv.Token != nil {
		return lv.Token.Literal
	}
	panic("Could not determine token literals for lvalue")
}
func (lv *LValue) String() string {
	out := bytes.Buffer{}
	out.WriteString(lv.Ident.String())
	for _, id := range lv.Idents {
		out.WriteString(".")
		out.WriteString(id.String())
	}
	return out.String()
}
func (lv *LValue) GetType(symTable *st.SymbolTable) types.Type {
	var entry st.Entry
	varName := lv.Ident.TokenLiteral()
	for {
		if entry = symTable.Contains(varName); entry == nil {
			if symTable.Parent == nil {
				return types.UnknownTySig
			} else {
				symTable = symTable.Parent
			}
		} else {
			break
		}
	}
	if lv.Idents == nil {
		return entry.GetEntryType()
	}
	// here entry is the entry of the first id in Idents
	symTable = entry.GetScopeST()
	remaining := lv.Idents[1:]
	for idx, id := range remaining {
		if entry = symTable.Contains(id.String()); entry == nil {
			return types.UnknownTySig
		} else {
			if idx == len(lv.Idents)-1 {
				break
			}
			symTable = entry.GetScopeST()
		}
	}
	return entry.GetEntryType()
}
func (lv *LValue) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	if lv.GetType(symTable) == types.UnknownTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (LValue) inner field has not been defined", lv.Token.LineNum))
	}
	return errors
}

type Expression struct {
	Token  *token.Token
	Left   *BoolTerm
	Rights []BoolTerm
}

func (exp *Expression) TokenLiteral() string {
	if exp.Token != nil {
		return exp.Token.Literal
	}
	panic("Could not determine token literals for expression")
}
func (exp *Expression) String() string {
	out := bytes.Buffer{}
	out.WriteString(exp.Left.String())
	for _, boolTerm := range exp.Rights {
		out.WriteString("||")
		out.WriteString(boolTerm.String())
	}
	return out.String()
}
func (exp *Expression) GetType(symTable *st.SymbolTable) types.Type {
	leftType := exp.Left.GetType(symTable)

	for _, rTerm := range exp.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (exp *Expression) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = exp.Left.TypeCheck(errors, symTable)
	leftMostTy := exp.Left.GetType(symTable)
	if len(exp.Rights) != 0 {
		// OR operation, needs bool types on both sides
		if leftMostTy != types.BoolTySig {
			errors = append(errors, fmt.Sprintf("[%v]: (Expression) expected bool type, found %v (%v)",
				exp.Token.LineNum, leftMostTy.GetName(), exp.Left.String()))
			return errors
		}
	}
	// check every expression is in the same type (types.BoolTySig)
	for _, curr := range exp.Rights {
		errors = curr.TypeCheck(errors, symTable)
		currTy := curr.GetType(symTable)
		if currTy != leftMostTy {
			errors = append(errors, fmt.Sprintf("[%v]: (Expression) expected %v type, found %v (%v)",
				exp.Token.LineNum, leftMostTy.GetName(), currTy.GetName(), curr.String()))
			break
		}
	}
	return errors
}

type BoolTerm struct {
	Token  *token.Token
	Left   *EqualTerm
	Rights []EqualTerm
}

func (bt *BoolTerm) TokenLiteral() string {
	if bt.Token != nil {
		return bt.Token.Literal
	}
	panic("Could not determine token literals for boolTerm")
}
func (bt *BoolTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(bt.Left.String())
	for _, equalTerm := range bt.Rights {
		out.WriteString("&&")
		out.WriteString(equalTerm.String())
	}
	return out.String()
}
func (bt *BoolTerm) GetType(symTable *st.SymbolTable) types.Type {
	leftType := bt.Left.GetType(symTable)

	for _, rTerm := range bt.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (bt *BoolTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = bt.Left.TypeCheck(errors, symTable)
	leftMostTy := bt.Left.GetType(symTable)
	if len(bt.Rights) != 0 {
		// OR operation, needs bool types on both sides
		if leftMostTy != types.BoolTySig {
			errors = append(errors, fmt.Sprintf("[%v]: (BoolTerm) expected bool type, found %v (%v)",
				bt.Token.LineNum, leftMostTy.GetName(), bt.Left.String()))
			return errors
		}
	}
	// check every expression is in the same type (types.BoolTySig)
	for _, curr := range bt.Rights {
		errors = curr.TypeCheck(errors, symTable)
		currTy := curr.GetType(symTable)
		if currTy != leftMostTy {
			errors = append(errors, fmt.Sprintf("[%v]: (BoolTerm) expected %v type, found %v (%v)",
				bt.Token.LineNum, leftMostTy.GetName(), currTy.GetName(), curr.String()))
			break
		}
	}
	return errors
}

type EqualTerm struct {
	Token         *token.Token
	Left          *RelationTerm
	EqualOperator []string // '=='|'!='
	Rights        []RelationTerm
}

func (et *EqualTerm) TokenLiteral() string {
	if et.Token != nil {
		return et.Token.Literal
	}
	panic("Could not determine token literals for equalTerm")
}
func (et *EqualTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(et.Left.String())
	for i, operator := range et.EqualOperator {
		relationTerm := et.Rights[i]
		out.WriteString(operator)
		out.WriteString(relationTerm.String())
	}
	return out.String()
}
func (et *EqualTerm) GetType(symTable *st.SymbolTable) types.Type {
	leftType := et.Left.GetType(symTable)

	for _, rTerm := range et.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (et *EqualTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = et.Left.TypeCheck(errors, symTable)
	for _, rTerm := range et.Rights {
		errors = rTerm.TypeCheck(errors, symTable)
	}
	return errors
}

type RelationTerm struct {
	Token             *token.Token
	Left              *SimpleTerm
	RelationOperators []string // '>'| '<' | '<=' | '>='
	Rights            []SimpleTerm
}

func (rt *RelationTerm) TokenLiteral() string {
	if rt.Token != nil {
		return rt.Token.Literal
	}
	panic("Could not determine token literals for relationTerm")
}
func (rt *RelationTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(rt.Left.String())
	for i, operator := range rt.RelationOperators {
		simpleTerm := rt.Rights[i]
		out.WriteString(operator)
		out.WriteString(simpleTerm.String())
	}
	return out.String()
}
func (rt *RelationTerm) GetType(symTable *st.SymbolTable) types.Type {
	leftType := rt.Left.GetType(symTable)

	for _, rTerm := range rt.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (rt *RelationTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = rt.Left.TypeCheck(errors, symTable)
	if len(rt.Rights) == 0 {
		return errors
	}
	// with + or - operations, every Term should be int type
	leftMostTy := rt.Left.GetType(symTable)
	if leftMostTy != types.IntTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (RelationTerm) expected int, found %v (%v)",
			rt.Token.LineNum, leftMostTy.GetName(), rt.Left.String()))
		return errors
	}
	for _, rTerm := range rt.Rights {
		errors = rTerm.TypeCheck(errors, symTable)
		currTy := rTerm.GetType(symTable)
		if currTy != types.IntTySig {
			errors = append(errors, fmt.Sprintf("[%v]: (RelationTerm) expected int, found %v (%v)",
				rTerm.Token.LineNum, currTy.GetName(), rTerm.String()))
			return errors
		}
	}
	return errors
}

type SimpleTerm struct {
	Token               *token.Token
	Left                *Term
	SimpleTermOperators []string // '+' | '-'
	Rights              []Term
}

func (st *SimpleTerm) TokenLiteral() string {
	if st.Token != nil {
		return st.Token.Literal
	}
	panic("Could not determine token literals for simpleTerm")
}
func (st *SimpleTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(st.Left.String())
	for i, operator := range st.SimpleTermOperators {
		term := st.Rights[i]
		out.WriteString(operator)
		out.WriteString(term.String())
	}
	return out.String()
}
func (st *SimpleTerm) GetType(symTable *st.SymbolTable) types.Type {
	leftType := st.Left.GetType(symTable)

	for _, rTerm := range st.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (st *SimpleTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = st.Left.TypeCheck(errors, symTable)
	if len(st.Rights) == 0 {
		return errors
	}
	// with + or - operations, every Term should be int type
	leftMostTy := st.Left.GetType(symTable)
	if leftMostTy != types.IntTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (SimpleTerm) expected int, found %v (%v)",
			st.Token.LineNum, leftMostTy.GetName(), st.Left.String()))
		return errors
	}
	for _, rTerm := range st.Rights {
		errors = rTerm.TypeCheck(errors, symTable)
		currTy := rTerm.GetType(symTable)
		if currTy != types.IntTySig {
			errors = append(errors, fmt.Sprintf("[%v]: (SimpleTerm) expected int, found %v (%v)",
				rTerm.Token.LineNum, currTy.GetName(), rTerm.String()))
			return errors
		}
	}
	return errors
}

type Term struct {
	Token         *token.Token
	Left          *UnaryTerm
	TermOperators []string // '*' | '/'
	Rights        []UnaryTerm
}

func (t *Term) TokenLiteral() string {
	if t.Token != nil {
		return t.Token.Literal
	}
	panic("Could not determine token literals for term")
}
func (t *Term) String() string {
	out := bytes.Buffer{}
	out.WriteString(t.Left.String())
	for i, operator := range t.TermOperators {
		unaryTerm := t.Rights[i]
		out.WriteString(operator)
		out.WriteString(unaryTerm.String())
	}
	return out.String()
}
func (t *Term) GetType(symTable *st.SymbolTable) types.Type {
	leftType := t.Left.GetType(symTable)

	for _, rTerm := range t.Rights {
		rightType := rTerm.GetType(symTable)
		if leftType != rightType {
			return types.UnknownTySig
		}
	}
	return leftType
}
func (t *Term) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = t.Left.TypeCheck(errors, symTable)
	if len(t.Rights) == 0 {
		return errors
	}
	// with * or / operations, every UnaryTerm should be int type
	leftMostTy := t.Left.GetType(symTable)
	if leftMostTy != types.IntTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (Term) expected int, found %v (%v)",
			t.Token.LineNum, leftMostTy.GetName(), t.Left.String()))
		return errors
	}
	for _, rTerm := range t.Rights {
		errors = rTerm.TypeCheck(errors, symTable)
		currTy := rTerm.GetType(symTable)
		if currTy != types.IntTySig {
			errors = append(errors, fmt.Sprintf("[%v]: (Term) expected int, found %v (%v)",
				rTerm.Token.LineNum, currTy.GetName(), rTerm.String()))
			return errors
		}
	}
	return errors
}

type UnaryTerm struct {
	Token         *token.Token
	UnaryOperator string // '!' | '-' | '' <- default
	SelectorTerm  *SelectorTerm
}

func (ut *UnaryTerm) TokenLiteral() string {
	if ut.Token != nil {
		return ut.Token.Literal
	}
	panic("Could not determine token literals for unaryTerm")
}
func (ut *UnaryTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(ut.UnaryOperator)
	out.WriteString(ut.SelectorTerm.String())
	return out.String()
}
func (ut *UnaryTerm) GetType(symTable *st.SymbolTable) types.Type {
	if ut.UnaryOperator == "!" {
		if ut.SelectorTerm.GetType(symTable) == types.BoolTySig {
			return types.BoolTySig
		} else {
			return types.UnknownTySig
		}
	} else if ut.UnaryOperator == "-" {
		if ut.SelectorTerm.GetType(symTable) == types.IntTySig {
			return types.IntTySig
		} else {
			return types.UnknownTySig
		}
	} else {
		return ut.SelectorTerm.GetType(symTable)
	}
}
func (ut *UnaryTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = ut.SelectorTerm.TypeCheck(errors, symTable)
	if ut.GetType(symTable) == types.UnknownTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (UnaryTerm) Unkown type", ut.Token.LineNum))
	}
	return errors
}

type SelectorTerm struct {
	Token  *token.Token
	Fact   *Factor
	Idents []IdentLiteral
}

func (selt *SelectorTerm) TokenLiteral() string {
	if selt.Token != nil {
		return selt.Token.Literal
	}
	panic("Could not determine token literals for selectorTerm")
}
func (selt *SelectorTerm) String() string {
	out := bytes.Buffer{}
	out.WriteString(selt.Fact.String())
	for _, id := range selt.Idents {
		out.WriteString(".")
		out.WriteString(id.String())
	}
	return out.String()
}
func (selt *SelectorTerm) GetType(symTable *st.SymbolTable) types.Type {
	facType := selt.Fact.GetType(symTable)
	if len(selt.Idents) == 0 {
		return facType
	} else if facType == types.StructTySig {
		var entry st.Entry
		varName := selt.Fact.String()
		for {
			if entry = symTable.Contains(varName); entry == nil {
				if symTable.Parent == nil {
					return types.UnknownTySig
				} else {
					symTable = symTable.Parent
				}
			} else {
				break
			}
		}
		if selt.Idents == nil {
			return entry.GetEntryType()
		}
		// here entry is the entry of the first id in Idents
		symTable = entry.GetScopeST()
		remaining := selt.Idents[1:]
		for idx, id := range remaining {
			if entry = symTable.Contains(id.String()); entry == nil {
				return types.UnknownTySig
			} else {
				if idx == len(selt.Idents)-1 {
					break
				}
				symTable = entry.GetScopeST()
			}
		}
		return entry.GetEntryType()
	}
	return types.UnknownTySig
}
func (selt *SelectorTerm) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = selt.Fact.TypeCheck(errors, symTable)
	if selt.GetType(symTable) == types.UnknownTySig {
		errors = append(errors, fmt.Sprintf("[%v]: (SelectorTerm) Unknown type", selt.Fact.Token.LineNum))
		return errors
	}
	return errors
}

type Factor struct {
	Token *token.Token
	Expr  Expr
}

func (f *Factor) TokenLiteral() string {
	if f.Token != nil {
		return f.Token.Literal
	}
	panic("Could not determine token literal for factor")
}
func (f *Factor) String() string {
	out := bytes.Buffer{}
	out.WriteString(f.Expr.String())
	return out.String()
}
func (f *Factor) GetType(symTable *st.SymbolTable) types.Type {
	return f.Expr.GetType(symTable)
}
func (f *Factor) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = f.Expr.TypeCheck(errors, symTable)
	return errors
}

func NewType(typeLit string) *Type                                { return &Type{nil, typeLit} }
func NewArgs(exprs []Expression) *Arguments                       { return &Arguments{nil, exprs} }
func NewLvalue(ident IdentLiteral, idents []IdentLiteral) *LValue { return &LValue{nil, ident, idents} }
func NewExpression(l *BoolTerm, rs []BoolTerm) *Expression {
	return &Expression{nil, l, rs}
}
func NewBoolTerm(l *EqualTerm, rs []EqualTerm) *BoolTerm { return &BoolTerm{nil, l, rs} }
func NewEqualTerm(l *RelationTerm, operators []string, rs []RelationTerm) *EqualTerm {
	return &EqualTerm{nil, l, operators, rs}
}
func NewRelationTerm(l *SimpleTerm, operators []string, rs []SimpleTerm) *RelationTerm {
	return &RelationTerm{nil, l, operators, rs}
}
func NewSimpleTerm(l *Term, operators []string, rs []Term) *SimpleTerm {
	return &SimpleTerm{nil, l, operators, rs}
}
func NewTerm(l *UnaryTerm, operators []string, rs []UnaryTerm) *Term {
	return &Term{nil, l, operators, rs}
}
func NewUnaryTerm(operator string, selectorTerm *SelectorTerm) *UnaryTerm {
	return &UnaryTerm{nil, operator, selectorTerm}
}
func NewSelectorTerm(factor *Factor, idents []IdentLiteral) *SelectorTerm {
	return &SelectorTerm{nil, factor, idents}
}
func NewFactor(expr *Expr) *Factor { return &Factor{nil, *expr} }

// InvocExpr : invocation in Factor ('id' [Arguments])
type InvocExpr struct {
	Token     *token.Token
	Ident     IdentLiteral
	InnerArgs *Arguments
}

func (ie *InvocExpr) TokenLiteral() string {
	if ie.Token != nil {
		return ie.Token.Literal
	}
	panic("Could not determine token literal for invocation expression inside Factor")
}
func (ie *InvocExpr) String() string {
	out := bytes.Buffer{}
	out.WriteString(ie.Ident.String())
	out.WriteString(ie.InnerArgs.String())
	return out.String()
}
func (ie *InvocExpr) GetType(symTable *st.SymbolTable) types.Type {
	if funcEntry := symTable.Contains(ie.Ident.Id); funcEntry != nil {
		return funcEntry.GetReturnTy()
	}
	return types.UnknownTySig
}
func (ie *InvocExpr) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	// refer from Invocation.TypeCheck
	funcName := ie.Ident.TokenLiteral()
	entry := symTable.Contains(funcName)
	if entry == nil {
		errors = append(errors, fmt.Sprintf("[%v]: function %v has not been defined", ie.Token.LineNum, funcName))
	} else {
		symTable = entry.GetScopeST()
		errors = ie.InnerArgs.TypeCheck(errors, symTable)
	}
	return errors
}

// PriorityExpression : '(' Expression ')' (inside Factor)
type PriorityExpression struct {
	Token           *token.Token
	InnerExpression *Expression
}

func (pe *PriorityExpression) TokenLiteral() string {
	if pe.Token != nil {
		return pe.Token.Literal
	}
	panic("Could not determine token literal for expression inside Factor")
}
func (pe *PriorityExpression) String() string {
	out := bytes.Buffer{}
	out.WriteString("(")
	out.WriteString(pe.InnerExpression.String())
	out.WriteString(")")
	return out.String()
}
func (pe *PriorityExpression) GetType(symTable *st.SymbolTable) types.Type {
	return pe.InnerExpression.GetType(symTable)
}
func (pe *PriorityExpression) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	errors = pe.InnerExpression.TypeCheck(errors, symTable)
	return errors
}

// NilNode : nil (keyword "nil")
type NilNode struct {
	Token *token.Token
}

func (n *NilNode) TokenLiteral() string                        { return n.Token.Literal }
func (n *NilNode) String() string                              { return n.Token.Literal }
func (n *NilNode) GetType(symTable *st.SymbolTable) types.Type { return types.VoidTySig }
func (n *NilNode) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

// BoolLiteral : True/False
type BoolLiteral struct {
	Token *token.Token
	Value bool
}

func (bl *BoolLiteral) TokenLiteral() string                        { return bl.Token.Literal }
func (bl *BoolLiteral) String() string                              { return bl.Token.Literal }
func (bl *BoolLiteral) GetType(symTable *st.SymbolTable) types.Type { return types.BoolTySig }
func (bl *BoolLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

// IntLiteral : number (integer)
type IntLiteral struct {
	Token *token.Token
	Value int64
}

func (il *IntLiteral) TokenLiteral() string                        { return il.Token.Literal }
func (il *IntLiteral) String() string                              { return il.Token.Literal }
func (il *IntLiteral) GetType(symTable *st.SymbolTable) types.Type { return types.IntTySig }
func (il *IntLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	return errors
}

// IdentLiteral : identifier
type IdentLiteral struct {
	Token *token.Token
	Id    string
}

func (idl *IdentLiteral) TokenLiteral() string { return idl.Token.Literal }
func (idl *IdentLiteral) String() string       { return idl.Token.Literal }
func (idl *IdentLiteral) GetType(symTable *st.SymbolTable) types.Type {
	var entry st.Entry
	for {
		if entry = symTable.Contains(idl.TokenLiteral()); entry == nil {
			if symTable.Parent == nil {
				return types.UnknownTySig
			} else {
				symTable = symTable.Parent
			}
		} else {
			break
		}
	}
	return entry.GetEntryType()
}
func (idl *IdentLiteral) TypeCheck(errors []string, symTable *st.SymbolTable) []string {
	if idl.GetType(symTable) == types.UnknownTySig {
		errors = append(errors, fmt.Sprintf("[%v]: %v has not been defined.", idl.Token.LineNum, idl.Id))
	}
	return errors
}
