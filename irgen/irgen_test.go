package irgen_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/irgen"
	"github.com/mewmew/uc/sem"
)

func TestGen(t *testing.T) {
	golden := []struct {
		path string
		want string
	}{
		// Identifier usage.
		{
			path: "../testdata/extra/irgen/int_ident_use.c",
			want: "../testdata/extra/irgen/int_ident_use.ll",
		},
		{
			path: "../testdata/extra/irgen/int_ident_def.c",
			want: "../testdata/extra/irgen/int_ident_def.ll",
		},
		{
			path: "../testdata/extra/irgen/array_ident_use.c",
			want: "../testdata/extra/irgen/array_ident_use.ll",
		},
		// NOTE: Not part of grammar for uC.
		//{
		//	path: "../testdata/extra/irgen/array_ident_def.c",
		//	want: "../testdata/extra/irgen/array_ident_def.ll",
		//},
		{
			path: "../testdata/extra/irgen/index_expr_use.c",
			want: "../testdata/extra/irgen/index_expr_use.ll",
		},
		{
			path: "../testdata/extra/irgen/index_expr_def.c",
			want: "../testdata/extra/irgen/index_expr_def.ll",
		},
		{
			path: "../testdata/extra/irgen/global_index_expr_use.c",
			want: "../testdata/extra/irgen/global_index_expr_use.ll",
		},
		{
			path: "../testdata/extra/irgen/global_index_expr_def.c",
			want: "../testdata/extra/irgen/global_index_expr_def.ll",
		},

		// Global variable declarations.
		{
			path: "../testdata/extra/irgen/global_def.c",
			want: "../testdata/extra/irgen/global_def.ll",
		},
		{
			path: "../testdata/extra/irgen/tentative_def.c",
			want: "../testdata/extra/irgen/tentative_def.ll",
		},
		// Return statements.
		{
			path: "../testdata/extra/irgen/void_ret.c",
			want: "../testdata/extra/irgen/void_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/implicit_void_ret.c",
			want: "../testdata/extra/irgen/implicit_void_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/int_ret.c",
			want: "../testdata/extra/irgen/int_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/expr_ret.c",
			want: "../testdata/extra/irgen/expr_ret.ll",
		},
		{
			path: "../testdata/extra/irgen/global_ret.c",
			want: "../testdata/extra/irgen/global_ret.ll",
		},
		// Local variable declarations.
		{
			path: "../testdata/extra/irgen/local_def.c",
			want: "../testdata/extra/irgen/local_def.ll",
		},
		// If and if-else statements.
		{
			path: "../testdata/extra/irgen/if_stmt.c",
			want: "../testdata/extra/irgen/if_stmt.ll",
		},
		{
			path: "../testdata/extra/irgen/if_else_stmt.c",
			want: "../testdata/extra/irgen/if_else_stmt.ll",
		},
		// Parenthesized expressions.
		{
			path: "../testdata/extra/irgen/paren_expr.c",
			want: "../testdata/extra/irgen/paren_expr.ll",
		},
		// Unary expressions.
		{
			path: "../testdata/extra/irgen/unary_expr_sub.c",
			want: "../testdata/extra/irgen/unary_expr_sub.ll",
		},
		{
			path: "../testdata/extra/irgen/unary_expr_not.c",
			want: "../testdata/extra/irgen/unary_expr_not.ll",
		},
		// Binary expressions.
		{
			path: "../testdata/extra/irgen/binary_expr_add.c",
			want: "../testdata/extra/irgen/binary_expr_add.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_sub.c",
			want: "../testdata/extra/irgen/binary_expr_sub.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_mul.c",
			want: "../testdata/extra/irgen/binary_expr_mul.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_div.c",
			want: "../testdata/extra/irgen/binary_expr_div.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_lt.c",
			want: "../testdata/extra/irgen/binary_expr_lt.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_gt.c",
			want: "../testdata/extra/irgen/binary_expr_gt.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_le.c",
			want: "../testdata/extra/irgen/binary_expr_le.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_ge.c",
			want: "../testdata/extra/irgen/binary_expr_ge.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_ne.c",
			want: "../testdata/extra/irgen/binary_expr_ne.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_eq.c",
			want: "../testdata/extra/irgen/binary_expr_eq.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_land.c",
			want: "../testdata/extra/irgen/binary_expr_land.ll",
		},
		{
			path: "../testdata/extra/irgen/binary_expr_assign.c",
			want: "../testdata/extra/irgen/binary_expr_assign.ll",
		},
		{
			path: "../testdata/extra/irgen/while_stmt.c",
			want: "../testdata/extra/irgen/while_stmt.ll",
		},
		// Function parameters.
		{
			path: "../testdata/extra/irgen/func_param.c",
			want: "../testdata/extra/irgen/func_param.ll",
		},
		// Nested variable declarations.
		{
			path: "../testdata/extra/irgen/nested_var_decl.c",
			want: "../testdata/extra/irgen/nested_var_decl.ll",
		},
		// Call expressions.
		{
			path: "../testdata/extra/irgen/call_expr.c",
			want: "../testdata/extra/irgen/call_expr.ll",
		},
		{
			path: "../testdata/extra/irgen/call_expr_cast.c",
			want: "../testdata/extra/irgen/call_expr_cast.ll",
		},
		// NOTE: Correct output. The only difference is that Clang emits all
		// alloca instructions at the beginning of the entry block. Thus, disabled
		// for now.
		//{
		//	path: "../testdata/extra/irgen/call_expr_multi_args.c",
		//	want: "../testdata/extra/irgen/call_expr_multi_args.ll",
		//},
		// NOTE: Correct output. The only difference is that Clang emits all
		// alloca instructions at the beginning of the entry block. Thus, disabled
		// for now.
		//{
		//	path: "../testdata/extra/irgen/call_expr_multi_args_cast.c",
		//	want: "../testdata/extra/irgen/call_expr_multi_args_cast.ll",
		//},
		// Expression statements.
		{
			path: "../testdata/extra/irgen/expr_stmt.c",
			want: "../testdata/extra/irgen/expr_stmt.ll",
		},
		// Index expressions.
		{
			path: "../testdata/extra/irgen/index_expr.c",
			want: "../testdata/extra/irgen/index_expr.ll",
		},
		// Array arguments.
		{
			path: "../testdata/extra/irgen/array_arg.c",
			want: "../testdata/extra/irgen/array_arg.ll",
		},
		// Array parameters.
		// TODO: Re-enable; https://github.com/mewmew/uc/issues/73
		{
			path: "../testdata/extra/irgen/array_param.c",
			want: "../testdata/extra/irgen/array_param.ll",
		},
		// Bug fixes.
		{
			path: "../testdata/extra/irgen/issue_68_nested_if.c",
			want: "../testdata/extra/irgen/issue_68_nested_if.ll",
		},
		{
			path: "../testdata/extra/irgen/issue_68_nested_if_while.c",
			want: "../testdata/extra/irgen/issue_68_nested_if_while.ll",
		},
		// TODO: Re-enable when better support for logical AND expressions have
		// been implemented.
		//{
		//	path: "../testdata/extra/irgen/issue_68_if_land.c",
		//	want: "../testdata/extra/irgen/issue_68_if_land.ll",
		//},
		{
			path: "../testdata/extra/irgen/issue_68_while_land.c",
			want: "../testdata/extra/irgen/issue_68_while_land.ll",
		},
		// NOTE: The test case "issue_70_if_ret.c" no longer panics on SetTerm
		// with a nil basic block. The test case still doesn't pass, but this is
		// simply because Clang allocates a dedicated variable for return values,
		// which differs from how uclang does. Thus, the test case remains
		// disabled for now.
		//{
		//	path: "../testdata/extra/irgen/issue_70_if_ret.c",
		//	want: "../testdata/extra/irgen/issue_70_if_ret.ll",
		//},
		// NOTE: The test case "issue_70_while_ret.c" no longer panics on SetTerm
		// with a nil basic block. The test case still doesn't pass, but this is
		// simply because Clang allocates a dedicated variable for return values,
		// which differs from how uclang does. Thus, the test case remains
		// disabled for now.
		//{
		//	path: "../testdata/extra/irgen/issue_70_while_ret.c",
		//	want: "../testdata/extra/irgen/issue_70_while_ret.ll",
		//},
		{
			path: "../testdata/extra/irgen/issue_69_trunc_arg.c",
			want: "../testdata/extra/irgen/issue_69_trunc_arg.ll",
		},
		// TODO: Re-enable; https://github.com/mewmew/uc/issues/73
		{
			path: "../testdata/extra/irgen/issue_73_pointer_pointer_use.c",
			want: "../testdata/extra/irgen/issue_73_pointer_pointer_use.ll",
		},
		// TODO: Re-enable; https://github.com/mewmew/uc/issues/73
		{
			path: "../testdata/extra/irgen/issue_73_pointer_pointer_ref.c",
			want: "../testdata/extra/irgen/issue_73_pointer_pointer_ref.ll",
		},
		{
			path: "../testdata/extra/irgen/global_array_arg.c",
			want: "../testdata/extra/irgen/global_array_arg.ll",
		},
		{
			path: "../testdata/extra/irgen/global_array_ident_use.c",
			want: "../testdata/extra/irgen/global_array_ident_use.ll",
		},
		// NOTE: Not part of grammar for uC.
		// {
		// 	path: "../testdata/extra/irgen/global_array_ident_def.c",
		// 	want: "../testdata/extra/irgen/global_array_ident_def.ll",
		// },
		{
			path: "../testdata/extra/irgen/global_array_param.c",
			want: "../testdata/extra/irgen/global_array_param.ll",
		},
	}

	for _, g := range golden {
		// Lex input.
		buf, err := ioutil.ReadFile(g.path)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		input := string(buf)
		s := scanner.NewFromString(input)

		// Parse input.
		p := parser.NewParser()
		f, err := p.Parse(s)
		if err != nil {
			t.Errorf("%q: parse error: %v", g.path, err)
			continue
		}
		file := f.(*ast.File)

		// Verify input.
		info, err := sem.Check(file)
		if err != nil {
			t.Errorf("%q: semantic analysis error: %v", g.path, err)
			continue
		}

		// Generate IR.
		// TODO: Remove debug output.
		fmt.Printf("\n=== [ %s ] ====================================\n\n", g.want)
		module := irgen.Gen(file, info)

		// Compare generated module against gold standard.
		buf, err = ioutil.ReadFile(g.want)
		if err != nil {
			t.Errorf("%q: %v", g.path, err)
			continue
		}
		got, want := module.String(), string(buf)
		if got != want {
			t.Errorf("%q: module mismatch; expected `%v`, got `%v`", g.path, want, got)
			// TODO: Remove debug output.
			fmt.Println("### FAIL:", g.path)
			continue
		}
		// TODO: Remove debug output.
		fmt.Println("PASS:", g.path)
	}
}
