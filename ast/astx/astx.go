// Package astx implements utility functions for generating abstract syntax
// trees.
package astx

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/errors"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// NewFile returns a new µC source file, based on the following production rule.
//
//    File
//       : DeclStmts
//    ;
func NewFile(decls interface{}) (*ast.File, error) {
	if decls, ok := decls.([]ast.Decl); ok {
		return &ast.File{Decls: decls}, nil
	}
	return nil, errutil.Newf("invalid file declarations type; expected []ast.Decl, got %T", decls)
}

// NewDeclList returns a new declaration list, based on the following production
// rule.
//
//    DeclList
//       : Decl
//    ;
func NewDeclList(decl interface{}) ([]ast.Decl, error) {
	if decl, ok := decl.(ast.Decl); ok {
		return []ast.Decl{decl}, nil
	}
	return nil, errutil.Newf("invalid declaration list declaration type; expected ast.Decl, got %T", decl)
}

// AppendDecl appends decl to the declaration list, based on the following
// production rule.
//
//    DeclList
//       : DeclList Decl
//    ;
func AppendDecl(list, decl interface{}) ([]ast.Decl, error) {
	lst, ok := list.([]ast.Decl)
	if !ok {
		return nil, errutil.Newf("invalid declaration list type; expected []ast.Decl, got %T", list)
	}
	if decl, ok := decl.(ast.Decl); ok {
		return append(lst, decl), nil
	}
	return nil, errutil.Newf("invalid declaration list declaration type; expected ast.Decl, got %T", decl)
}

// NewFuncDecl returns a new function declaration node, based on the following
// production rule.
//
//    FuncDecl
//       : FuncHeader
//    ;
//
//    FuncHeader
//       : BasicType ident "(" Params ")"
//    ;
func NewFuncDecl(resultType, name, params interface{}) (*ast.FuncDecl, error) {
	resType, err := NewType(resultType)
	if err != nil {
		return nil, errutil.Newf("invalid function result type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid function name identifier; %v", err)
	}
	fields, ok := params.([]*types.Field)
	if !ok {
		return nil, errutil.Newf("invalid function parameters type; expected []*types.Field, got %T", params)
	}
	typ := &types.Func{Params: fields, Result: resType}
	return &ast.FuncDecl{Type: typ, Name: ident}, nil
}

// SetFuncBody sets the function body of the given function declaration, based
// on the following production rule.
//
//    FuncDef
//       : FuncHeader BlockStmt
//    ;
func SetFuncBody(f, body interface{}) (*ast.FuncDecl, error) {
	fn, ok := f.(*ast.FuncDecl)
	if !ok {
		return nil, errutil.Newf("invalid function declaration type; expected *ast.FuncDecl, got %T", f)
	}
	if body, ok := body.(*ast.BlockStmt); ok {
		fn.Body = body
		return fn, nil
	}
	return nil, errutil.Newf("invalid function body type; expected *ast.BlockStmt, got %T", body)

}

// NewScalarDecl returns a new scalar declaration node, based on the following
// production rule.
//
//    ScalarDecl
//       : TypeName ident
//    ;
func NewScalarDecl(typ, name interface{}) (*ast.VarDecl, error) {
	scalarType, err := NewType(typ)
	if err != nil {
		return nil, errutil.Newf("invalid scalar type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid scalar declaration identifier; %v", err)
	}
	return &ast.VarDecl{Type: scalarType, Name: ident}, nil
}

// NewArrayDecl returns a new array declaration node, based on the following
// production rule.
//
//    ArrayDecl
//       : TypeName ident "[" int_lit "]"
//    ;
func NewArrayDecl(elem, name, length interface{}) (*ast.VarDecl, error) {
	typ, err := NewArrayType(elem, length)
	if err != nil {
		return nil, errutil.Newf("invalid array type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid array declaration identifier; %v", err)
	}
	return &ast.VarDecl{Type: typ, Name: ident}, nil
}

