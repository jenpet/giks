package git

import (
	"giks/test/gittest"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const testGitDir = "../test/output/git-dir"

func TestGetMixin_whenMixinIsUnknown(t *testing.T) {
	m, err := GetMixin("UNKNOWN")
	assert.Nil(t, m, "unknown mixins should result in a nil response")
	assert.Error(t, err, "error expected when looking up an unknown mixin")
}

func TestStagedFilesLister(t *testing.T) {
	r := gittest.NewTestRepository(testGitDir)
	defer r.Clean()
	r.WriteFile("README", "please read me")
	r.AddAll()
	vars := map[string]string{}
	m, _ := GetMixin("STAGED_FILES")
	b, err := m.Enrich(r.AbsDir(), vars)
	assert.NoError(t, err, "no error expected when enriching vars")
	assert.True(t, b, "enriching vars should return no error")
	assert.Contains(t, strings.Split(vars["GIKS_MIXIN_STAGED_FILES"], " "), "README", "expected affected files to have added file")
}
