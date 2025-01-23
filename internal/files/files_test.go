package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileChecksum(t *testing.T) {
	checksum, err := FileChecksum("../../testdata/checksum_testfile")

	assert.Nil(t, err)
	assert.Equal(t,
		"e0ef7229e64c61596d8be928397e19fcc542ac920c4132106fb1ec2295dd73d1",
		checksum)

	checksum, err = FileChecksum("../../testdata/checksum_testfile_not_existing")
	assert.NotNil(t, err)
	assert.Equal(t, "", "")
}

func TestIsFileExisting(t *testing.T) {
	assert.True(t, IsFileExisting("../../testdata/checksum_testfile"))
	assert.False(t, IsFileExisting("../../testdata/checksum_testfile_not_existing"))
}