// NewVoidParam returns a new void parameter, based on the following production
// rule.
//
//    Params
//    	: TypeName // "void"
//    ;
func NewVoidParam(typ interface{}) ([]*types.Field, error) {
	if typ, ok := typ.(*types.Basic); ok {
		switch typ.Kind {
		case types.Void:
			// Valid void parameter.
			return []*types.Field{{Type: typ}}, nil
		default:
			return nil, errutil.Newf("invalid void parameter kind; expected Void, got %v", typ.Kind)
		}
	}
	return nil, errutil.Newf("invalid void parameter type; expected *types.Basic, got %T", typ)
}

// NewFieldList returns a new field list, based on the following production
// rule.
//
//    FieldList
//       : Field
//    ;
func NewFieldList(field interface{}) ([]*types.Field, error) {
	if field, ok := field.(*types.Field); ok {
		return []*types.Field{field}, nil
	}
	return nil, errutil.Newf("invalid field list field type; expected *types.Field, got %T", field)
}

// AppendField appends field to the field list, based on the following
// production rule.
//
//    FieldList
//       : FieldList "," Field
//    ;
func AppendField(list, field interface{}) ([]*types.Field, error) {
	lst, ok := list.([]*types.Field)
	if !ok {
		return nil, errutil.Newf("invalid field list type; expected []*types.Field, got %T", list)
	}
	if field, ok := field.(*types.Field); ok {
		return append(lst, field), nil
	}
	return nil, errutil.Newf("invalid field list field type; expected *types.Field, got %T", field)
}

// NewField returns a new field, based on the following production rules.
//
//    ParamDecl
//       : ScalarDecl
//       | TypeName ident "[" "]"
//    ;
func NewField(decl interface{}) (*types.Field, error) {
	if decl, ok := decl.(*ast.VarDecl); ok {
		return &types.Field{Type: decl.Type, Name: decl.Name.Name}, nil
	}
	return nil, errutil.Newf("invalid field type; expected *ast.VarDecl, got %T", decl)
}

// NewExprStmt returns a new expression statement, based on the following
// production rule.
//
//    Stmt
//       : Expr ";"
//    ;
func NewExprStmt(x interface{}) (*ast.ExprStmt, error) {
	if x, ok := x.(ast.Expr); ok {
		return &ast.ExprStmt{X: x}, nil
	}
	return nil, errutil.Newf("invalid expression statement expression type; expected ast.Expr, got %T", x)
}

// NewReturnStmt returns a new return statement, based on the following
// production rule.
//
//    Stmt
//       : "return" Expr ";"
//    ;
func NewReturnStmt(result interface{}) (*ast.ReturnStmt, error) {
	if result, ok := result.(ast.Expr); ok {
		return &ast.ReturnStmt{Result: result}, nil
	}
	return nil, errutil.Newf("invalid return statement result type; expected ast.Expr, got %T", result)
}

// NewWhileStmt returns a new while statement, based on the following production
// rule.
//
//    Stmt
//       : "while" Condition Stmt
//    ;
func NewWhileStmt(cond, body interface{}) (*ast.WhileStmt, error) {
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid while statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := body.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid while statement body type; expected ast.Stmt, got %T", body)
	}
	return &ast.WhileStmt{Cond: condExpr, Body: bodyStmt}, nil
}

// NewIfStmt returns a new if statement, based on the following production
// rules.
//
//    Stmt
//       : "if" Condition Stmt ElsePart
//    ;
//
//    ElsePart
//       : empty
//       | "else" Stmt
//    ;
func NewIfStmt(cond, trueBranch, falseBranch interface{}) (*ast.IfStmt, error) {
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid if statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := trueBranch.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid if statement body type; expected ast.Stmt, got %T", trueBranch)
	}
	if falseBranch == nil {
		return &ast.IfStmt{Cond: condExpr, Body: bodyStmt}, nil
	}
	if elseStmt, ok := falseBranch.(ast.Stmt); ok {
		return &ast.IfStmt{Cond: condExpr, Body: bodyStmt, Else: elseStmt}, nil
	}
	return nil, errutil.Newf("invalid if statement else-body type; expected ast.Stmt, got %T", falseBranch)
}

