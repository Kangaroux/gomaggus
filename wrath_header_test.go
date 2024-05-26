package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateKey(t *testing.T) {
	rows := loadTestData("test_data/header/wrath_generate_key.csv")

	for _, row := range rows {
		sessionKey := hexToByteArray(row[0], false)
		fixedKey := hexToByteArray(row[1], false)
		expected := hexToByteArray(row[2], false)
		wrath := &WrathHeaderCrypto{
			sessionKey: sessionKey,
		}

		assert.Equal(t, expected.Bytes(), wrath.generateKey(fixedKey))
	}
}
