package internal

import (
	"sync"

	"github.com/dlclark/regexp2"
	sitter "github.com/smacker/go-tree-sitter"
)

var JavaInputRegex, _ = CompileRegex2("((RequestBody|getParameter|PathVariable|RequestParam|RequestHeader|CookieValue|ModelAttribute|getQuery|getHeader|getCookie|getPathVariable|getRequestBody|getRequestParam))")
var CSharpInputRegex, _ = CompileRegex2("((FromBody|FromQuery|FromRoute|FromHeader|FromForm))")
var GoInputRegex, _ = CompileRegex2("((Request|FormValue|FormFile|Get|Post|Put|Delete|Patch|Head|Options|Param|QueryParam|GetQuery|GetForm|GetFile|GetHeader|GetParam|GetPath|GetPostForm|GetPostFormFile|GetPostFormValue|GetPostFormFileValue|GetPostFormFileValues|GetPostFormValue|GetPostFormValues|GetPostFormFileValues|GetPostFormFileValue))")
var PythonInputRegex, _ = CompileRegex2("(request\\.GET|request\\.POST|request\\.FILES|request\\.data|request\\.form|request\\.values|request\\.cookies|request\\.headers|request\\.args|request\\.json|flask\\.request|django\\.HttpRequest|QueryDict\\.get|QueryDict\\.getlist|request\\.query_params|request\\.data|request\\.FILES\\['.*?'\\]|request\\.POST\\['.*?'\\]|request.FILES\\['.*?'\\]|request.COOKIES\\['.*?'\\]|request.GET\\['.*?'\\])")
var JavaScriptInputRegex, _ = CompileRegex2("(req\\.query|req\\.params|req\\.body|req\\.cookies|req\\.headers|req\\.get|req\\.post|req\\.file|req\\.files|req\\.param|req\\.queryParam|req\\.queryParams)")
var isParserEnabled = true

var isParallelFlag = true

var NumOfFiles = 0
var NumOfFilesProcessed = 0
var LineCounter = 0
var IsFinised = false
var RegForCheckingSimpleRegex *regexp2.Regexp

var SExpressionQuery = make(map[int]*sitter.Query)
var mutex sync.Mutex
var initExpressionMutex sync.Mutex
var Queries []Query

var Verbose = false
var RuleIdToEvaluate = 0
var Similarity = false
var ListRules = false
var DefaultResultsFile = "./results.json"

var ResultFormat = "csv"

type MatchDT struct {
	MatchVal      string
	Line          int
	StartPosition Position
	EndPosition   Position
	isTainted     bool
}

type OneFileData struct {
	FileName            string
	SourceCode          string
	SourceRunes         []rune
	LowerCaseSourceCode string
	SymbolTableDS       SymbolTable
	Language            int
}

type SymbolTable struct {
	StringLiterals          map[string][]LinePragma
	ObjectCreateExprs       map[string][]LinePragma
	MethodInvocations       map[string][]LinePragma
	Parameters              map[string][]LinePragma
	VarWithStringLiteralVal map[string][]Value
	VariablesTypes          map[string][]Value
	MethodDeclaration       map[string][]LinePragma
	MethodDefinition        map[string][]LinePragma
	Comments                map[string][]LinePragma
	IsInitialized           bool

	TaintedVariables map[string][]LinePragma
}

type Value struct {
	Val        string
	LinePragma LinePragma
}

type LinePragma struct {
	StartLine   uint32
	StartColumn uint32
	EndLine     uint32
	EndColumn   uint32
	FileName    string
}

type PatternComponent struct {
	Tag                string          `json:"tag"`
	Value              string          `json:"value"`
	Name               string          `json:"name,omitempty"`
	ValueCompiledRegex *regexp2.Regexp `json:"-"`
}

type DetectionSequence struct {
	PatternSequenceArr []PatternComponent `json:"detectionSequence"`
}

type Query struct {
	Id                string              `json:"id"`
	RuleId            int                 `json:"ruleId"`
	Language          string              `json:"language"`
	FileExtensions    string              `json:"fileExtensions"`
	RuleName          string              `json:"ruleName"`
	Severity          string              `json:"severity"`
	Description       string              `json:"description"`
	RemediationAdvice string              `json:"remediationAdvice"`
	TriggerPattern    string              `json:"triggerPattern"`
	DetectionPatterns []DetectionSequence `json:"detectionPatterns"`
	ExclusionPattern  ExclusionPattern    `json:"exclusionPattern"`

	TriggerPatternRegex *regexp2.Regexp `json:"-"`
	ExclusionPatternRex *regexp2.Regexp `json:"-"`

	TriggerExpr   ParsedSimpleRegex `json:"-"`
	ExclusionExpr ParsedSimpleRegex `json:"-"`
}

type ExclusionPattern struct {
	Scope string `json:"scope"`
	Range string `json:"range"`
	Value string `json:"value"`
}

type ParsedSimpleRegex struct {
	IsCaseInsensitive bool
	IsComplexRegex    bool
	PatternArr        []string
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Result struct {
	RuleId            int      `json:"rule_id"`
	Language          string   `json:"language"`
	RuleName          string   `json:"rule_name"`
	FileName          string   `json:"file_name"`
	Line              string   `json:"line"`
	ProblematicLine   string   `json:"problematic_line"`
	StartPosition     Position `json:"start_position"`
	EndPosition       Position `json:"end_position"`
	RemediationAdvice string   `json:"remediation_advice"`
	SimilarityID      string   `json:"similarity_id,omitempty"`
	Severity          string   `json:"severity"`
	Description       string   `json:"description"`
}

type Macro struct {
	Name   string              `json:"name"`
	Macros []DetectionSequence `json:"macros"`
}