// NewBlockStmt returns a new block statement, based on the following production
// rule.
//
//    BlockStmt
//       : "{" BlockItems "}"
//    ;
func NewBlockStmt(items interface{}) (*ast.BlockStmt, error) {
	if items == nil {
		return &ast.BlockStmt{}, nil
	}
	if items, ok := items.([]ast.BlockItem); ok {
		return &ast.BlockStmt{Items: items}, nil
	}
	return nil, errutil.Newf("invalid block statements type; expected []ast.BlockItem, got %T", items)
}

// NewBlockItemList returns a new block item list, based on the following
// production rule.
//
//    BlockItemList
//       : BlockItem
//    ;
func NewBlockItemList(item interface{}) ([]ast.BlockItem, error) {
	if item, ok := item.(ast.BlockItem); ok {
		return []ast.BlockItem{item}, nil
	}
	return nil, errutil.Newf("invalid block item list block item type; expected ast.BlockItem, got %T", item)
}

// AppendBlockItem appends item to the block item list, based on the following
// production rule.
//
//    BlockItemList
//       : BlockItemList BlockItem
//    ;
func AppendBlockItem(list, item interface{}) ([]ast.BlockItem, error) {
	lst, ok := list.([]ast.BlockItem)
	if !ok {
		return nil, errutil.Newf("invalid block item list type; expected []ast.BlockItem, got %T", list)
	}
	if item, ok := item.(ast.BlockItem); ok {
		return append(lst, item), nil
	}
	return nil, errutil.Newf("invalid block item list block item type; expected ast.BlockItem, got %T", item)
}

// NewBinaryExpr returns a new binary experssion node, based on the following
// production rules.
//
//    Expr2R
//       : Expr2R "=" Expr5L
//    ;
//
//    Expr5L
//       : Expr5L "&&" Expr9L
//    ;
//
//    Expr9L
//       : Expr9L "==" Expr10L
//       | Expr9L "!=" Expr10L
//    ;
//
//    Expr10L
//       : Expr10L "<" Expr12L
//       | Expr10L ">" Expr12L
//       | Expr10L "<=" Expr12L
//       | Expr10L ">=" Expr12L
//    ;
//
//    Expr12L
//       : Expr12L "+" Expr13L
//       | Expr12L "-" Expr13L
//    ;
//
//    Expr13L
//       : Expr13L "*" Expr14
//       | Expr13L "/" Expr14
//    ;
func NewBinaryExpr(x interface{}, op token.Kind, y interface{}) (*ast.BinaryExpr, error) {
	switch op {
	case token.Assign,
		token.Land,
		token.Eq, token.Ne,
		token.Lt, token.Gt, token.Le, token.Ge,
		token.Add, token.Sub,
		token.Mul, token.Div:
		// Valid binary operator.
	default:
		return nil, errutil.Newf("invalid binary operator; expected Assign, Land, Eq, Ne, Lt, Gt, Le, Ge, Add, Sub, Mul or Div, got %v", op)
	}
	arg0, ok := x.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid first binary operand type; expected ast.Expr, got %T", x)
	}
	arg1, ok := y.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid second binary operand type; expected ast.Expr, got %T", y)
	}
	return &ast.BinaryExpr{X: arg0, Op: op, Y: arg1}, nil
}

// NewUnaryExpr returns a new unary experssion node, based on the following
// production rules.
//
//    Expr14
//       : "-" Expr15
//       | "!" Expr15
//    ;
func NewUnaryExpr(op token.Kind, x interface{}) (*ast.UnaryExpr, error) {
	switch op {
	case token.Sub, token.Not:
		// Valid unary operator.
	default:
		return nil, errutil.Newf("invalid unary operator; expected Sub or Not, got %v", op)
	}
	if x, ok := x.(ast.Expr); ok {
		return &ast.UnaryExpr{Op: op, X: x}, nil
	}
	return nil, errutil.Newf("invalid unary operand type; expected ast.Expr, got %T", x)
}

// TODO: Add char_lit production rule to NewBasicLit doc comment once handled
// explicitly in uc.bnf.

