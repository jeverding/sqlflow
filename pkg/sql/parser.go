// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goyacc -p sql -o parser.go sql.y. DO NOT EDIT.

//line sql.y:2
package sql

import __yyfmt__ "fmt"

//line sql.y:2

import (
	"fmt"
	"strings"
	"sync"
)

/* expr defines an expression as a Lisp list.  If len(val)>0,
   it is an atomic expression, in particular, NUMBER, IDENT,
   or STRING, defined by typ and val; otherwise, it is a
   Lisp S-expression. */
type expr struct {
	typ  int
	val  string
	sexp exprlist
}

type exprlist []*expr

/* construct an atomic expr */
func atomic(typ int, val string) *expr {
	return &expr{
		typ: typ,
		val: val,
	}
}

/* construct a funcall expr */
func funcall(name string, oprd exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(IDENT, name)}, oprd...),
	}
}

/* construct a unary expr */
func unary(typ int, op string, od1 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1),
	}
}

/* construct a binary expr */
func binary(typ int, od1 *expr, op string, od2 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1, od2),
	}
}

/* construct a variadic expr */
func variadic(typ int, op string, ods exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, ods...),
	}
}

type extendedSelect struct {
	extended bool
	train    bool
	analyze  bool
	standardSelect
	trainClause
	predictClause
	analyzeClause
}

type standardSelect struct {
	fields exprlist
	tables []string
	where  *expr
	limit  string
}

type trainClause struct {
	estimator  string
	trainAttrs attrs
	columns    columnClause
	label      string
	save       string
}

/* If no FOR in the COLUMN, the key is "" */
type columnClause map[string]exprlist
type filedClause exprlist

type attrs map[string]*expr

type predictClause struct {
	predAttrs attrs
	model     string
	into      string
}

type analyzeClause struct {
	analyzeAttrs attrs
	trainedModel string
	explainer    string
}

var parseResult *extendedSelect

func attrsUnion(as1, as2 attrs) attrs {
	for k, v := range as2 {
		as1[k] = v
	}
	return as1
}

//line sql.y:111
type sqlSymType struct {
	yys  int
	val  string /* NUMBER, IDENT, STRING, and keywords */
	flds exprlist
	tbls []string
	expr *expr
	expl exprlist
	atrs attrs
	eslt extendedSelect
	slct standardSelect
	tran trainClause
	colc columnClause
	labc string
	infr predictClause
	anal analyzeClause
}

const SELECT = 57346
const FROM = 57347
const WHERE = 57348
const LIMIT = 57349
const TRAIN = 57350
const PREDICT = 57351
const ANALYZE = 57352
const WITH = 57353
const COLUMN = 57354
const LABEL = 57355
const USING = 57356
const INTO = 57357
const FOR = 57358
const AS = 57359
const IDENT = 57360
const NUMBER = 57361
const STRING = 57362
const AND = 57363
const OR = 57364
const GE = 57365
const LE = 57366
const NE = 57367
const NOT = 57368
const POWER = 57369
const UMINUS = 57370

var sqlToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SELECT",
	"FROM",
	"WHERE",
	"LIMIT",
	"TRAIN",
	"PREDICT",
	"ANALYZE",
	"WITH",
	"COLUMN",
	"LABEL",
	"USING",
	"INTO",
	"FOR",
	"AS",
	"IDENT",
	"NUMBER",
	"STRING",
	"AND",
	"OR",
	"'>'",
	"'<'",
	"'='",
	"'!'",
	"GE",
	"LE",
	"NE",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"NOT",
	"POWER",
	"UMINUS",
	"';'",
	"'('",
	"')'",
	"','",
	"'['",
	"']'",
	"'\"'",
	"'\\''",
}
var sqlStatenames = [...]string{}

const sqlEofCode = 1
const sqlErrCode = 2
const sqlInitialStackSize = 16

//line sql.y:307

/* Like Lisp's builtin function cdr. */
func (e *expr) cdr() (r []string) {
	for i := 1; i < len(e.sexp); i++ {
		r = append(r, e.sexp[i].String())
	}
	return r
}

/* Convert exprlist to string slice. */
func (el exprlist) Strings() (r []string) {
	for i := 0; i < len(el); i++ {
		r = append(r, el[i].String())
	}
	return r
}

