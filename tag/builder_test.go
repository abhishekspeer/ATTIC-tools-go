package tag

import (
	"github.com/vladvelici/spdx-go/spdx"
	"testing"
)

func sameStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sameLicence(a, b spdx.AnyLicenceInfo) bool {
	switch ta := a.(type) {
	default:
		return false
	case spdx.LicenceReference:
		if tb, ok := b.(spdx.LicenceReference); ok && tb == ta {
			return true
		}
		return false
	case spdx.DisjunctiveLicenceList:
		if tb, ok := b.(spdx.DisjunctiveLicenceList); ok && len(ta) == len(tb) {
			for i, lica := range ta {
				licb := tb[i]
				if !sameLicence(lica, licb) {
					return false
				}
			}
			return true
		}
		return false
	case spdx.ConjunctiveLicenceList:
		if tb, ok := b.(spdx.ConjunctiveLicenceList); ok && len(ta) == len(tb) {
			for i, lica := range ta {
				licb := tb[i]
				if !sameLicence(lica, licb) {
					return false
				}
			}
			return true
		}
		return false
	}
}

func joinStrSlice(a []string, glue string) string {
	if len(a) == 0 {
		return ""
	}
	result := a[0]
	for i := 1; i < len(a); i++ {
		result += glue + a[i]
	}
	return result
}

func TestUpd(t *testing.T) {
	a := ""
	f := upd(&a)
	f("world")
	if a != "world" {
		t.Fail()
	}
}

func TestUpdFail(t *testing.T) {
	a := ""
	f := upd(&a)
	f("hello")
	err := f("world")
	if err != ErrAlreadyDefined {
		t.Errorf("Different error: %s", err)
	}
}

func TestUpdList(t *testing.T) {
	arr := []string{"1", "2", "3"}
	f := updList(&arr)
	f("4")
	if len(arr) != 4 || arr[3] != "4" {
		t.Fail()
	}
}

func TestVerifCodeNoExcludes(t *testing.T) {
	vc := new(spdx.VerificationCode)
	value := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	f := verifCode(vc)
	err := f(value)
	if err != nil {
		t.Errorf("Error should be nil but found %s.", err)
	}
	if vc.Value != value {
		t.Error("Verification code value different than given value")
	}

	if len(vc.ExcludedFiles) != 0 {
		t.Errorf("ExcludedFiles should be null/empty but found %s.", vc.ExcludedFiles)
	}
}

func TestVerifCodeWithExcludes(t *testing.T) {
	vc := new(spdx.VerificationCode)
	value := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	excludes := " (excludes: abc.txt, file.spdx)"
	f := verifCode(vc)
	err := f(value + excludes)

	if err != nil {
		t.Errorf("Error should be nil but found %s.", err)
	}
	if vc.Value != value {
		t.Error("Verification code value different than given value")
	}

	if !sameStrSlice(vc.ExcludedFiles, []string{"abc.txt", "file.spdx"}) {
		t.Error("Different ExcludedFiles. Elements found: %s.", vc.ExcludedFiles)
	}
}

func TestVerifCodeWithExcludes2(t *testing.T) {
	vc := new(spdx.VerificationCode)
	value := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	excludes := " (abc.txt, file.spdx)"
	f := verifCode(vc)

	err := f(value + excludes)

	if err != nil {
		t.Errorf("Error should be nil but found %s.", err)
	}

	if vc.Value != value {
		t.Error("Verification code value different than given value")
	}

	if !sameStrSlice(vc.ExcludedFiles, []string{"abc.txt", "file.spdx"}) {
		t.Error("Different ExcludedFiles. Elements found: %s.", vc.ExcludedFiles)
	}
}

func TestVerifCodeInvalidExcludesNoClosedParentheses(t *testing.T) {
	vc := new(spdx.VerificationCode)
	value := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	excludes := " ("
	f := verifCode(vc)
	err := f(value + excludes)

	if err != ErrNoClosedParen {
		t.Errorf("Error should be UnclosedParentheses but found %s.", err)
	}
}