// NewBasicLit returns a new basic literal experssion node of the given kind,
// based on the following production rule.
//
//    Expr15
//       : int_lit
//    ;
func NewBasicLit(val interface{}, kind token.Kind) (*ast.BasicLit, error) {
	s, err := tokenString(val)
	if err != nil {
		return nil, errutil.Newf("invalid basic literal value; %v", err)
	}
	switch kind {
	case token.CharLit, token.IntLit:
		// Valid kind.
	default:
		return nil, errutil.Newf("invalid basic literal kind; expected CharLit or IntLit, got %v", kind)
	}
	return &ast.BasicLit{Kind: kind, Val: s}, nil
}

// NewIdent returns a new identifier experssion node, based on the following
// production rule.
//
//    Expr15
//       : ident
//    ;
func NewIdent(name interface{}) (*ast.Ident, error) {
	s, err := tokenString(name)
	if err != nil {
		return nil, errutil.Newf("invalid identifier; %v", err)
	}
	return &ast.Ident{Name: s}, nil
}

// NewIndexExpr returns a new index expression, based on the following
// production rule.
//
//    Expr15
//       : ident "[" Expr "]"
//    ;
func NewIndexExpr(name, index interface{}) (*ast.IndexExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid array name; %v", err)
	}
	if index, ok := index.(ast.Expr); ok {
		return &ast.IndexExpr{Name: ident, Index: index}, nil
	}
	return nil, errutil.Newf("invalid index expression type; expected ast.Expr, got %T", index)
}

// NewCallExpr returns a new call expression, based on the following production
// rule.
//
//    Expr15
//       : ident "(" Args ")"
//    ;
func NewCallExpr(name, args interface{}) (*ast.CallExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid function name; %v", err)
	}
	if args == nil {
		return &ast.CallExpr{Name: ident}, nil
	}
	if args, ok := args.([]ast.Expr); ok {
		return &ast.CallExpr{Name: ident, Args: args}, nil
	}
	return nil, errutil.Newf("invalid function arguments type; expected []ast.Expr, got %T", args)
}

// NewParenExpr returns a new parenthesized expression, based on the following
// production rule.
//
//    ParenExpr
//       : "(" Expr ")"
//    ;
func NewParenExpr(x interface{}) (*ast.ParenExpr, error) {
	if x, ok := x.(ast.Expr); ok {
		return &ast.ParenExpr{X: x}, nil
	}
	return nil, errutil.Newf("invalid parenthesized expression type; expected ast.Expr, got %T", x)
}

// NewExprList returns a new expression list, based on the following production
// rule.
//
//    ExprList
//       : Expr
//    ;
func NewExprList(x interface{}) ([]ast.Expr, error) {
	if x, ok := x.(ast.Expr); ok {
		return []ast.Expr{x}, nil
	}
	return nil, errutil.Newf("invalid expression list expression type; expected ast.Expr, got %T", x)
}

// AppendExpr appends x to the expression list, based on the following
// production rule.
//
//    ExprList
//       : ExprList "," Expr
//    ;
func AppendExpr(list, x interface{}) ([]ast.Expr, error) {
	lst, ok := list.([]ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid expression list type; expected []ast.Expr, got %T", list)
	}
	if x, ok := x.(ast.Expr); ok {
		return append(lst, x), nil
	}
	return nil, errutil.Newf("invalid expression list expression type; expected ast.Expr, got %T", x)
}

// tokenString returns the lexeme of the given token.
func tokenString(tok interface{}) (string, error) {
	if tok, ok := tok.(*gocctoken.Token); ok {
		return string(tok.Lit), nil
	}
	return "", errutil.Newf("invalid tok type; expected *gocctoken.Token, got %T", tok)
}

// TODO: Remove NewExpectedError if it is not in use after the parser error
// handling has matured.

// NewExpectedError returns a user-friendly parse error of the following form.
//
//    pos: unexpected "foo", expected bar
//
// where pos is the position (in bytes) within the input stream, foo the next
// token in the input stream and bar a description of the expected token or
// production.
func NewExpectedError(err interface{}, expected string) (ast.Node, error) {
	if err, ok := err.(*errors.Error); ok {
		pos := err.ErrorToken.Pos.Offset
		return nil, fmt.Errorf("%d: unexpected %q, expected %v", pos, string(err.ErrorToken.Lit), expected)
	}
	return nil, errutil.Newf("invalid error type; expected *errors.Error, got %T", err)
}