func (e *expr) String() string {
	if e.typ == 0 { /* a compound expression */
		switch e.sexp[0].typ {
		case '+', '*', '/', '%', '=', '<', '>', '!', LE, GE, AND, OR:
			if len(e.sexp) != 3 {
				log.Panicf("Expecting binary expression, got %.10q", e.sexp)
			}
			return fmt.Sprintf("%s %s %s", e.sexp[1], e.sexp[0].val, e.sexp[2])
		case '-':
			switch len(e.sexp) {
			case 2:
				return fmt.Sprintf(" -%s", e.sexp[1])
			case 3:
				return fmt.Sprintf("%s - %s", e.sexp[1], e.sexp[2])
			default:
				log.Panicf("Expecting either unary or binary -, got %.10q", e.sexp)
			}
		case '(':
			if len(e.sexp) != 2 {
				log.Panicf("Expecting ( ) as unary operator, got %.10q", e.sexp)
			}
			return fmt.Sprintf("(%s)", e.sexp[1])
		case '[':
			return "[" + strings.Join(e.cdr(), ", ") + "]"
		case NOT:
			return fmt.Sprintf("NOT %s", e.sexp[1])
		case IDENT: /* function call */
			return e.sexp[0].val + "(" + strings.Join(e.cdr(), ", ") + ")"
		}
	} else {
		return fmt.Sprintf("%s", e.val)
	}

	log.Panicf("Cannot print an unknown expression")
	return ""
}

func (s standardSelect) String() string {
	r := "SELECT "
	if len(s.fields) == 0 {
		r += "*"
	} else {
		for i := 0; i < len(s.fields); i++ {
			r += s.fields[i].String()
			if i != len(s.fields)-1 {
				r += ", "
			}
		}
	}
	r += "\nFROM " + strings.Join(s.tables, ", ")
	if s.where != nil {
		r += fmt.Sprintf("\nWHERE %s", s.where)
	}
	if len(s.limit) > 0 {
		r += fmt.Sprintf("\nLIMIT %s", s.limit)
	}
	return r
}

// sqlReentrantParser makes sqlParser, generated by goyacc and using a
// global variable parseResult to return the result, reentrant.
type sqlSyncParser struct {
	pr sqlParser
}

func newParser() *sqlSyncParser {
	return &sqlSyncParser{sqlNewParser()}
}

var mu sync.Mutex

func (p *sqlSyncParser) Parse(s string) (r *extendedSelect, e error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			e, ok = r.(error)
			if !ok {
				e = fmt.Errorf("%v", r)
			}
		}
	}()

	mu.Lock()
	defer mu.Unlock()

	p.pr.Parse(newLexer(s))
	return parseResult, nil
}

//line yacctab:1
var sqlExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const sqlPrivate = 57344

const sqlLast = 191