func TestChecksum(t *testing.T) {
	cksum := new(spdx.Checksum)
	val := "d6a770ba38583ed4bb4525bd96e50461655d2758"
	algo := "SHA1"

	f := checksum(cksum)
	err := f(algo + ": " + val)

	if err != nil {
		t.Errorf("Error should be nil but found %s.", err)
	}

	if cksum.Value != "d6a770ba38583ed4bb4525bd96e50461655d2758" {
		t.Errorf("Checksum value is wrong. Found: '%s'. Expected: '%s'.", cksum.Value, val)
	}

	if cksum.Algo != algo {
		t.Errorf("Algo is wrong. Found: '%s'. Expected: '%s'.", cksum.Algo, algo)
	}
}

func TestChecksumInvalid(t *testing.T) {
	cksum := new(spdx.Checksum)
	f := checksum(cksum)
	err := f("d6a770ba38583ed4bb")

	if err != ErrInvalidChecksum {
		t.Errorf("Invalid error found: %s.", err)
	}
}

// findMatchingParenSet tests:

func TestFindMachingParen(t *testing.T) {
	input := "a(f())d"
	//        0123456
	o, c := findMatchingParenSet(input)
	if o != 1 || c != 5 {
		t.Errorf("Expected o=1 and c=5 but have o=%d and c=%d", o, c)
	}
}

func TestFindMachingParenUnbalanced(t *testing.T) {
	input := "a(f()d"
	//        012345
	o, c := findMatchingParenSet(input)
	if o != 1 || c >= 0 {
		t.Errorf("Expected o=1 and c<0 but have o=%d and c=%d", o, c)
	}
}

func TestFindMachingParenNextElem(t *testing.T) {
	input := "a()"
	//        012
	o, c := findMatchingParenSet(input)
	if o != 1 || c != 2 {
		t.Errorf("Expected o=1 and c=2 but have o=%d and c=%d", o, c)
	}
}

func TestFindMachingParenNoParenthesis(t *testing.T) {
	input := "blahblah"
	o, c := findMatchingParenSet(input)
	if o >= 0 || c >= 0 || c >= o {
		t.Errorf("Expected c<o<0 but have o=%d and c=%d", o, c)
	}
}

