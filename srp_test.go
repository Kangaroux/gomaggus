package main

import (
	"encoding/csv"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustDecodeHex(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func hexToByteArray(s string, bigEndian bool) *ByteArray {
	return NewByteArray(mustDecodeHex(s), 0, bigEndian)
}

func loadTestData(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	return rows
}

func Test_calcX(t *testing.T) {
	t.Run("generated test data", func(t *testing.T) {
		rows := loadTestData("test_data/calculate_x.csv")

		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := hexToByteArray(row[2], false)
			expected := hexToByteArray(row[3], false).BigInt()

			assert.Equal(t, expected, calcX(username, password, salt))
		}
	})

	t.Run("case insensitive username and password", func(t *testing.T) {
		salt := hexToByteArray("F3DCABA1165E23534CDC7D709E87B409C505C28D26D3DC14247796BE29CC4D24", false)

		assert.Equal(t, calcX("username", "password", salt), calcX("USERNAME", "PASSWORD", salt))
	})
}

func Test_passVerify(t *testing.T) {
	t.Run("generated test data", func(t *testing.T) {
		rows := loadTestData("test_data/calculate_verifier.csv")

		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := hexToByteArray(row[2], false)
			expected := hexToByteArray(row[3], false).BigInt()

			assert.Equal(t, expected, passVerify(username, password, salt))
		}
	})

	t.Run("case insensitive username and password", func(t *testing.T) {
		salt := hexToByteArray("F3DCABA1165E23534CDC7D709E87B409C505C28D26D3DC14247796BE29CC4D24", false)

		assert.Equal(t, passVerify("username", "password", salt), passVerify("USERNAME", "PASSWORD", salt))
	})
}

func Test_calcServerPublicKey(t *testing.T) {
	rows := loadTestData("test_data/calculate_server_public_key.csv")

	for _, row := range rows {
		verifier := hexToByteArray(row[0], false).BigInt()
		serverPrivateKey := hexToByteArray(row[1], false).BigInt()
		expected := hexToByteArray(row[2], false).BigInt()

		assert.Equal(t, expected, calcServerPublicKey(verifier, serverPrivateKey))
	}
}

func Test_calcClientSKey(t *testing.T) {
	rows := loadTestData("test_data/calculate_client_s.csv")

	for _, row := range rows {
		serverPublicKey := hexToByteArray(row[0], false).BigInt()
		clientPrivateKey := hexToByteArray(row[1], false).BigInt()
		x := hexToByteArray(row[2], false).BigInt()
		u := hexToByteArray(row[3], false).BigInt()
		expected := hexToByteArray(row[4], false)

		assert.Equal(t, expected, calcClientSKey(clientPrivateKey, serverPublicKey, x, u))
	}
}

func Test_calcServerSKey(t *testing.T) {
	rows := loadTestData("test_data/calculate_server_s.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		serverPrivateKey := hexToByteArray(row[1], false).BigInt()
		verifier := hexToByteArray(row[2], false).BigInt()
		u := hexToByteArray(row[3], false).BigInt()
		expected := hexToByteArray(row[4], false)

		assert.Equal(t, expected, calcServerSKey(clientPublicKey, verifier, u, serverPrivateKey))
	}
}

func Test_calcU(t *testing.T) {
	rows := loadTestData("test_data/calculate_u.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		serverPublicKey := hexToByteArray(row[1], false).BigInt()
		expected := hexToByteArray(row[2], false).BigInt()

		assert.Equal(t, expected, calcU(clientPublicKey, serverPublicKey))
	}
}

