package iqshell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTsWithNewDomain(t *testing.T) {
	line := "http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts"
	newDomain := ""
	removeSparePreSlash := true
	expected := "/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts"

	checkReplaceTsWithNewDomain(t, line, newDomain, removeSparePreSlash, expected)
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		false,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		false,
		"//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		true,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")

	checkReplaceTsWithNewDomain(t,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		true,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		false,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		true,
		"/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		false,
		"//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		false,
		"hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"",
		true,
		"hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")

	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com",
		false,
		"http://puts.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com",
		true,
		"http://puts.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com",
		true,
		"http://puts.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com",
		false,
		"http://puts.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")

	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com/",
		false,
		"http://puts.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com/",
		true,
		"http://puts.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"https://puts.clouddn.com/",
		true,
		"https://puts.clouddn.com/hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
	checkReplaceTsWithNewDomain(t,
		"http://putsdf.bkt.clouddn.com//hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts",
		"http://puts.clouddn.com/",
		false,
		"http://puts.clouddn.com///hfuewdjhfjdsekgske_Tde7/Hke29839_df/0001.ts")
}

func checkReplaceTsWithNewDomain(t *testing.T, line, newDomain string,
	removeSparePreSlash bool, expected string) {
	newLine := replaceTsNewDomain(line, newDomain, removeSparePreSlash)
	assert.Equal(t, expected, newLine)
}
