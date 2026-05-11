package gen

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mtgto/pediaroute-go/internal/app/core"
	"github.com/xwb1989/sqlparser"
)

// writeGzipSQL writes SQL content as a gzip-compressed temp file and returns the path.
func writeGzipSQL(t *testing.T, sql string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.sql.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	if _, err := gw.Write([]byte(sql)); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestParsePageStatement(t *testing.T) {
	// namespace=1 (Talk:X) should be filtered; id=4 (D) is a redirect
	stmt, err := sqlparser.Parse(
		"INSERT INTO `page` VALUES (1,0,'A',0),(2,0,'B',0),(99,1,'Talk:X',0),(4,0,'D',1)",
	)
	if err != nil {
		t.Fatal(err)
	}
	pages, err := parsePageStatement(stmt.(*sqlparser.Insert))
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages (namespace 1 filtered), got %d", len(pages))
	}
	if pages[0].id != 1 || pages[0].title != "A" || pages[0].isRedirect {
		t.Errorf("page[0]: got %+v", pages[0])
	}
	if pages[1].id != 2 || pages[1].title != "B" || pages[1].isRedirect {
		t.Errorf("page[1]: got %+v", pages[1])
	}
	if pages[2].id != 4 || pages[2].title != "D" || !pages[2].isRedirect {
		t.Errorf("page[2] (redirect): got %+v", pages[2])
	}
}

func TestParsePageLinkStatement(t *testing.T) {
	// namespace=1 links and fromNamespace=1 links should be filtered
	stmt, err := sqlparser.Parse(
		"INSERT INTO `pagelinks` VALUES (1,0,'B',0),(1,1,'Talk:B',0),(99,0,'C',1),(2,0,'C',0)",
	)
	if err != nil {
		t.Fatal(err)
	}
	links, err := parsePageLinkStatement(stmt.(*sqlparser.Insert))
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("expected 2 links (non-zero namespaces filtered), got %d", len(links))
	}
	if links[0].from != 1 || links[0].title != "B" {
		t.Errorf("link[0]: got %+v", links[0])
	}
	if links[1].from != 2 || links[1].title != "C" {
		t.Errorf("link[1]: got %+v", links[1])
	}
}

func TestLoadPages(t *testing.T) {
	// Lowercase sort: "a" (→"a") < "B" (→"b") < "C" (→"c")
	f := writeGzipSQL(t, "INSERT INTO `page` VALUES (2,0,'B',0),(1,0,'a',0),(3,0,'C',0)")
	pages := loadPages(f)
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	if pages[0].title != "a" || pages[1].title != "B" || pages[2].title != "C" {
		t.Errorf("unexpected sort order: %q %q %q", pages[0].title, pages[1].title, pages[2].title)
	}
}

func TestBuildIndex(t *testing.T) {
	// pages[0]="C"(id=3), pages[1]="A"(id=1), pages[2]="B"(id=2) — order matches array index
	pages := []page{{id: 3, title: "C"}, {id: 1, title: "A"}, {id: 2, title: "B"}}

	idIdx := buildIndex(pages, func(p page) int32 { return p.id })
	if idIdx[0].key != 1 || idIdx[1].key != 2 || idIdx[2].key != 3 {
		t.Errorf("idIndex not sorted: %v", idIdx)
	}
	// id=2 is pages[2] (title="B"), so arrayIndex == 2
	i, ok := lookupIndex(idIdx, int32(2))
	if !ok || i != 2 {
		t.Errorf("lookupIndex id=2: got (%d, %v), want (2, true)", i, ok)
	}
	_, ok = lookupIndex(idIdx, int32(99))
	if ok {
		t.Error("lookupIndex id=99 should return false")
	}

	titleIdx := buildIndex(pages, func(p page) string { return p.title })
	if titleIdx[0].key != "A" || titleIdx[1].key != "B" || titleIdx[2].key != "C" {
		t.Errorf("titleIndex not sorted: %v", titleIdx)
	}
	// "B" is pages[2] (id=3), so arrayIndex == 2
	i, ok = lookupIndex(titleIdx, "B")
	if !ok || i != 2 {
		t.Errorf("lookupIndex title='B': got (%d, %v), want (2, true)", i, ok)
	}
	_, ok = lookupIndex(titleIdx, "Z")
	if ok {
		t.Error("lookupIndex title='Z' should return false")
	}
}