func Test_prepareInterleave(t *testing.T) {
	type testCase struct {
		S        []byte
		expected []byte
	}

	testCases := []testCase{
		{[]byte{0, 1}, []byte{}},
		{[]byte{0, 1, 2, 3}, []byte{2, 3}},
		{[]byte{1, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{0, 0, 0, 0}, []byte{}},
	}

	for _, tc := range testCases {
		S := NewByteArray(tc.S, 0, false)
		assert.Equal(t, tc.expected, prepareInterleave(S))
	}
}

func Test_calcInterleave(t *testing.T) {
	rows := loadTestData("test_data/calculate_interleaved.csv")

	for _, row := range rows {
		S := hexToByteArray(row[0], false)
		expected := hexToByteArray(row[1], false)

		assert.Equal(t, expected, calcInterleave(S))
	}
}

func Test_calcServerSessionKey(t *testing.T) {
	rows := loadTestData("test_data/calculate_server_session_key.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		verifier := hexToByteArray(row[1], false).BigInt()
		serverPrivateKey := hexToByteArray(row[2], false).BigInt()
		expected := hexToByteArray(row[3], false)
		serverPublicKey := calcServerPublicKey(verifier, serverPrivateKey)

		assert.Equal(
			t,
			expected,
			calcServerSessionKey(clientPublicKey,
				serverPublicKey,
				verifier,
				serverPrivateKey,
			),
		)
	}
}

func Test_calcClientSessionKey(t *testing.T) {
	rows := loadTestData("test_data/calculate_client_session_key.csv")

	for _, row := range rows {
		username := row[0]
		password := row[1]
		clientPublicKey := hexToByteArray(row[2], false).BigInt()
		clientPrivateKey := hexToByteArray(row[3], false).BigInt()
		serverPublicKey := hexToByteArray(row[4], false).BigInt()
		salt := hexToByteArray(row[5], false)
		expected := hexToByteArray(row[6], false)

		assert.Equal(
			t,
			expected,
			calcClientSessionKey(username,
				password,
				serverPublicKey,
				clientPrivateKey,
				clientPublicKey,
				salt,
			),
		)
	}
}

func Test_calcServerProof(t *testing.T) {
	// rows := loadTestData("test_data/calculate_client_session_key.csv")

	rows := [][]string{
		strings.Split("BFD1AC65C8DAAAD88BF9DFF9AF8D1DCDF11DFD0C7E398EDCDF5DBBD08EFB39D3 7EBBC190D9AB2DC0CD891372CB30DF1ED35CDA1E 4876E68F9FCCB6CA9BC9C9BCEBDB36F2358B6EAD0F17881D811891A9888E8E5B10E1162CE8B58293 269E3A3EF5DCD15944F043513BDA20D20FEBA2E0", " "),
		strings.Split("3EC64FE225897DD6B7FFE9AA548384268CD217B22E15EBA642DC4E36E84758C8 AAF217B55AAA57C81E3CBA6EC9C5BEBA709C0EBA D59F6E6F2AEB568CF77297A65D913BCB4FAFA90BBB992C46ED79982881B71CBCB1420B0BD9E76F04 1A5F24F41D8A8FF2A04A0C8A8C52B554C31ABAA1", " "),
		strings.Split("DD1D9BDE1B2EED7D773F7823F5AC1CBB13A2DD01BE27B24CDF86A7E1E04AD93C 3AC09CE24E137EDBD6C269CBBC8EEBCE8AF9B06C 037E0DD26B9CE0198A2D55FDACA65D466E453FB531648A7B5CAF9FA48099EAAABAF42486D4DB1AB3 DFD6617610E81725C1C613732202ECF4A1B87161", " "),
		strings.Split("10A16DBD377C6CD1F0CF8F3FF3AA0AA22B0BEAEE1F33DC0B416F2ACC35DEE558 9E7D1A4CB8CC0FD5FF7BF72CF7ECDAFE16EF22C5 4368337C996127EEC14402D7E6CC7EEBCAD9C723D1AFE02E2DAEB20CF7169F7E6B7BE50AD14FFE9D B8376B31F701777E944F8C6BC388157951CB6BB4", " "),
		strings.Split("FB26EF0F80DA26B3C094CF20FCB5BBAF6885125CC332A7DC6EE1DF1423FD9DBD F2ADD1DBDD0E3A92A5BFF1CEC3EED937AEE23C3D 27300BD445AC212EAAAA44610844F9837DF386250259C084F27F6ED376FABDF1CC80B33F1E7C4BFF 7B6ADCF53FB097CFA7A7B9EBC5CB99470DA0EAEB", " "),
		strings.Split("AAED23EA7C055A61EFCEFC8B336F9D4AECCA1BBDABBEF14C7AFF7AFE9B254B59 FEB5044F42A5E7B30BEBFB4FA0AC9FACA26CC0BC 09803B4FBD24EC5266E718C2010D1C0A3AAA95001FB7D10C921C03F6E0AC75F8DBD54B22A8B4FAE6 8B1554D13B79F3AF3BCACB7CF57610D89A2F9F0E", " "),
		strings.Split("FD6C9440CCF2EAB56A2DC5DCB25AC9FAFA38C8DEC2F9E9FEC5BDEC3C46AFDCD4 F7CA0C4D9829C7E29FD587FF21B5FE3EAEAB2B25 0CF646DF2A27FC04249FD8D6812CD4F789FBA2D90CFAC013B80FA05D187767DCEBC0BDFF5EFE75DD AB877B2075CB73B7BE23F2BB8FF03A8896FB4525", " "),
		strings.Split("8F7CD51CFDECFFB8DBC9D4F5DBB52B23AAE6F2AD2BBFBEBB51DD0BB3B2AB96F2 BB7DD9B5ACD7EEFBCEAEDBF07DBEA1814CFBD239 BDC037CC442B814F2A2B66C5A1B2B5588C2587C65CD038D7B486E449F92E08FBE584AE3A337F7B99 98318CCBE9B793012F30F9930DC5E8F6EB51F304", " "),
		strings.Split("F887D27B820BFE61DDB6ABBBC6F765BB988ABEDFCDAD7DFF02CB05F74DD1A366 A39B0BA4C47E9F784E674498D6817BE944C0FF9F DE9B0EFD6E1637F0B8478642205D7C28EF501CD1C20401AE9CD1F5BA28ECC783B11F4AEA984CBF90 E10D52614A26B9E9D92EC057977EE7EC87421C0D", " "),
		strings.Split("51CAB57DE3FC1236FDD0EB5DBDF77D9AA51DA107CC02D41BAA2CC8AA42E7ABBF 9F4CAFCBBCB9DCB13355944DEF44D6C4BFC617A9 64DE52F7CE9BF15AA94E7EF1A2F39206904C3D685AA2765DD04D7237EC3C44991D9E11CF2D7E4B3D DAB3009F03FBC57E4B95395A749136F97C528E47", " "),
	}

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], true).LittleEndian().BigInt()
		clientProof := hexToByteArray(row[1], true)
		sessionKey := hexToByteArray(row[2], true)
		expected := hexToByteArray(row[3], true).LittleEndian()

		assert.Equal(
			t,
			expected,
			calcServerProof(clientPublicKey,
				clientProof,
				sessionKey,
			),
		)
	}
}