var sqlAct = [...]int{

	33, 116, 64, 115, 15, 94, 93, 90, 27, 26,
	28, 89, 91, 25, 92, 109, 91, 63, 43, 108,
	131, 35, 101, 91, 41, 34, 42, 70, 106, 30,
	72, 21, 36, 57, 31, 32, 128, 60, 61, 8,
	10, 9, 11, 12, 13, 129, 104, 75, 76, 77,
	78, 79, 80, 81, 82, 83, 84, 85, 86, 87,
	73, 129, 20, 19, 105, 27, 26, 28, 107, 46,
	47, 48, 4, 118, 97, 126, 99, 127, 35, 59,
	58, 24, 34, 17, 134, 132, 30, 117, 100, 36,
	62, 31, 32, 104, 102, 104, 130, 18, 27, 26,
	28, 44, 45, 46, 47, 48, 119, 124, 122, 120,
	114, 35, 121, 119, 95, 34, 125, 98, 96, 30,
	74, 71, 36, 39, 31, 32, 38, 37, 23, 40,
	119, 133, 55, 56, 51, 50, 49, 123, 53, 52,
	54, 44, 45, 46, 47, 48, 55, 56, 51, 50,
	49, 88, 53, 52, 54, 44, 45, 46, 47, 48,
	51, 50, 49, 65, 53, 52, 54, 44, 45, 46,
	47, 48, 112, 113, 69, 111, 67, 68, 3, 66,
	14, 29, 22, 16, 7, 6, 110, 103, 5, 2,
	1,
}
var sqlPact = [...]int{

	174, -1000, 34, 65, -1000, 25, 24, -7, 110, 62,
	80, 109, 108, 105, -1000, 112, -17, -13, -1000, -1000,
	-1000, -1000, -23, -1000, -1000, 125, -1000, -13, -1000, -1000,
	80, 60, 59, -1000, 80, 80, 47, 152, 165, 163,
	-12, 103, -10, 102, 80, 80, 80, 80, 80, 80,
	80, 80, 80, 80, 80, 80, 80, 111, -33, -38,
	-1000, -1000, -1000, -29, 125, 96, 100, 96, 99, 96,
	80, -1000, -1000, -18, -1000, 37, 37, -1000, -1000, -1000,
	71, 71, 71, 71, 71, 71, 137, 137, -1000, -1000,
	-1000, 80, -1000, 52, -1000, 3, -1000, 54, -1000, 5,
	-25, -1000, 125, 160, 96, 55, 80, 94, 90, -1000,
	122, 89, 55, 57, -1000, 20, -1000, -1000, -13, -1000,
	125, -1000, -1000, 78, -1000, 4, -1000, -1000, 67, 55,
	-1000, 66, -1000, -1000, -1000,
}
var sqlPgo = [...]int{

	0, 190, 189, 188, 187, 186, 185, 184, 183, 182,
	2, 0, 1, 17, 181, 3, 180, 5, 6,
}
var sqlR1 = [...]int{

	0, 1, 1, 1, 1, 2, 2, 2, 2, 3,
	3, 6, 6, 7, 7, 4, 4, 4, 16, 16,
	8, 8, 8, 12, 12, 12, 15, 15, 5, 5,
	9, 9, 17, 18, 18, 11, 11, 13, 13, 14,
	14, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10,
}
var sqlR2 = [...]int{

	0, 2, 3, 3, 3, 2, 3, 3, 3, 8,
	7, 4, 6, 4, 6, 2, 4, 5, 5, 1,
	1, 1, 3, 1, 1, 1, 1, 3, 2, 2,
	1, 3, 3, 1, 3, 3, 4, 1, 3, 2,
	3, 1, 1, 1, 1, 3, 3, 3, 1, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 2, 2,
}
var sqlChk = [...]int{

	-1000, -1, -2, 4, 38, -3, -6, -7, 5, 7,
	6, 8, 9, 10, -16, -11, -8, 18, 32, 38,
	38, 38, -9, 18, 19, -10, 19, 18, 20, -14,
	39, 44, 45, -11, 35, 31, 42, 18, 18, 18,
	17, 41, 39, 41, 30, 31, 32, 33, 34, 25,
	24, 23, 28, 27, 29, 21, 22, -10, 20, 20,
	-10, -10, 43, -13, -10, 11, 14, 11, 14, 11,
	39, 18, 40, -13, 18, -10, -10, -10, -10, -10,
	-10, -10, -10, -10, -10, -10, -10, -10, 40, 44,
	45, 41, 43, -18, -17, 18, 18, -18, 18, -18,
	-13, 40, -10, -4, 41, 12, 25, 14, 14, 40,
	-5, 15, 12, 13, -17, -15, -12, 32, 18, -11,
	-10, 18, 18, 15, 18, -15, 18, 20, 16, 41,
	18, 16, 18, -12, 18,
}
var sqlDef = [...]int{

	0, -2, 0, 0, 1, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 5, 0, 19, 21, 20, 2,
	3, 4, 6, 30, 7, 8, 41, 42, 43, 44,
	0, 0, 0, 48, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	62, 63, 39, 0, 37, 0, 0, 0, 0, 0,
	0, 22, 35, 0, 31, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 58, 59, 60, 61, 45, 46,
	47, 0, 40, 0, 33, 0, 11, 0, 13, 0,
	0, 36, 38, 0, 0, 0, 0, 0, 0, 18,
	0, 0, 0, 0, 34, 15, 26, 23, 24, 25,
	32, 12, 14, 0, 10, 0, 28, 29, 0, 0,
	9, 0, 16, 27, 17,
}
var sqlTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 26, 44, 3, 3, 34, 3, 45,
	39, 40, 32, 30, 41, 31, 3, 33, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 38,
	24, 25, 23, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 42, 3, 43,
}
var sqlTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 27, 28, 29, 35, 36, 37,
}
var sqlTok3 = [...]int{
	0,
}

var sqlErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	sqlDebug        = 0
	sqlErrorVerbose = false
)

type sqlLexer interface {
	Lex(lval *sqlSymType) int
	Error(s string)
}