// TestRun is an end-to-end test that runs Run() and verifies all output files.
//
// Graph:  A -> B, C
//         B -> C
//         D -> A   (D is a redirect page)
//         link to "NoSuchPage" is silently ignored
//         namespace-1 page "Talk:X" is excluded
//
// Pages sorted by lowercase title: A(0) B(1) C(2) D(3)
//
// link.dat forward section (indices):  1 2 | 2 | 0
//
//	page 0 (A): fwdIdx=0 len=2  → B(1), C(2)
//	page 1 (B): fwdIdx=2 len=1  → C(2)
//	page 2 (C): fwdIdx=0 len=0  (no outgoing)
//	page 3 (D): fwdIdx=3 len=1  → A(0)
//
// link.dat backward section (indices): 3 | 0 | 0 1
//
//	page 0 (A): bwdIdx=4 len=1  ← D(3)
//	page 1 (B): bwdIdx=5 len=1  ← A(0)
//	page 2 (C): bwdIdx=6 len=2  ← A(0), B(1)
//	page 3 (D): bwdIdx=0 len=0  (no incoming)
func TestRun(t *testing.T) {
	pageSQLFile := writeGzipSQL(t,
		"INSERT INTO `page` VALUES (1,0,'A',0),(2,0,'B',0),(3,0,'C',0),(4,0,'D',1),(99,1,'Talk:X',0)")
	pageLinkSQLFile := writeGzipSQL(t,
		"INSERT INTO `pagelinks` VALUES (1,0,'B',0),(1,0,'C',0),(2,0,'C',0),(4,0,'A',0),(1,0,'NoSuchPage',0)")
	outDir := t.TempDir()

	if err := Run("test", pageSQLFile, pageLinkSQLFile, outDir); err != nil {
		t.Fatal(err)
	}

	// --- config.json ---
	configBytes, err := os.ReadFile(filepath.Join(outDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lang core.Language
	if err := json.Unmarshal(configBytes, &lang); err != nil {
		t.Fatal(err)
	}
	if lang.PageCount != 4 {
		t.Errorf("PageCount: want 4, got %d", lang.PageCount)
	}
	if lang.LinkCount != 4 {
		t.Errorf("LinkCount: want 4, got %d", lang.LinkCount)
	}

	// --- title.dat ---
	// sorted: A B C D, each null-terminated
	wantTitle := []byte("A\x00B\x00C\x00D\x00")
	gotTitle, _ := os.ReadFile(filepath.Join(outDir, "title.dat"))
	if !bytes.Equal(gotTitle, wantTitle) {
		t.Errorf("title.dat: want %q, got %q", wantTitle, gotTitle)
	}

	// --- link.dat ---
	var wantLinkBuf bytes.Buffer
	for _, v := range []uint32{1, 2, 2, 0, 3, 0, 0, 1} {
		binary.Write(&wantLinkBuf, binary.LittleEndian, v)
	}
	gotLink, _ := os.ReadFile(filepath.Join(outDir, "link.dat"))
	if !bytes.Equal(gotLink, wantLinkBuf.Bytes()) {
		t.Errorf("link.dat:\n  want %x\n  got  %x", wantLinkBuf.Bytes(), gotLink)
	}

	// --- page.dat ---
	// 27 bytes per page: int32 + uint32 + uint16 + bool + uint32*4
	var wantPageBuf bytes.Buffer
	type pageRecord struct {
		id         int32
		titleOff   uint32
		titleLen   uint16
		isRedirect bool
		fwdIdx     uint32
		fwdLen     uint32
		bwdIdx     uint32
		bwdLen     uint32
	}
	for _, p := range []pageRecord{
		{1, 0, 1, false, 0, 2, 4, 1},
		{2, 2, 1, false, 2, 1, 5, 1},
		{3, 4, 1, false, 0, 0, 6, 2},
		{4, 6, 1, true, 3, 1, 0, 0},
	} {
		binary.Write(&wantPageBuf, binary.LittleEndian, p.id)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.titleOff)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.titleLen)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.isRedirect)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.fwdIdx)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.fwdLen)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.bwdIdx)
		binary.Write(&wantPageBuf, binary.LittleEndian, p.bwdLen)
	}
	gotPage, _ := os.ReadFile(filepath.Join(outDir, "page.dat"))
	if !bytes.Equal(gotPage, wantPageBuf.Bytes()) {
		t.Errorf("page.dat:\n  want %x\n  got  %x", wantPageBuf.Bytes(), gotPage)
	}
}

// TestRunMultipleStatements verifies that multiple INSERT statements in one SQL file
// are all processed correctly (real Wikipedia dumps split data across many INSERTs).
func TestRunMultipleStatements(t *testing.T) {
	pageSQLFile := writeGzipSQL(t,
		"INSERT INTO `page` VALUES (1,0,'A',0),(2,0,'B',0);\n"+
			"INSERT INTO `page` VALUES (3,0,'C',0);")
	pageLinkSQLFile := writeGzipSQL(t,
		"INSERT INTO `pagelinks` VALUES (1,0,'B',0);\n"+
			"INSERT INTO `pagelinks` VALUES (2,0,'C',0);")
	outDir := t.TempDir()

	if err := Run("test", pageSQLFile, pageLinkSQLFile, outDir); err != nil {
		t.Fatal(err)
	}

	configBytes, _ := os.ReadFile(filepath.Join(outDir, "config.json"))
	var lang core.Language
	if err := json.Unmarshal(configBytes, &lang); err != nil {
		t.Fatal(err)
	}
	if lang.PageCount != 3 {
		t.Errorf("PageCount: want 3, got %d", lang.PageCount)
	}
	if lang.LinkCount != 2 {
		t.Errorf("LinkCount: want 2, got %d", lang.LinkCount)
	}
}
