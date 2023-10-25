package actions_test

import (
	"rggo/fileSystem/ugz/actions"
	"testing"
)

func TestCheckFlagsAreNotEmpty(t *testing.T) {
	t.Run(
		"Src flag is empty",
		func(t *testing.T) {
			src := ""
			dst := "some_dst"
			actual := actions.CheckFlagsAreNotEmpty(src, dst)
			assertFalse(t, actual)
		})

	t.Run(
		"Dst flag is empty",
		func(t *testing.T) {
			src := "some_src"
			dst := ""
			actual := actions.CheckFlagsAreNotEmpty(src, dst)
			assertFalse(t, actual)
		})

	t.Run(
		"Src and dst flags are not empty",
		func(t *testing.T) {
			src := "some_src"
			dst := "some_dst"
			actual := actions.CheckFlagsAreNotEmpty(src, dst)
			assertTrue(t, actual)
		})
}

func TestExcludePath(t *testing.T) {
	t.Run(
		"Path is directory",
		func(t *testing.T) {
			info := actions.FilterInfo{IsDir: true}
			actual, err := actions.ExcludePath(info)
			assertNotNullError(t, err)
			assertTrue(t, actual)
		})
}

func TestExcludePathByExtension(t *testing.T) {
	testCases := []struct {
		name string
		info actions.FilterInfo
	}{
		{
			"File without extension",
			actions.FilterInfo{IsDir: false, Extension: ""}},
		{
			"File with 'log' extension",
			actions.FilterInfo{IsDir: false, Extension: "log"}},
		{
			"File with 'go' extension",
			actions.FilterInfo{IsDir: false, Extension: "go"}},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				actual, err := actions.ExcludePath(tc.info)
				assertNotNullError(t, err)
				assertTrue(t, actual)
			})
	}
}

func TestIncludePath(t *testing.T) {
	testCases := []struct {
		name string
		info actions.FilterInfo
	}{
		{
			"Archive extension 'gz'",
			actions.FilterInfo{IsDir: false, Extension: "gz"}},
		{
			"Archive extension 'GZ'",
			actions.FilterInfo{IsDir: false, Extension: "GZ"}},
		{
			"Archive extension 'Gz'",
			actions.FilterInfo{IsDir: false, Extension: "Gz"}},
		{
			"Archive extension 'gZ'",
			actions.FilterInfo{IsDir: false, Extension: "gZ"}},
	}
	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				actual, err := actions.ExcludePath(tc.info)
				assertNotNullError(t, err)
				assertFalse(t, actual)
			})
	}
}

func TestGetFileWithoutGZExtension(t *testing.T) {
	t.Run(
		"Input filename dir.log.gz",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("dir.log.gz")
			assertEqualStrings(t, "dir.log", actual)
		})

	t.Run(
		"Input filename test.log.GZ",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("test.Log.GZ")
			assertEqualStrings(t, "test.Log", actual)
		})

	t.Run(
		"Input filename test.go.gZ",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("test.go.gZ")
			assertEqualStrings(t, "test.go", actual)
		})

	t.Run(
		"Input filename test.go.Gz",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("test.go.Gz")
			assertEqualStrings(t, "test.go", actual)
		})

	t.Run(
		"Input filename Test",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("Test")
			assertEqualStrings(t, "Test", actual)
		})

	t.Run(
		"Input filename test.go",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("test.go")
			assertEqualStrings(t, "test.go", actual)
		})

	t.Run(
		"Input empty string",
		func(t *testing.T) {
			actual := actions.GetFileWithoutGZ("")
			assertEqualStrings(t, "", actual)
		})
}

func assertTrue(t *testing.T, actual bool) {
	if actual {
		return
	}

	t.Error("Expected true")
}

func assertFalse(t *testing.T, actual bool) {
	if !actual {
		return
	}

	t.Error("Expected false")
}

func assertNotNullError(t *testing.T, err error) {
	if err == nil {
		return
	}

	t.Error("Value is not null")
}

func assertEqualStrings(t *testing.T, expected, actual string) {
	if expected == actual {
		return
	}

	t.Errorf("\nExpected: '%s'\n Actual: '%s'\n", expected, actual)
}
