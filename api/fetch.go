package main

import (
	"cmp"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"iter"
	"log"
	"maps"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/MarkRosemaker/cassette"
	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/openapi"
	codegen "github.com/MarkRosemaker/openapi-codegen"
	compress "github.com/MarkRosemaker/openapi-compress"
	edit "github.com/MarkRosemaker/openapi-edit"
	enrich "github.com/MarkRosemaker/openapi-enrich"
	flatten "github.com/MarkRosemaker/openapi-flatten"
	"github.com/ettle/strcase"
	"golang.org/x/sync/errgroup"
)

const (
	pathOpenAPI     = "api/openapi.json"
	schemaRefPrefix = "#/components/schemas/"
)

var isClaudeCode = os.Getenv("CLAUDECODE") != ""

func main() {
	eg := errgroup.Group{}

	ctx := context.Background()
	eg.Go(func() error { return fetchLLMs(ctx) })
	eg.Go(func() error { return persistOpenAPI(ctx) })

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	doc, err := fixOpenAPI()
	if err != nil {
		log.Fatal(err)
	}

	if err := generateCode(doc); err != nil {
		log.Fatal(err)
	}
}

func fetchLLMs(ctx context.Context) error {
	if isClaudeCode {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/llms.txt", nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	f, err := os.Create("api/llms.txt")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, rsp.Body); err != nil {
		return err
	}

	return nil
}