type sqlParser interface {
	Parse(sqlLexer) int
	Lookahead() int
}

type sqlParserImpl struct {
	lval  sqlSymType
	stack [sqlInitialStackSize]sqlSymType
	char  int
}

func (p *sqlParserImpl) Lookahead() int {
	return p.char
}

func sqlNewParser() sqlParser {
	return &sqlParserImpl{}
}

const sqlFlag = -1000

func sqlTokname(c int) string {
	if c >= 1 && c-1 < len(sqlToknames) {
		if sqlToknames[c-1] != "" {
			return sqlToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func sqlStatname(s int) string {
	if s >= 0 && s < len(sqlStatenames) {
		if sqlStatenames[s] != "" {
			return sqlStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func sqlErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !sqlErrorVerbose {
		return "syntax error"
	}

	for _, e := range sqlErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + sqlTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := sqlPact[state]
	for tok := TOKSTART; tok-1 < len(sqlToknames); tok++ {
		if n := base + tok; n >= 0 && n < sqlLast && sqlChk[sqlAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if sqlDef[state] == -2 {
		i := 0
		for sqlExca[i] != -1 || sqlExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; sqlExca[i] >= 0; i += 2 {
			tok := sqlExca[i]
			if tok < TOKSTART || sqlExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if sqlExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += sqlTokname(tok)
	}
	return res
}

func sqllex1(lex sqlLexer, lval *sqlSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = sqlTok1[0]
		goto out
	}
	if char < len(sqlTok1) {
		token = sqlTok1[char]
		goto out
	}
	if char >= sqlPrivate {
		if char < sqlPrivate+len(sqlTok2) {
			token = sqlTok2[char-sqlPrivate]
			goto out
		}
	}
	for i := 0; i < len(sqlTok3); i += 2 {
		token = sqlTok3[i+0]
		if token == char {
			token = sqlTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = sqlTok2[1] /* unknown char */
	}
	if sqlDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", sqlTokname(token), uint(char))
	}
	return char, token
}

func sqlParse(sqllex sqlLexer) int {
	return sqlNewParser().Parse(sqllex)
}

func (sqlrcvr *sqlParserImpl) Parse(sqllex sqlLexer) int {
	var sqln int
	var sqlVAL sqlSymType
	var sqlDollar []sqlSymType
	_ = sqlDollar // silence set and not used
	sqlS := sqlrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	sqlstate := 0
	sqlrcvr.char = -1
	sqltoken := -1 // sqlrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		sqlstate = -1
		sqlrcvr.char = -1
		sqltoken = -1
	}()
	sqlp := -1
	goto sqlstack

ret0:
	return 0

ret1:
	return 1

sqlstack:
	/* put a state and value onto the stack */
	if sqlDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", sqlTokname(sqltoken), sqlStatname(sqlstate))
	}

	sqlp++
	if sqlp >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlS[sqlp] = sqlVAL
	sqlS[sqlp].yys = sqlstate

sqlnewstate:
	sqln = sqlPact[sqlstate]
	if sqln <= sqlFlag {
		goto sqldefault /* simple state */
	}
	if sqlrcvr.char < 0 {
		sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
	}
	sqln += sqltoken
	if sqln < 0 || sqln >= sqlLast {
		goto sqldefault
	}
	sqln = sqlAct[sqln]
	if sqlChk[sqln] == sqltoken { /* valid shift */
		sqlrcvr.char = -1
		sqltoken = -1
		sqlVAL = sqlrcvr.lval
		sqlstate = sqln
		if Errflag > 0 {
			Errflag--
		}
		goto sqlstack
	}

sqldefault:
	/* default state action */
	sqln = sqlDef[sqlstate]
	if sqln == -2 {
		if sqlrcvr.char < 0 {
			sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if sqlExca[xi+0] == -1 && sqlExca[xi+1] == sqlstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			sqln = sqlExca[xi+0]
			if sqln < 0 || sqln == sqltoken {
				break
			}
		}
		sqln = sqlExca[xi+1]
		if sqln < 0 {
			goto ret0
		}
	}
	if sqln == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			sqllex.Error(sqlErrorMessage(sqlstate, sqltoken))
			Nerrs++
			if sqlDebug >= 1 {
				__yyfmt__.Printf("%s", sqlStatname(sqlstate))
				__yyfmt__.Printf(" saw %s\n", sqlTokname(sqltoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for sqlp >= 0 {
				sqln = sqlPact[sqlS[sqlp].yys] + sqlErrCode
				if sqln >= 0 && sqln < sqlLast {
					sqlstate = sqlAct[sqln] /* simulate a shift of "error" */
					if sqlChk[sqlstate] == sqlErrCode {
						goto sqlstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if sqlDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", sqlS[sqlp].yys)
				}
				sqlp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if sqlDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", sqlTokname(sqltoken))
			}
			if sqltoken == sqlEofCode {
				goto ret1
			}
			sqlrcvr.char = -1
			sqltoken = -1
			goto sqlnewstate /* try again in the same state */
		}
	}

	/* reduction by production sqln */
	if sqlDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", sqln, sqlStatname(sqlstate))
	}

	sqlnt := sqln
	sqlpt := sqlp
	_ = sqlpt // guard against "declared and not used"

	sqlp -= sqlR2[sqln]
	// sqlp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if sqlp+1 >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlVAL = sqlS[sqlp+1]

	/* consult goto table to find next state */
	sqln = sqlR1[sqln]
	sqlg := sqlPgo[sqln]
	sqlj := sqlg + sqlS[sqlp].yys + 1

	if sqlj >= sqlLast {
		sqlstate = sqlAct[sqlg]
	} else {
		sqlstate = sqlAct[sqlj]
		if sqlChk[sqlstate] != -sqln {
			sqlstate = sqlAct[sqlg]
		}
	}
	// dummy call; replaced with literal code
	switch sqlnt {

	case 1:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:155
		{
			parseResult = &extendedSelect{
				extended:       false,
				standardSelect: sqlDollar[1].slct}
		}
	case 2:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:160
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          true,
				standardSelect: sqlDollar[1].slct,
				trainClause:    sqlDollar[2].tran}
		}
	case 3:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:167
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          false,
				standardSelect: sqlDollar[1].slct,
				predictClause:  sqlDollar[2].infr}
		}
	case 4:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:174
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          false,
				analyze:        true,
				standardSelect: sqlDollar[1].slct,
				analyzeClause:  sqlDollar[2].anal}
		}
	case 5:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:185
		{
			sqlVAL.slct.fields = sqlDollar[2].expl
		}
	case 6:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:186
		{
			sqlVAL.slct.tables = sqlDollar[3].tbls
		}
	case 7:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:187
		{
			sqlVAL.slct.limit = sqlDollar[3].val
		}
	case 8:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:188
		{
			sqlVAL.slct.where = sqlDollar[3].expr
		}
	case 9:
		sqlDollar = sqlS[sqlpt-8 : sqlpt+1]
//line sql.y:192
		{
			sqlVAL.tran.estimator = sqlDollar[2].val
			sqlVAL.tran.trainAttrs = sqlDollar[4].atrs
			sqlVAL.tran.columns = sqlDollar[5].colc
			sqlVAL.tran.label = sqlDollar[6].labc
			sqlVAL.tran.save = sqlDollar[8].val
		}
	case 10:
		sqlDollar = sqlS[sqlpt-7 : sqlpt+1]
//line sql.y:199
		{
			sqlVAL.tran.estimator = sqlDollar[2].val
			sqlVAL.tran.trainAttrs = sqlDollar[4].atrs
			sqlVAL.tran.columns = sqlDollar[5].colc
			sqlVAL.tran.save = sqlDollar[7].val
		}
	case 11:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:208
		{
			sqlVAL.infr.into = sqlDollar[2].val
			sqlVAL.infr.model = sqlDollar[4].val
		}
	case 12:
		sqlDollar = sqlS[sqlpt-6 : sqlpt+1]
//line sql.y:209
		{
			sqlVAL.infr.into = sqlDollar[2].val
			sqlVAL.infr.predAttrs = sqlDollar[4].atrs
			sqlVAL.infr.model = sqlDollar[6].val
		}
	case 13:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:213
		{
			sqlVAL.anal.trainedModel = sqlDollar[2].val
			sqlVAL.anal.explainer = sqlDollar[4].val
		}
	case 14:
		sqlDollar = sqlS[sqlpt-6 : sqlpt+1]
//line sql.y:214
		{
			sqlVAL.anal.trainedModel = sqlDollar[2].val
			sqlVAL.anal.analyzeAttrs = sqlDollar[4].atrs
			sqlVAL.anal.explainer = sqlDollar[6].val
		}
	case 15:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:218
		{
			sqlVAL.colc = map[string]exprlist{"feature_columns": sqlDollar[2].expl}
		}
	case 16:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:219
		{
			sqlVAL.colc = map[string]exprlist{sqlDollar[4].val: sqlDollar[2].expl}
		}
	case 17:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:220
		{
			sqlVAL.colc[sqlDollar[5].val] = sqlDollar[3].expl
		}
	case 18:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:224
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr, atomic(IDENT, "AS"), funcall("", sqlDollar[4].expl)}
		}
	case 19:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:227
		{
			sqlVAL.expl = sqlDollar[1].flds
		}
	case 20:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:231
		{
			sqlVAL.flds = append(sqlVAL.flds, atomic(IDENT, "*"))
		}
	case 21:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:232
		{
			sqlVAL.flds = append(sqlVAL.flds, atomic(IDENT, sqlDollar[1].val))
		}
	case 22:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:233
		{
			sqlVAL.flds = append(sqlDollar[1].flds, atomic(IDENT, sqlDollar[3].val))
		}
	case 23:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:237
		{
			sqlVAL.expr = atomic(IDENT, "*")
		}
	case 24:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:238
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 25:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:239
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 26:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:243
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 27:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:244
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 28:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:248
		{
			sqlVAL.labc = sqlDollar[2].val
		}
	case 29:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:249
		{
			sqlVAL.labc = sqlDollar[2].val[1 : len(sqlDollar[2].val)-1]
		}
	case 30:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:253
		{
			sqlVAL.tbls = []string{sqlDollar[1].val}
		}
	case 31:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:254
		{
			sqlVAL.tbls = append(sqlDollar[1].tbls, sqlDollar[3].val)
		}
	case 32:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:258
		{
			sqlVAL.atrs = attrs{sqlDollar[1].val: sqlDollar[3].expr}
		}
	case 33:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:262
		{
			sqlVAL.atrs = sqlDollar[1].atrs
		}
	case 34:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:263
		{
			sqlVAL.atrs = attrsUnion(sqlDollar[1].atrs, sqlDollar[3].atrs)
		}
	case 35:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:267
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, nil)
		}
	case 36:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:268
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, sqlDollar[3].expl)
		}
	case 37:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:272
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 38:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:273
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 39:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:277
		{
			sqlVAL.expl = nil
		}
	case 40:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:278
		{
			sqlVAL.expl = sqlDollar[2].expl
		}
	case 41:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:282
		{
			sqlVAL.expr = atomic(NUMBER, sqlDollar[1].val)
		}
	case 42:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:283
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 43:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:284
		{
			sqlVAL.expr = atomic(STRING, sqlDollar[1].val)
		}
	case 44:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:285
		{
			sqlVAL.expr = variadic('[', "square", sqlDollar[1].expl)
		}
	case 45:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:286
		{
			sqlVAL.expr = unary('(', "paren", sqlDollar[2].expr)
		}
	case 46:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:287
		{
			sqlVAL.expr = unary('"', "quota", atomic(STRING, sqlDollar[2].val))
		}
	case 47:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:288
		{
			sqlVAL.expr = unary('\'', "quota", atomic(STRING, sqlDollar[2].val))
		}
	case 48:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:289
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 49:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:290
		{
			sqlVAL.expr = binary('+', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 50:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:291
		{
			sqlVAL.expr = binary('-', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 51:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:292
		{
			sqlVAL.expr = binary('*', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 52:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:293
		{
			sqlVAL.expr = binary('/', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 53:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:294
		{
			sqlVAL.expr = binary('%', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 54:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:295
		{
			sqlVAL.expr = binary('=', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 55:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:296
		{
			sqlVAL.expr = binary('<', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 56:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:297
		{
			sqlVAL.expr = binary('>', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 57:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:298
		{
			sqlVAL.expr = binary(LE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 58:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:299
		{
			sqlVAL.expr = binary(GE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 59:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:300
		{
			sqlVAL.expr = binary(NE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 60:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:301
		{
			sqlVAL.expr = binary(AND, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 61:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:302
		{
			sqlVAL.expr = binary(OR, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 62:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:303
		{
			sqlVAL.expr = unary(NOT, sqlDollar[1].val, sqlDollar[2].expr)
		}
	case 63:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:304
		{
			sqlVAL.expr = unary('-', sqlDollar[1].val, sqlDollar[2].expr)
		}
	}
	goto sqlstack /* stack new state and value */
}