func TestConjDisjAND(t *testing.T) {
	input := "a and b"
	conj, disj := conjOrDisjSet(input)
	if conj != true || disj != false {
		t.Errorf("Expected cojunctive. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestConjDisjOR(t *testing.T) {
	input := "a or b"
	conj, disj := conjOrDisjSet(input)
	if conj != false || disj != true {
		t.Errorf("Expected disjunctive. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestConjDisjBOTH(t *testing.T) {
	input := "a and b or c"
	conj, disj := conjOrDisjSet(input)
	if conj != true || disj != true {
		t.Errorf("Expected conjunctive and disjunctive. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestConjDisjParenOR(t *testing.T) {
	input := "(a and b) or c"
	conj, disj := conjOrDisjSet(input)
	if conj != false || disj != true {
		t.Errorf("Expected disjunctive. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestConjDisjMultipleParen(t *testing.T) {
	input := "(a and b and (c or d)) or (a or d)"
	conj, disj := conjOrDisjSet(input)
	if conj != false || disj != true {
		t.Errorf("Expected disjunctive. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestConjDisjNone(t *testing.T) {
	input := "GPLv3"
	conj, disj := conjOrDisjSet(input)
	if conj != false || disj != false {
		t.Errorf("Expected neither. Found: conj=%b, disj=%b", conj, disj)
	}
}

func TestLicenceSetSplit(t *testing.T) {
	input := "a and (b or c)"
	expected := []string{"a", "(b or c)"}
	output := licenceSetSplit(andSeprator, input)

	if !sameStrSlice(expected, output) {
		t.Errorf("Expected %s but found %s", expected, output)
	}
}

func TestLicenceSetSplitSameTypeInParen(t *testing.T) {
	input := "a and (b and c) and d"
	expected := []string{"a", "(b and c)", "d"}
	output := licenceSetSplit(andSeprator, input)

	if !sameStrSlice(expected, output) {
		t.Errorf("Expected %s but found %s", expected, output)
	}
}

func TestLicenceSetSplitNoParen(t *testing.T) {
	input := "a and b  and  c "
	expected := []string{"a", "b", "c"}
	output := licenceSetSplit(andSeprator, input)

	if !sameStrSlice(expected, output) {
		t.Errorf("Expected %s but found %s", expected, output)
	}
}

func TestLicenceSetSplitNoSeparator(t *testing.T) {
	input := "a"
	expected := []string{"a"}
	output := licenceSetSplit(andSeprator, input)

	if !sameStrSlice(expected, output) {
		t.Errorf("Expected %s but found %s", expected, output)
	}
}

func TestParseLicenceSetOr(t *testing.T) {
	input := "(GPLv3 or LicenseRef-1)"
	expected := spdx.DisjunctiveLicenceList{spdx.NewLicenceReference("GPLv3"), spdx.NewLicenceReference("LicenseRef-1")}
	output, err := parseLicenceSet(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !sameLicence(output, expected) {
		t.Errorf("\nExpected: %T : %s\nbut found %T : %s\n", expected, expected, output, output)
	}
}

func TestParseLicenceSetAnd(t *testing.T) {
	input := "(GPLv3 and LicenseRef-1)"
	expected := spdx.ConjunctiveLicenceList{spdx.NewLicenceReference("GPLv3"), spdx.NewLicenceReference("LicenseRef-1")}
	output, err := parseLicenceSet(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !sameLicence(output, expected) {
		t.Errorf("\nExpected: %T : %s\nbut found %T : %s\n", expected, expected, output, output)
	}
}

func TestParseLicenceSetNested(t *testing.T) {
	input := "(GPLv3 or (LicenseRef-1 and LicenseRef-3) or LicenseRef-2)"

	expected := spdx.DisjunctiveLicenceList{
		spdx.NewLicenceReference("GPLv3"),
		spdx.ConjunctiveLicenceList{
			spdx.NewLicenceReference("LicenseRef-1"),
			spdx.NewLicenceReference("LicenseRef-3"),
		},
		spdx.NewLicenceReference("LicenseRef-2"),
	}

	output, err := parseLicenceSet(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !sameLicence(output, expected) {
		t.Errorf("\nExpected: %T : %s\nbut found %T : %s\n", expected, expected, output, output)
	}
}

func TestParseLicenceSetSingleValue(t *testing.T) {
	input := "one value "

	expected := spdx.NewLicenceReference("one value")

	output, err := parseLicenceSet(input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !sameLicence(output, expected) {
		t.Errorf("\nExpected: %T : %s\nbut found %T : %s\n", expected, expected, output, output)
	}
}

func TestParseLicenceSetEmptyLicence(t *testing.T) {
	input := " "

	_, err := parseLicenceSet(input)
	if err != ErrEmptyLicence {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestParseLicenceStringEmptyLicence(t *testing.T) {
	input := " "

	_, err := parseLicenceString(input)
	if err != ErrEmptyLicence {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestParseLicenceStringUnbalancedParentheses(t *testing.T) {
	input := " (()"

	_, err := parseLicenceString(input)
	if err != ErrNoClosedParen {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestLicenceSetConjunctionAndDisjunction(t *testing.T) {
	input := "a and b or c"
	_, err := parseLicenceSet(input)
	if err != ErrConjunctionAndDisjunction {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestDoc(t *testing.T) {
	m := map[string][]string{
		"SpecVersion":        {"1.2"},
		"DataLicense":        {"CC-1.0"},
		"DocumentComment":    {"hahaha hahaha"},
		"Creator":            {"Organization: F", "Person: Test Testaculous"},
		"Created":            {"04/05/06"},
		"CreatorComment":     {"This is some creator comment...\n blah blah"},
		"LicenseListVersion": {"1.2.3"},
	}

	// transform data to input pair
	input := make([]pair, 0)
	for k, vals := range m {
		for _, v := range vals {
			input = append(input, pair{k, v})
		}
	}

	doc, err := build(input)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if doc.SpecVersion != m["SpecVersion"][0] {
		t.Errorf("h Invalid doc.SpecVersion: '%s'", doc.SpecVersion)
	}
	if doc.DataLicence != m["DataLicense"][0] {
		t.Errorf("Invalid doc.DataLicence: '%s'", doc.DataLicence)
	}
	if doc.Comment != m["DocumentComment"][0] {
		t.Errorf("Invalid doc.Comment: '%s'", doc.Comment)
	}
	if !sameStrSlice(doc.CreationInfo.Creator, m["Creator"]) {
		t.Errorf("Invalid doc.CreationInfo.Creator: (len=%d) '%s'", len(doc.CreationInfo.Creator), doc.CreationInfo.Creator)
	}
	if doc.CreationInfo.Created != m["Created"][0] {
		t.Errorf("Invalid doc.CreationInfo.Created: '%s'", doc.CreationInfo.Created)
	}
	if doc.CreationInfo.Comment != m["CreatorComment"][0] {
		t.Errorf("Invalid doc.CreationInfo.Comment: '%s'", doc.CreationInfo.Comment)
	}
	if doc.CreationInfo.LicenceListVersion != m["LicenseListVersion"][0] {
		t.Errorf("Invalid doc.LicenceListVersion: '%s'", doc.CreationInfo.LicenceListVersion)
	}
}

func TestSamePropertyTwice(t *testing.T) {
	input := []pair{
		{"SpecVersion", "1.2"},
		{"SpecVersion", "1.3"},
	}

	_, err := build(input)

	if err != ErrAlreadyDefined {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestDocNestedPackage(t *testing.T) {
	cksum := spdx.Checksum{
		Algo:  "SHA1",
		Value: "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
	}
	input := []pair{
		{"SpecVersion", "1.2"},
		{"PackageName", "spdx-tools-go"},
		{"PackageVersion", "1.2"},
		{"PackageFileName", "spdx-tools-go.tar.gz"},
		{"PackageSupplier", "Organization: Linux Foundation"},
		{"PackageOriginator", "Organization: Linux Foundation"},
		{"PackageDownloadLocation", "http://git.spdx.org/spdx-tools-go.git"},
		{"PackageHomePage", "http://git.spdx.org/spdx-tools-go.git"},
		{"PackageSourceInfo", ":)"},
		{"PackageLicenseComments", "example comment"},
		{"PackageCopyrightText", "sample copyright text"},
		{"PackageSummary", "spdx library for go"},
		{"PackageDescription", "spdx parser library for go supporting tag and rdf formats"},
		{"PackageVerificationCode", "fdsa431fdsa43fdsa432"},
		{"PackageLicenseDeclared", "MIT1"},
		{"PackageLicenseConcluded", "MIT2"},
		// non-string comparable attributes:
		{"PackageChecksum", cksum.Algo + ": " + cksum.Value},
		{"PackageLicenseInfoFromFiles", "Apache2"},
		{"PackageLicenseInfoFromFiles", "MIT"},
	}

	doc, err := build(input)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(doc.Packages) != 1 {
		t.Error("None or more packages than expected.")
	}
	pkg := doc.Packages[0]

	strExpected := []string{
		doc.SpecVersion,
		pkg.Name,
		pkg.Version,
		pkg.FileName,
		pkg.Supplier,
		pkg.Originator,
		pkg.DownloadLocation,
		pkg.HomePage,
		pkg.SourceInfo,
		pkg.LicenceComments,
		pkg.CopyrightText,
		pkg.Summary,
		pkg.Description,
		pkg.VerificationCode.Value,
		pkg.LicenceDeclared.LicenceId(),
		pkg.LicenceConcluded.LicenceId(),
	}

	for i, val := range strExpected {
		p := input[i]
		if val != p.Value {
			t.Errorf("Invalid %s: '%s'.", p.Key, val)
		}
	}

	if *pkg.Checksum != cksum {
		t.Errorf("Invalid PackageChecksum: '%s'.", *pkg.Checksum)
	}

	if len(pkg.LicenceInfoFromFiles) != 2 ||
		pkg.LicenceInfoFromFiles[0].LicenceId() != "Apache2" ||
		pkg.LicenceInfoFromFiles[1].LicenceId() != "MIT" {

		t.Errorf("Invalid PackageLicenseInfoFromFiles: '%s'", pkg.LicenceInfoFromFiles)
	}
}

func TestDocNestedFiles(t *testing.T) {
	cksum := spdx.Checksum{
		Algo:  "SHA1",
		Value: "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12",
	}

	fContributor := []string{"Person: A", "Person: B", "Organization: LF"}
	fDependency := []string{"f0.txt", "f1.txt", "f2.txt"}
	licInfo := []string{"MIT", "Apache", "LicenseRef-0"}
	input := []pair{
		{"SpecVersion", "1.2"},
		{"PackageName", "spdx-tools-go"},

		{"FileName", "a.txt"},
		{"FileType", "BINARY"},
		{"LicenseConcluded", "MIT"},
		{"LicenseComments", "file licence comments"},
		{"FileCopyrightText", "some copyright text"},
		{"FileComment", "some file comment"},
		{"FileNotice", "example file notice"},

		// nested attributes
		{"ArtifactOfProjectName", "spdx-tools-go"},
		{"ArtifactOfProjectUri", "http://git.spdx.org/spdx-tools-go.git"},
		{"ArtifactOfProjectHomepage", "http://git.spdx.org/spdx-tools-go.git"},

		{"ArtifactOfProjectName", "spdx"},
		{"ArtifactOfProjectName", "spdx-tools"},

		// second file
		{"FileName", "b.go"},
		{"FileType", "SOURCE"},

		// non-string checks
		{"FileChecksum", cksum.Algo + ": " + cksum.Value},
	}

	for _, v := range fContributor {
		input = append(input, pair{"FileContributor", v})
	}

	for _, v := range fDependency {
		input = append(input, pair{"FileDependency", v})
	}

	for _, v := range licInfo {
		input = append(input, pair{"LicenseInfoInFile", v})
	}

	doc, err := build(input)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(doc.Packages) != 1 {
		t.Error("None or more packages than expected.")
	}
	pkg := doc.Packages[0]

	if len(doc.Files) != 2 {
		t.Errorf("Package has %s file(s). Expected %s", len(doc.Files), 2)
	}
	file1 := doc.Files[0]
	file2 := doc.Files[1]

	if len(file1.ArtifactOf) != 3 {
		t.Errorf("File 1 has %s ArtifactOf. Expected %s", len(file1.ArtifactOf), 3)
	}
	if len(file2.ArtifactOf) != 0 {
		t.Errorf("File 2 has %s ArtifactOf. Expected %s", len(file2.ArtifactOf), 0)
	}

	// check all string properties
	strExpected := []string{
		doc.SpecVersion,
		pkg.Name,

		file1.Name,
		file1.Type,
		file1.LicenceConcluded.LicenceId(),
		file1.LicenceComments,
		file1.CopyrightText,
		file1.Comment,
		file1.Notice,
		file1.ArtifactOf[0].Name,
		file1.ArtifactOf[0].ProjectUri,
		file1.ArtifactOf[0].HomePage,
		file1.ArtifactOf[1].Name,
		file1.ArtifactOf[2].Name,

		file2.Name,
		file2.Type,
	}

	for i, val := range strExpected {
		p := input[i]
		if val != p.Value {
			t.Errorf("Invalid %s: '%s'.", p.Key, val)
		}
	}

	// same LicenceInfoInFile
	if len(file2.LicenceInfoInFile) != len(licInfo) {
		t.Errorf("Wrong file2.LicenceInfoInFile: '%s'", file2.LicenceInfoInFile)
	} else {
		fail := false
		for i, lic := range file2.LicenceInfoInFile {
			if lic.LicenceId() != licInfo[i] {
				fail = true
				break
			}
		}
		if fail {
			t.Errorf("Wrong file2.LicenceInfoInFile: '%s'", file2.LicenceInfoInFile)
		}
	}

	// sameContributors
	if !sameStrSlice(fContributor, file2.Contributor) {
		t.Errorf("Wrong file2.Contributor: '%s'. Expected: '%s'", file2.Contributor, fContributor)
	}

	if *file2.Checksum != cksum {
		t.Errorf("Invalid PackageChecksum: '%s'.", *pkg.Checksum)
	}

}

func TestDocLicenceId(t *testing.T) {
	ids := [...]string{"LicenseRef-1", "LicenseRef-2"}
	xRef := [...][]string{{"A", "B", "Ca"}, {}}
	names := [...][]string{{"MIT", "Apache", "LicenseRef-0"}, {"Lic", "Example"}}

	input := make([]pair, 0)

	for i, id := range ids {
		input = append(input, pair{"LicenseID", id})
		for _, ref := range xRef[i] {
			input = append(input, pair{"LicenseCrossReference", ref})
		}
		for _, name := range names[i] {
			input = append(input, pair{"LicenseName", name})
		}
	}

	doc, err := build(input)

	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
		t.FailNow()
	}

	if len(doc.ExtractedLicenceInfo) != len(ids) {
		t.Errorf("Wrong ExtractedLicenceInfo (len=%s): '%s'", len(doc.ExtractedLicenceInfo), doc.ExtractedLicenceInfo)
		t.FailNow()
	}

	for i, extr := range doc.ExtractedLicenceInfo {
		if extr.Id != ids[i] {
			t.Errorf("(#%d) Wrong ID: '%s' (expected '%s')", i, extr.Id, ids[i])
		}
		if !sameStrSlice(extr.Name, names[i]) {
			t.Errorf("(#%d) Wrong Name: '%s' (expected '%s')", i, extr.Name, names[i])
		}
		if !sameStrSlice(extr.CrossReference, xRef[i]) {
			t.Errorf("(#%d) Wrong CrossReference: '%s' (expected '%s')", i, extr.CrossReference, xRef[i])
		}
	}
}

func TestDocInvalidProperty(t *testing.T) {
	input := []pair{{"SomeInvalidProperty", "value"}}
	_, err := build(input)
	if err == nil {
		t.Fail()
	}
}

func TestLicenceAlreadyDefined(t *testing.T) {
	input := []pair{
		{"PackageName", "some package"},
		{"PackageLicenseDeclared", "decl1"},
		{"PackageLicenseDeclared", "decl2"},
	}
	_, err := build(input)
	if err != ErrAlreadyDefined {
		t.Errorf("Unexpected error '%s'.", err)
	}
}

func TestDocLicenceConjAndDisj(t *testing.T) {
	input := []pair{
		{"PackageName", "some package"},
		{"PackageLicenseDeclared", "(decl1 and decl2 or decl3)"},
	}
	_, err := build(input)
	if err != ErrConjunctionAndDisjunction {
		t.Errorf("Unexpected error '%s'.", err)
	}
}

func TestDocLicenceInvalidSet(t *testing.T) {
	input := []pair{
		{"PackageName", "some package"},
		{"PackageLicenseDeclared", "a and ()"},
	}
	_, err := build(input)
	if err != ErrEmptyLicence {
		t.Errorf("Unexpected error '%s'.", err)
	}
}

func TestChecksumAlreadyDefined(t *testing.T) {
	input := []pair{
		{"PackageName", "some package"},
		{"PackageChecksum", "SHA1: d6a770ba38583ed4bb4525bd96e50471655d2758"},
		{"PackageChecksum", "SHA1: d6a770ba38583ed4bb4525bd96e50471655d2758"},
	}
	_, err := build(input)
	if err != ErrAlreadyDefined {
		t.Errorf("Unexpected error '%s'.", err)
	}
}

func TestVerifCodeAlreadyDefined(t *testing.T) {
	input := []pair{
		{"PackageName", "some package"},
		{"PackageVerificationCode", "6a770ba38583ed4bb4525bd96e50471655d2758"},
		{"PackageVerificationCode", "d6a770ba38583ed4bb4525bd96e50471655d2758"},
	}
	_, err := build(input)
	if err != ErrAlreadyDefined {
		t.Errorf("Unexpected error '%s'.", err)
	}
}

func TestReviewer(t *testing.T) {

	reviews := []spdx.Review{
		{"a", "b", "c"},
		{"d", "e", "f"},
	}

	input := make([]pair, 0, 6)
	for _, rev := range reviews {
		input = append(input, pair{"Reviewer", rev.Reviewer})
		input = append(input, pair{"ReviewDate", rev.Date})
		input = append(input, pair{"ReviewComment", rev.Comment})
	}

	doc, err := build(input)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		t.FailNow()
	}

	if len(doc.Reviews) != len(reviews) {
		t.Errorf("Invalid reviews (len=%d): '%s' (expected (len=%d): '%s')", len(doc.Reviews), doc.Reviews, len(reviews), reviews)
		t.FailNow()
	}

	for i, rev := range doc.Reviews {
		if *rev != reviews[i] {
			t.Errorf("Invalid reviews (len=%d): '%s' (expected (len=%d): '%s')", len(doc.Reviews), doc.Reviews, len(reviews), reviews)
			t.FailNow()
		}
	}
}