func readOfficial(ctx context.Context) (io.ReadCloser, error) {
	if isClaudeCode {
		return os.Open("api/openapi-official.json")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://developers.notion.com/openapi.json", nil)
	if err != nil {
		return nil, err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

func fetchOpenAPI(ctx context.Context) (*openapi.Document, error) {
	rc, err := readOfficial(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return openapi.LoadFromReader(rc)
}

func persistOpenAPI(ctx context.Context) error {
	doc, err := fetchOpenAPI(ctx)
	if err != nil {
		return err
	}

	var w io.Writer

	f, err := os.Create(pathOpenAPI)
	if err != nil {
		return err
	}
	defer f.Close()

	if isClaudeCode {
		w = f
	} else {
		f2, err := os.Create("api/openapi-official.json")
		if err != nil {
			return err
		}
		defer f2.Close()

		w = io.MultiWriter(f, f2)
	}

	return doc.WriteJSON(w)
}

// shortNames names what has no better name than its position in a union or allOf, or one too long for what it is,
// short and as Notion calls it: by the prefix the names hold, for the schemas within it too, as
// MeetingNotesDateFilterOperator for MeetingNotesPropertyFilterFilterOneOf2Operator. A name in exactNames is renamed
// whole and goes before a prefix. An entry that names no schema fails the run.
var shortNames = map[string]string{
	"AgentListResultsItemOneOf0":   "CustomAgent",
	"AgentListResultsItemNotionAi": "NotionAIAgent",

	"AgentListResultsItemOneOf0CreatedByOneOf0": "AgentCreator", // the user or bot that created it

	"CreateACommentAllOf0OneOf0": "PageComment",
	"CreateACommentAllOf0OneOf1": "PageMarkdownComment",
	"CreateACommentAllOf0OneOf2": "DiscussionComment",
	"CreateACommentAllOf0OneOf3": "DiscussionMarkdownComment",
	"CreateACommentAllOf0":       "NewComment",
	"UpdateACommentAllOf0":       "CommentUpdate",
	"UpdatePageMarkdownAllOf0":   "MarkdownEdit",

	"CreateMeetingNoteAllOf0OneOf0": "MeetingNoteInPage",
	"CreateMeetingNoteAllOf0Source": "MeetingNoteFromBlock",
	"CreateMeetingNoteAllOf0":       "MeetingNoteSource",
	"CreateMeetingNoteOkOneOf0":     "MeetingNoteBlock",

	"MeetingNotesPropertyFilterFilterString": "MeetingNotesTextFilter",
	"MeetingNotesPropertyFilterFilterPerson": "MeetingNotesPersonFilter",
	"MeetingNotesPropertyFilterFilterOneOf2": "MeetingNotesDateFilter",
	"MeetingNotesPropertyFilterFilterOneOf3": "MeetingNotesDateRangeFilter",

	"MeetingNotesPropertyFilterFilterOneOf3ValueOneOf2": "MeetingNotesDateRangeFilterRelative",
	"QueryMeetingNotesFilterFiltersItemOneOf1":          "MeetingNotesFilterGroup",

	"DataSourceObjectSyncMetadataHealthAnyOf":  "SyncHealth", // what the failing and the healthy share
	"DataSourceObjectSyncMetadataHealthPaused": "SyncHealthPaused",

	"RollupPropertyConfigurationRollupAllOfOneOf":   "RollupByNames",
	"RollupPropertyConfigurationRollupAllOf0OneOf1": "RollupByRelationID",
	"RollupPropertyConfigurationRollupAllOf0OneOf2": "RollupByRollupID",
	"RollupPropertyConfigurationRollupAllOf0OneOf3": "RollupByIDs",

	"UpdateABlockAnyOf30":                    "BlockTrashUpdate",
	"UpdateADataSourcePropertiesValueOneOf0": "PropertyUpdate",
	"UpdateADataSourcePropertiesValueName":   "PropertyRename",
	"UpdateAgentStatus2":                     "AgentStatusUpdate", // beside the UpdateAgentStatus it sets
	"UpdateSessionRequestMessage":            "SessionMessage",
	"UpdateSessionRequestOneOf1":             "SessionAnswer",
	"UpdateSessionRequestOneOf2":             "SessionContinuation",
}

// exactNames are the names of shortNames whose schema alone is renamed, as their prefix starts the names of others.
var exactNames = map[string]string{
	"DataSourceObjectSyncMetadataHealthAnyOf":         "SyncHealthFailing",
	"DataSourceObjectResponseSyncMetadataHealthAnyOf": "SyncHealthOK",

	// upstream's responses, whose names differ from the requests fetch builds only in case, as PageIconResponse holds
	"contentWithRichTextColorAndIcon": "ContentWithRichTextColorAndIconResponse",
	"customEmojiPageIcon":             "CustomEmojiPageIconResponse",
	"externalPageIcon":                "ExternalPageIconResponse",
}

// shortenNames renames the component schemas as shortNames and exactNames say, then drops ObjectResponse or Response
// from a name, as User for userObjectResponse, and a word said twice in a row, as BlockObjectToDo for
// BlockObjectToDoToDo, where the name without it is no other schema's. A response that the request of the same name
// differs from keeps its suffix, as BoardViewConfigResponse beside BoardViewConfig: there it tells the two apart.
func shortenNames(doc *openapi.Document) error {
	unused := map[string]bool{}
	for k := range shortNames {
		unused[k] = true
	}

	for k := range exactNames {
		unused[k] = true
	}

	prefixes := slices.SortedFunc(maps.Keys(shortNames), func(a, b string) int { return len(b) - len(a) })

	byTable := func(n string) string {
		if to, ok := exactNames[n]; ok {
			delete(unused, n)
			return to
		}

		for _, prefix := range prefixes {
			if rest, ok := strings.CutPrefix(n, prefix); ok {
				delete(unused, prefix)
				return shortNames[prefix] + rest
			}
		}

		return n
	}

	names := slices.Collect(maps.Keys(doc.Components.Schemas))
	slices.Sort(names)

	taken := map[string]bool{} // by Go name, as codegen writes it
	for _, n := range names {
		taken[strcase.ToGoPascal(byTable(n))] = true
	}

	to := map[string]string{}
	claim := func(n, short string) bool {
		if taken[short] {
			return false
		}

		taken[short] = true
		to[n] = short

		return true
	}

	for _, n := range names {
		if t := byTable(n); t != n {
			to[n] = t
			continue
		}

		for _, suffix := range []string{"ObjectResponse", "Response"} {
			if before, ok := strings.CutSuffix(n, suffix); ok && before != "" && claim(n, strcase.ToGoPascal(before)) {
				break
			}
		}

		if _, ok := to[n]; !ok {
			if short := withoutRepeats(strcase.ToGoPascal(n)); short != strcase.ToGoPascal(n) {
				claim(n, short)
			}
		}
	}

	// flatten named what a schema holds after it, as BotInfoResponseOwner, which follows its new name, as BotInfoOwner
	renamed := slices.SortedFunc(maps.Keys(to), func(a, b string) int { return len(b) - len(a) })
	for _, n := range names {
		if _, ok := to[n]; ok {
			continue
		}

		goName := strcase.ToGoPascal(n)
		for _, parent := range renamed {
			rest, ok := strings.CutPrefix(goName, strcase.ToGoPascal(parent))
			if ok && rest != "" && unicode.IsUpper(rune(rest[0])) {
				claim(n, withoutRepeats(strcase.ToGoPascal(to[parent])+rest))
				break
			}
		}
	}

	if err := edit.RenameSchemas(doc, to); err != nil {
		return err
	}

	// codegen would number a second schema of the same Go name, as it did contentWithRichTextColorAndIcon2
	byGoName := map[string]string{}
	for n := range doc.Components.Schemas.ByIndex() {
		if other, ok := byGoName[strcase.ToGoPascal(n)]; ok {
			return fmt.Errorf("%s and %s would have the same Go name: name one in shortNames or exactNames", other, n)
		}

		byGoName[strcase.ToGoPascal(n)] = n
	}

	if len(unused) > 0 {
		return fmt.Errorf("shortNames or exactNames name schemas the specification no longer has: %s",
			strings.Join(slices.Sorted(maps.Keys(unused)), ", "))
	}

	return nil
}

// reWord matches a word of a Go name, as To, Do and URL.
var reWord = regexp.MustCompile(`[A-Z]+[a-z0-9]*|[a-z0-9]+`)

// withoutRepeats is the Go name with every run of words said twice in a row said once, as BlockObjectToDo for
// BlockObjectToDoToDo.
func withoutRepeats(name string) string {
	words := reWord.FindAllString(name, -1)

	for k := len(words) / 2; k >= 1; k-- {
		for i := 0; i+2*k <= len(words); i++ {
			if slices.Equal(words[i:i+k], words[i+k:i+2*k]) {
				words = slices.Delete(words, i+k, i+2*k)
				i--
			}
		}
	}

	return strings.Join(words, "")
}

func stripSuffixes(doc *openapi.Document) error {
	return edit.RenameSchemasFunc(doc, func(n string) string {
		if before, ok := strings.CutSuffix(n, "Request"); ok {
			return strcase.ToGoPascal(before)
		}

		return n
	})
}

// consolidateErrors replaces the schemas whose names start with any of prefixes, each an allOf of the common error and
// a code and status, with one schema called name that allows every code and status they did. Its code refers to a
// schema of its own, called codeName.
func consolidateErrors(doc *openapi.Document, prefixes []string, name, common, codeName string) error {
	var names []string

	for n := range doc.Components.Schemas.ByIndex() {
		if slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(n, p) }) {
			names = append(names, n)
		}
	}

	if len(names) == 0 {
		return fmt.Errorf("no schemas start with any of %q", prefixes)
	}

	if _, ok := doc.Components.Schemas[common]; !ok {
		return fmt.Errorf("schema %q not found", common)
	}

	codes, statuses := &jsonSet{}, &jsonSet{}

	for _, n := range names {
		if err := collectError(doc.Components.Schemas[n], common, codes, statuses); err != nil {
			return componentErr(n, err)
		}
	}

	codeSchema := &openapi.Schema{Type: openapi.TypeString, Enum: codes.values}
	doc.Components.Schemas.Set(codeName, codeSchema)

	// set one by one, since a map literal leaves their order to chance
	props := openapi.Schemas{}
	props.Set("code", edit.SchemaRef(doc, codeName))
	props.Set("status", &openapi.Schema{Type: openapi.TypeInteger, Enum: statuses.values})

	first := names[0]
	doc.Components.Schemas[first].Replace(&openapi.Schema{AllOf: openapi.SchemaList{
		edit.SchemaRef(doc, common),
		{Type: openapi.TypeObject, Properties: props, Required: []string{"code", "status"}},
	}})

	redirect := map[string]string{}
	for _, n := range names[1:] {
		redirect[n] = first
	}

	if err := edit.RedirectSchemas(doc, redirect); err != nil {
		return fmt.Errorf("redirecting schemas: %w", err)
	}

	if err := edit.RenameSchema(doc, first, name); err != nil {
		return fmt.Errorf("renaming %s to %s: %w", first, name, err)
	}

	return nil
}

// collectError adds the codes and statuses s allows to codes and statuses, failing on anything it would not keep.
func collectError(s *openapi.Schema, common string, codes, statuses *jsonSet) error {
	if len(s.OneOf) > 0 {
		for i, alt := range s.OneOf {
			if err := collectError(alt, common, codes, statuses); err != nil {
				return &errpath.ErrField{Field: "oneOf", Err: &errpath.ErrIndex{Index: i, Err: err}}
			}
		}

		return nil
	}

	if len(s.AllOf) != 2 || s.AllOf[0].Ref == nil || s.AllOf[0].Ref.Identifier != schemaRefPrefix+common {
		return &errpath.ErrField{Field: "allOf", Err: fmt.Errorf("want %s and one schema", common)}
	}

	for prop, p := range s.AllOf[1].Properties.ByIndex() {
		switch prop {
		case "code":
			codes.add(p)
		case "status":
			statuses.add(p)
		case "additional_data": // a narrower shape of what the common error already allows
		default:
			return &errpath.ErrField{Field: "allOf", Err: &errpath.ErrIndex{Index: 1, Err: &errpath.ErrField{
				Field: "properties", Err: &errpath.ErrKey{Key: prop, Err: errors.New("unexpected property")},
			}}}
		}
	}

	return nil
}

// nameEnum gives s's enum values the Go names names, which must match them one for one. It fails if s already has
// extensions, so a change upstream that adds some is noticed instead of overwritten.
// See https://github.com/oapi-codegen/oapi-codegen/blob/main/docs/extensions.md#x-enum-varnames--x-enumnames.
func nameEnum(s *openapi.Schema, names ...string) error {
	if len(s.Enum) != len(names) {
		return &errpath.ErrField{Field: "enum", Err: fmt.Errorf("has %d values, want %d", len(s.Enum), len(names))}
	}

	if len(s.Extensions) > 0 {
		return fmt.Errorf("already has extensions: %s", s.Extensions)
	}

	ext, err := json.Marshal(map[string][]string{"x-enum-varnames": names})
	if err != nil {
		return err
	}

	s.Extensions = ext

	return nil
}

// propertyErr reports err as one of the property prop of the component schema name.
func propertyErr(name, prop string, err error) error {
	return componentErr(name, &errpath.ErrField{Field: "properties", Err: &errpath.ErrKey{Key: prop, Err: err}})
}

// componentErr reports err as one of the component schema name.
func componentErr(name string, err error) error {
	return &errpath.ErrField{Field: "components", Err: &errpath.ErrField{
		Field: "schemas", Err: &errpath.ErrKey{Key: name, Err: err},
	}}
}

// namePartialUnions names each union of an object and its partial form, which responses spell out inline and flatten
// would name after the operation, or which the official spec names after where it is used, such as PageOrPartial for
// a page or the partial page Notion returns in its place. It is an anyOf, since an object matches its partial form too.
func namePartialUnions(doc *openapi.Document) error {
	var partials []string

	for name := range doc.Components.Schemas.ByIndex() {
		if strings.HasPrefix(name, "partial") && strings.HasSuffix(name, "ObjectResponse") {
			partials = append(partials, name)
		}
	}

	for _, name := range partials {
		ref := schemaRefPrefix + name
		union := strcase.ToGoPascal(strings.TrimSuffix(strings.TrimPrefix(name, "partial"), "ObjectResponse")) + "OrPartial"

		isUnion := func(s *openapi.Schema) bool {
			alts := alternatives(s)

			return len(alts) == 2 && alts[0].Ref != nil && alts[1].Ref != nil &&
				(alts[0].Ref.Identifier == ref || alts[1].Ref.Identifier == ref)
		}

		if err := edit.ExtractSchema(doc, union, isUnion); err != nil && !errors.Is(err, edit.ErrNoMatch) {
			return fmt.Errorf("naming the union with %s: %w", name, err)
		}

		// the official spec names some of them itself, such as userValueResponse
		var named []string

		for n, s := range doc.Components.Schemas.ByIndex() {
			if n != union && isUnion(s) {
				named = append(named, n)
			}
		}

		if _, ok := doc.Components.Schemas[union]; !ok && len(named) > 0 {
			if err := edit.RenameSchema(doc, named[0], union); err != nil {
				return err
			}

			named = named[1:]
		}

		for _, n := range named {
			if err := edit.RedirectSchema(doc, n, union, ""); err != nil {
				return err
			}
		}

		u, ok := doc.Components.Schemas[union]
		if !ok {
			continue
		}

		// some say oneOf, but an object matches its partial form too
		u.AnyOf, u.OneOf = alternatives(u), nil
	}

	return nil
}

// nameBranches moves each inline branch of a union into the component schemas, under its parent's
// name and what tells it apart from the other branches -- see [telling].
// flatten would name it after its title, which other unions' branches share, or its position, which says nothing.
func nameBranches(doc *openapi.Document) {
	for name, s := range doc.Components.Schemas.ByIndex() {
		nameBranchesIn(doc, s, name)
	}

	for s, name := range bodies(doc) {
		nameBranchesIn(doc, s, name)
	}
}

// bodies yields the schema of each request and response body, with the name flatten gives it.
func bodies(doc *openapi.Document) iter.Seq2[*openapi.Schema, string] {
	return func(yield func(*openapi.Schema, string) bool) {
		inContent := func(c openapi.Content, name string) bool {
			for _, mt := range c {
				if mt.Schema != nil && !yield(mt.Schema, strcase.ToGoPascal(cmp.Or(mt.Schema.Title, name))) {
					return false
				}
			}

			return true
		}

		for _, p := range doc.Paths.ByIndex() {
			for _, op := range p.Operations {
				if op.RequestBody != nil && op.RequestBody.Value != nil &&
					!inContent(op.RequestBody.Value.Content, op.OperationID) {
					return
				}

				for code, rsp := range op.Responses.ByIndex() {
					if rsp.Value != nil &&
						!inContent(rsp.Value.Content, op.OperationID+" "+cmp.Or(code.StatusText(), string(code))) {
						return
					}
				}
			}
		}
	}
}

// nameLists names each list Notion responds with after what it lists, such as BlockList, and gives it the request_id
// Notion sends with every response. Lists of the same are then one schema.
func nameLists(doc *openapi.Document) {
	for s := range bodies(doc) {
		object, _ := constString(s.Properties["object"])
		of, ok := constString(s.Properties["type"])
		if object != "list" || !ok || s.Ref != nil {
			continue
		}

		if _, ok := s.Properties["request_id"]; !ok {
			s.Properties.Set("request_id", requestID())
		}

		name := strcase.ToGoPascal(of) + "List"
		if same, ok := doc.Components.Schemas[name]; ok && equalJSON(same, s) {
			s.Replace(edit.SchemaRef(doc, name))
		} else {
			moveToComponents(doc, s, name)
		}
	}
}

// nameObjects names each object Notion responds with after what it is, such as Session. Of the bodies that are the
// same object, the one with the most properties takes the name, as the others are partial forms of it.
func nameObjects(doc *openapi.Document) {
	fullest := map[string]*openapi.Schema{}

	for s := range bodies(doc) {
		object, ok := constString(s.Properties["object"])
		if !ok || object == "list" || s.Ref != nil {
			continue
		}

		if f, ok := fullest[object]; !ok || len(s.Properties) > len(f.Properties) {
			fullest[object] = s
		}
	}

	for _, object := range slices.Sorted(maps.Keys(fullest)) {
		moveToComponents(doc, fullest[object], strcase.ToGoPascal(object))
	}
}

// setResults makes the results of a search or a query of a data source a list of PageOrDataSource.
func setResults(doc *openapi.Document) error {
	const list, item = "PageOrDataSourceList", "PageOrDataSource"

	if err := edit.RenameSchema(doc, "pageOrDataSourceListResponse", list); err != nil {
		return err
	}

	s, ok := doc.Components.Schemas[list]
	if !ok {
		return componentErr(list, errors.New("not found"))
	}

	results, ok := s.Properties["results"]
	if !ok {
		return propertyErr(list, "results", errors.New("not found"))
	}

	if _, ok := doc.Components.Schemas[item]; !ok {
		if results.Items == nil {
			return propertyErr(list, "results", errors.New("lists no items"))
		}

		moveToComponents(doc, results.Items, item)
	}

	results.Replace(&openapi.Schema{Type: openapi.TypeArray, Items: edit.SchemaRef(doc, item)})

	return nil
}

// requestID is the schema of the request_id Notion sends with every response.
func requestID() *openapi.Schema {
	return &openapi.Schema{Type: openapi.TypeString, Format: openapi.FormatUUID}
}

// equalJSON reports whether a and b are written the same.
func equalJSON(a, b *openapi.Schema) bool {
	ja, errA := json.Marshal(a)
	jb, errB := json.Marshal(b)

	return errA == nil && errB == nil && string(ja) == string(jb)
}

// nameBranchesIn is nameBranches for s, the schema flatten would name name, and what it holds.
func nameBranchesIn(doc *openapi.Document, s *openapi.Schema, name string) {
	if s == nil || s.Ref != nil {
		return
	}

	s.OneOf, s.AnyOf = spliceUnions(s.OneOf, true), spliceUnions(s.AnyOf, false)
	hoistParts(s)

	for _, union := range []struct {
		alts openapi.SchemaList
		kind string
	}{{s.OneOf, "OneOf"}, {s.AnyOf, "AnyOf"}} {
		for i, alt := range union.alts {
			branch := fmt.Sprintf("%s%s%d", name, union.kind, i)

			if v, ok := telling(union.alts, i); ok && alt.Type == openapi.TypeObject {
				if moved, n, ok := moveToComponents(doc, alt, strcase.ToGoPascal(name+" "+v)); ok {
					alt, branch = moved, n
				}
			}

			nameBranchesIn(doc, alt, branch)
		}
	}

	for i, part := range s.AllOf {
		// codegen folds an inline part into s, so what it holds is named as if s held it
		for prop, p := range part.Properties.ByIndex() {
			if p.Type == openapi.TypeObject || len(alternatives(p)) > 0 {
				moveToComponents(doc, p, strcase.ToGoPascal(name+" "+prop))
			}
		}

		nameBranchesIn(doc, part, fmt.Sprintf("%sAllOf%d", name, i))
	}

	for prop, p := range s.Properties.ByIndex() {
		nameBranchesIn(doc, p, strcase.ToGoPascal(name+" "+strings.ReplaceAll(prop, "/", " ")))
	}

	nameBranchesIn(doc, s.Items, name+"Item")

	if s.AdditionalProperties != nil {
		nameBranchesIn(doc, s.AdditionalProperties.Schema, name+"Value")
	}
}

// moveToComponents moves s, if inline, into the component schemas as name, if free, leaving a reference to it in its
// place, and returns the moved schema.
func moveToComponents(doc *openapi.Document, s *openapi.Schema, name string) (*openapi.Schema, string, bool) {
	if s.Ref != nil || doc.Components.Schemas[name] != nil {
		return nil, "", false
	}

	moved := new(openapi.Schema)
	moved.Replace(s)
	doc.Components.Schemas.Set(name, moved)
	s.Replace(edit.SchemaRef(doc, name))

	return moved, name, true
}

// hoistParts moves the properties of each inline part of s's allOf that is only an object with properties onto s
// itself, which allows the same values, so that they are named after s rather than after the part's position.
func hoistParts(s *openapi.Schema) {
	s.AllOf = slices.DeleteFunc(s.AllOf, func(part *openapi.Schema) bool {
		if part.Ref != nil || part.Type != openapi.TypeObject || len(part.Properties) == 0 ||
			len(alternatives(part)) > 0 || len(part.AllOf) > 0 || part.AdditionalProperties != nil {
			return false
		}

		for prop := range part.Properties {
			if _, ok := s.Properties[prop]; ok {
				return false
			}
		}

		if s.Properties == nil {
			s.Properties = openapi.Schemas{}
		}

		for prop, p := range part.Properties.ByIndex() {
			s.Properties.Set(prop, p)
		}

		s.Type, s.Required = openapi.TypeObject, append(s.Required, part.Required...)

		return true
	})
}

// spliceUnions lists, instead of each inline alternative in alts that is only a union of the same kind, that union's
// alternatives, which allows the same values.
func spliceUnions(alts openapi.SchemaList, oneOf bool) openapi.SchemaList {
	var spliced openapi.SchemaList

	for _, alt := range alts {
		inner := alt.AnyOf
		if oneOf {
			inner = alt.OneOf
		}

		if alt.Ref == nil && len(inner) > 0 && onlyUnion(alt) {
			spliced = append(spliced, spliceUnions(inner, oneOf)...)
		} else {
			spliced = append(spliced, alt)
		}
	}

	return spliced
}

// onlyUnion reports whether s says nothing besides its oneOf or anyOf and a description.
func onlyUnion(s *openapi.Schema) bool {
	c := *s
	c.OneOf, c.AnyOf, c.Description = nil, nil, ""

	b, err := json.Marshal(&c)

	return err == nil && string(b) == "{}"
}

// telling is what tells the branch alts[i] of a union apart from the others: a value it fixes a member to and no
// other branch does, the member type's first, or else the one member it requires.
func telling(alts openapi.SchemaList, i int) (string, bool) {
	s := alts[i]

	props := []string{"type"}
	for prop := range s.Properties.ByIndex() {
		props = append(props, prop)
	}

	for _, prop := range props {
		if v, ok := constString(s.Properties[prop]); ok && !slices.ContainsFunc(alts, func(other *openapi.Schema) bool {
			w, ok := constString(deref(other).Properties[prop])
			return other != s && ok && w == v
		}) {
			return v, true
		}
	}

	if len(s.Required) == 1 {
		return s.Required[0], true
	}

	// an enum whose values all begin alike, as string_is and string_contains, unlike those of every other alternative
	for _, prop := range props {
		if v := enumPrefix(s.Properties[prop]); v != "" && !slices.ContainsFunc(alts, func(other *openapi.Schema) bool {
			return other != s && enumPrefix(deref(other).Properties[prop]) == v
		}) {
			return v, true
		}
	}

	return "", false
}

// enumPrefix is the words, joined by underscores, that every value of p's string enum of two or more begins with, as
// string of string_is and string_contains, or "" if they share none.
func enumPrefix(p *openapi.Schema) string {
	if p == nil || len(deref(p).Enum) < 2 {
		return ""
	}

	var prefix []string

	for i, e := range deref(p).Enum {
		var v string
		if json.Unmarshal(e, &v) != nil {
			return ""
		}

		words := strings.Split(v, "_")
		if i == 0 {
			prefix = words
			continue
		}

		n := 0
		for n < len(prefix) && n < len(words) && prefix[n] == words[n] {
			n++
		}

		prefix = prefix[:n]
	}

	return strings.Join(prefix, "_")
}

// constString is the string p fixes its value to, if it does, as its const or the one value of its enum.
func constString(p *openapi.Schema) (string, bool) {
	if p == nil {
		return "", false
	}

	c := p.Const
	if len(p.Enum) == 1 {
		c = p.Enum[0]
	}

	var v string

	return v, json.Unmarshal(c, &v) == nil
}

// deref is the schema s stands for: the one it refers to, if it is a reference.
func deref(s *openapi.Schema) *openapi.Schema {
	if s.Ref != nil && s.Ref.Value != nil {
		return s.Ref.Value
	}

	return s
}

// variants are the names of the component schemas the union name lists.
func variants(doc *openapi.Document, name string) ([]string, error) {
	s, ok := doc.Components.Schemas[name]
	if !ok {
		return nil, componentErr(name, errors.New("not found"))
	}

	var names []string

	for _, alt := range alternatives(s) {
		n, ok := strings.CutPrefix(alt.Ref.Identifier, schemaRefPrefix)
		if alt.Ref == nil || !ok {
			return nil, componentErr(name, errors.New("lists a variant that is not a component"))
		}

		names = append(names, n)
	}

	return names, nil
}

// alternatives are s's oneOf, or else its anyOf.
func alternatives(s *openapi.Schema) openapi.SchemaList {
	if len(s.OneOf) > 0 {
		return s.OneOf
	}

	return s.AnyOf
}

// addRequestID adds the request_id Notion returns on every top-level object, but the official spec leaves out, to the
// schemas names.
func addRequestID(doc *openapi.Document, names ...string) error {
	for _, n := range names {
		s, ok := doc.Components.Schemas[n]
		if !ok {
			return componentErr(n, errors.New("not found"))
		}

		if _, ok := s.Properties["request_id"]; ok {
			return propertyErr(n, "request_id", errors.New("already there"))
		}

		s.Properties.Set("request_id", requestID())
	}

	return nil
}

// dateProps are the properties Notion declares as dates but takes and sends as a date or a date-time: "an ISO 8601
// date, with optional time".
var dateProps = []struct {
	schema string
	props  []string
}{
	{"Date", []string{"start", "end"}},
	{"dateResponse", []string{"start", "end"}},
}

// allowDateTimes makes the dateProps a DateOrDateTime. It fails unless each of them is, or is one of, a string with
// the date format, so a change upstream is noticed.
func allowDateTimes(doc *openapi.Document) error {
	const name = "DateOrDateTime"

	doc.Components.Schemas.Set(name, &openapi.Schema{
		Description: "An ISO 8601 date, with an optional time.",
		AnyOf: openapi.SchemaList{
			{Type: openapi.TypeString, Format: openapi.FormatDate},
			{Type: openapi.TypeString, Format: openapi.FormatDateTime},
		},
	})

	for _, d := range dateProps {
		s, ok := doc.Components.Schemas[d.schema]
		if !ok {
			return componentErr(d.schema, errors.New("not found"))
		}

		for _, prop := range d.props {
			p, ok := s.Properties[prop]
			if !ok {
				return propertyErr(d.schema, prop, errors.New("not found"))
			}

			found := false

			for _, alt := range append(openapi.SchemaList{p}, p.OneOf...) {
				if alt.Type == openapi.TypeString && alt.Format == openapi.FormatDate {
					ref := edit.SchemaRef(doc, name)
					ref.Description = alt.Description
					alt.Replace(ref)

					found = true
				}
			}

			if !found {
				return propertyErr(d.schema, prop, errors.New("is not a date"))
			}
		}
	}

	return nil
}

// jsonSet holds the distinct values of a const or enum, in the order first seen.
type jsonSet struct {
	seen   map[string]bool
	values []jsontext.Value
}

func (s *jsonSet) add(p *openapi.Schema) {
	if s.seen == nil {
		s.seen = map[string]bool{}
	}

	vals := p.Enum
	if p.Const != nil {
		vals = append(vals, p.Const)
	}

	for _, v := range vals {
		if !s.seen[string(v)] {
			s.seen[string(v)] = true
			s.values = append(s.values, v)
		}
	}
}

// addForm adds to the union a form for the kind, made from its form model for the kind like: the same members but
// for the type, which is the kind, and like's member, which is member under the kind's name. It is named as model is,
// after the kind.
func addForm(doc *openapi.Document, union, model, like, kind string, member *openapi.Schema) error {
	alts, err := unionForms(doc, union)
	if err != nil {
		return err
	}

	m, ok := doc.Components.Schemas[model]
	if !ok {
		return componentErr(model, errors.New("not found"))
	}

	if _, ok := m.Properties[like]; !ok {
		return propertyErr(model, like, errors.New("not found"))
	}

	name := strings.Replace(model, strcase.ToCamel(like), strcase.ToCamel(kind), 1)
	if _, ok := doc.Components.Schemas[name]; ok {
		return componentErr(name, errors.New("already exists"))
	}

	props := openapi.Schemas{}

	for prop, p := range m.Properties.ByIndex() {
		switch prop {
		case "type":
			props.Set(prop, &openapi.Schema{Type: openapi.TypeString, Const: jsontext.Value(strconv.Quote(kind))})
		case like:
			props.Set(kind, member)
		default:
			props.Set(prop, p)
		}
	}

	required := slices.Clone(m.Required)
	if i := slices.Index(required, like); i >= 0 {
		required[i] = kind
	}

	doc.Components.Schemas.Set(name, &openapi.Schema{
		Title: strcase.ToCase(kind, strcase.TitleCase, ' '), Type: openapi.TypeObject,
		Properties: props, Required: required, AdditionalProperties: m.AdditionalProperties,
	})

	*alts = append(*alts, edit.SchemaRef(doc, name))

	return nil
}

// addPropertyKinds adds the kinds of property Notion has but the official spec does not list, made from those of
// last_edited_time.
func addPropertyKinds(doc *openapi.Document) error {
	nullableDateTime := func() *openapi.Schema {
		return &openapi.Schema{
			OneOf: openapi.SchemaList{{Type: openapi.TypeString, Format: openapi.FormatDateTime}, {Type: openapi.TypeNull}},
		}
	}

	for _, f := range []struct {
		union, model, kind string
		member             *openapi.Schema
	}{
		{"PropertyConfig", "lastEditedTimeDatabasePropertyConfigResponse", "place", edit.SchemaRef(doc, "emptyObject")},
		{"PropertyConfig", "lastEditedTimeDatabasePropertyConfigResponse", "last_visited_time", edit.SchemaRef(doc, "emptyObject")},
		{"simplePropertyValueResponse", "lastEditedTimeSimplePropertyValueResponse", "last_visited_time", nullableDateTime()},
		{"propertyItemObjectResponse", "lastEditedTimePropertyItemObjectResponse", "last_visited_time", nullableDateTime()},
	} {
		if err := addForm(doc, f.union, f.model, "last_edited_time", f.kind, f.member); err != nil {
			return err
		}
	}

	return nil
}

// unionForms are the forms of the union name: its oneOf, or that of the part of its allOf that has one.
func unionForms(doc *openapi.Document, name string) (*openapi.SchemaList, error) {
	s, ok := doc.Components.Schemas[name]
	if !ok {
		return nil, componentErr(name, errors.New("not found"))
	}

	switch {
	case len(s.OneOf) > 0:
		return &s.OneOf, nil
	case len(s.AnyOf) > 0:
		return &s.AnyOf, nil
	}

	for _, part := range s.AllOf {
		if len(part.OneOf) > 0 {
			return &part.OneOf, nil
		}
	}

	return nil, componentErr(name, errors.New("is no union"))
}

// extractArrayOf names every array of the component item, which the official spec spells out inline each time.
func extractArrayOf(doc *openapi.Document, name, item string) error {
	if err := edit.ExtractSchema(doc, name, func(s *openapi.Schema) bool {
		return s.Type == openapi.TypeArray && s.Items != nil && s.Items.Ref != nil &&
			s.Items.Ref.Identifier == schemaRefPrefix+item
	}); err != nil {
		return fmt.Errorf("naming the arrays of %s: %w", item, err)
	}

	return nil
}

// nameSorts makes the sorts of a data source's entries one Sort, whether a view's, read or written, or a query's.
func nameSorts(doc *openapi.Document) error {
	if err := edit.RenameSchemas(doc, map[string]string{
		"viewSortResponse":      "Sort",
		"propertySortResponse":  "PropertySort",
		"timestampSortResponse": "TimestampSort",
		"ViewPropertySorts":     "PropertySorts",
	}); err != nil {
		return fmt.Errorf("renaming sorts: %w", err)
	}

	// a view's sorts are written as they are read: its update allows property sorts only
	if err := edit.RedirectSchemas(doc, map[string]string{
		"ViewSort":         "Sort",
		"ViewPropertySort": "PropertySort",
	}); err != nil {
		return fmt.Errorf("redirecting sorts: %w", err)
	}

	doc.Components.Schemas["Sort"].Description = "A sort of a data source's entries, by a property or a timestamp."

	query, err := queryBody(doc)
	if err != nil {
		return err
	}

	sorts, ok := query.Properties["sorts"]
	if !ok || sorts.Items == nil {
		return errors.New("the query of a data source has no sorts")
	}

	sorts.Items = edit.SchemaRef(doc, "Sort")

	if err := extractArrayOf(doc, "Sorts", "Sort"); err != nil {
		return err
	}

	return edit.RedirectSchemas(doc, map[string]string{"ViewSorts": "Sorts"})
}

// nameFilter makes the filters of a data source's entries one Filter, whether a view's, read or written, or a query's.
// It is recursive: a view's filter may nest deeper than the two levels the official spec spells out for a query's.
func nameFilter(doc *openapi.Document) error {
	if _, ok := doc.Components.Schemas["Filter"]; ok {
		return componentErr("Filter", errors.New("already exists"))
	}

	doc.Components.Schemas.Set("Filter", &openapi.Schema{
		Description: "A filter of a data source's entries: by a property or a timestamp, or by any or all of other filters.",
	})
	setFilterForms(doc)

	query, err := queryBody(doc)
	if err != nil {
		return err
	}

	if _, ok := query.Properties["filter"]; !ok {
		return errors.New("the query of a data source has no filter")
	}

	query.Properties.Set("filter", edit.SchemaRef(doc, "Filter"))

	if err := edit.RedirectSchemas(doc, map[string]string{
		"ViewFilter":         "Filter",
		"viewFilterResponse": "Filter",
	}); err != nil {
		return fmt.Errorf("redirecting filters: %w", err)
	}

	edit.RemoveUnreferenced(doc, "groupFilterOperatorArray", "propertyOrTimestampFilterArray", "propertyOrTimestampFilter")

	return nil
}

// setFilterForms sets what Filter may be: any or all of other filters, a property filter, or a timestamp filter.
func setFilterForms(doc *openapi.Document) {
	for _, c := range []struct{ name, op, description string }{
		{"FilterOr", "or", "The entries any of the filters allow."},
		{"FilterAnd", "and", "The entries all of the filters allow."},
	} {
		props := openapi.Schemas{}
		props.Set(c.op, &openapi.Schema{
			Description: c.description, Type: openapi.TypeArray, MaxItems: new(uint(100)), Items: edit.SchemaRef(doc, "Filter"),
		})

		doc.Components.Schemas.Set(c.name, &openapi.Schema{
			Type: openapi.TypeObject, Properties: props, Required: []string{c.op},
			AdditionalProperties: &openapi.AdditionalProperties{},
		})
	}

	doc.Components.Schemas["Filter"].OneOf = openapi.SchemaList{
		edit.SchemaRef(doc, "FilterOr"), edit.SchemaRef(doc, "FilterAnd"), edit.SchemaRef(doc, "propertyFilter"), edit.SchemaRef(doc, "timestampFilter"),
	}
}

// queryBody is the schema of the body of a data source's query.
func queryBody(doc *openapi.Document) (*openapi.Schema, error) {
	const path = "/v1/data_sources/{data_source_id}/query"

	item, ok := doc.Paths[path]
	if !ok || item.Post == nil || item.Post.RequestBody == nil || item.Post.RequestBody.Value == nil {
		return nil, fmt.Errorf("POST %s has no request body", path)
	}

	mt, ok := item.Post.RequestBody.Value.Content["application/json"]
	if !ok || mt.Schema == nil {
		return nil, fmt.Errorf("POST %s has no JSON request body", path)
	}

	return mt.Schema, nil
}

func fixOpenAPI() (*openapi.Document, error) {
	doc, err := openapi.LoadFromFile(pathOpenAPI)
	if err != nil {
		return nil, err
	}

	// fix UUIDs
	if s, ok := doc.Components.Schemas["idRequest"]; ok {
		s.Format = openapi.FormatUUID
	}

	edit.WalkSchemas(doc, func(s *openapi.Schema) {
		if len(s.Properties) == 0 {
			return
		}

		s, ok := s.Properties["data_source_id"]
		if !ok {
			return
		}

		if s.Type == openapi.TypeString {
			s.Format = openapi.FormatUUID
		} else if len(s.OneOf) > 0 && s.OneOf[0].Type == openapi.TypeString {
			s.OneOf[0].Format = openapi.FormatUUID
		}
	})

	for _, p := range doc.Paths.ByIndex() {
		for _, op := range p.Operations {
			for code, override := range map[openapi.StatusCode]string{ //nolint:exhaustive
				"200": "",
				"202": "",
				"400": "",
				"401": "",
				"403": "",
				"404": "",
				"406": "",
				"409": "",
				"429": "",
				"500": "",
				"503": "",
				"504": "",
				"529": "",
			} {
				rsp := op.Responses[code]
				if rsp == nil {
					continue
				}

				rsp.Value.Description = cmp.Or(
					rsp.Value.Description,
					override,
					code.StatusText(),
				)
			}

			for _, p := range op.Parameters {
				if p.Value.Name == "start_cursor" {
					p.Value.Schema.Format = openapi.FormatUUID
				}
			}
		}
	}

	if err := stripSuffixes(doc); err != nil {
		return nil, err
	}

	if err := consolidateErrors(doc, []string{"error_api_", "error_oauth_"},
		"Error", "publicApiCommonErrorResponse", "ErrorCode"); err != nil {
		return nil, err
	}

	notionVersion, ok := doc.Components.Parameters["notionVersion"]
	if !ok {
		return nil, errors.New("parameter notionVersion not found")
	}

	if err := nameEnum(notionVersion.Value.Schema, "Current"); err != nil {
		return nil, &errpath.ErrField{Field: "components", Err: &errpath.ErrField{
			Field: "parameters", Err: &errpath.ErrKey{Key: "notionVersion", Err: &errpath.ErrField{Field: "schema", Err: err}},
		}}
	}

	if err := edit.RenameSchemas(doc, map[string]string{
		"blockObjectResponse":             "Block",
		"databaseObjectResponse":          "Database",
		"databasePropertyConfigResponse":  "PropertyConfig",
		"pageObjectResponse":              "Page",
		"pagePropertyValueWithIdResponse": "PropertyValue",
		"richTextItemResponse":            "RichText",
	}); err != nil {
		return nil, fmt.Errorf("renaming schemas: %w", err)
	}

	if err := addPropertyKinds(doc); err != nil {
		return nil, err
	}

	if err := extractArrayOf(doc, "RichTexts", "RichText"); err != nil {
		return nil, err
	}

	if err := nameSorts(doc); err != nil {
		return nil, err
	}

	if err := nameFilter(doc); err != nil {
		return nil, err
	}

	blocks, err := variants(doc, "Block")
	if err != nil {
		return nil, err
	}

	// each block variant, so that they keep differing only in their type and its member
	if err := addRequestID(doc, append([]string{"Page", "Database", "Error"}, blocks...)...); err != nil {
		return nil, err
	}

	if err := namePartialUnions(doc); err != nil {
		return nil, err
	}

	nameLists(doc)

	if err := setResults(doc); err != nil {
		return nil, err
	}

	nameObjects(doc)
	nameBranches(doc)

	if err := allowDateTimes(doc); err != nil {
		return nil, err
	}

	if err := applyPasses(doc); err != nil {
		return nil, err
	}

	if err := shortenNames(doc); err != nil {
		return nil, err
	}

	if err := doc.Validate(); err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}

	if err := doc.WriteToFile(pathOpenAPI); err != nil {
		return nil, fmt.Errorf("writing to file: %w", err)
	}

	return doc, nil
}

// applyPasses runs, after the edits above, the passes the openapi-* commands would, in the order a spec goes through
// them, then sorts as they do.
func applyPasses(doc *openapi.Document) error {
	ias, err := cassette.InteractionsReadFile("api/interactions.json")
	if err != nil {
		return err
	}

	// the placeholder for the next interaction to record has no response
	ias = slices.DeleteFunc(ias, func(ia *cassette.Interaction) bool { return ia.Response.StatusCode == 0 })

	flattenDoc := func(d *openapi.Document) error { return flatten.Document(d, flatten.Config{MarkOrigin: true}) }

	for _, pass := range []struct {
		name string
		run  func(*openapi.Document) error
	}{
		{"enrich", func(d *openapi.Document) error { return enrich.Enrich(d, ias) }},
		{"flatten", flattenDoc},
		{"compress", func(d *openapi.Document) error { return compress.Document(d, compress.Config{}) }},
		{"flatten again", flattenDoc},
	} {
		if err := pass.run(doc); err != nil {
			return fmt.Errorf("%s: %w", pass.name, err)
		}
	}

	for _, path := range doc.Paths {
		for _, op := range path.Operations {
			op.Responses.Sort()
		}
	}

	doc.Components.SortMaps()

	return nil
}

func generateCode(doc *openapi.Document) error {
	ias, err := cassette.InteractionsReadFile("api/interactions.json")
	if err != nil {
		return err
	}

	ias = slices.DeleteFunc(ias, func(ia *cassette.Interaction) bool { return ia.Response.StatusCode == 0 })

	return codegen.Generate(codegen.Config{
		Interactions: ias,
		Spec:         doc,
		PackageName:  "notion",
		OutputDir:    "pkg/notion",
		Types:        true,
		Client:       true,
		ClientTest:   true,
		Debug:        true,
	})
}
