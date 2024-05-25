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
	rows := loadTestData("test_data/calculate_server_proof.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		clientProof := hexToByteArray(row[1], false)
		sessionKey := hexToByteArray(row[2], false)
		expected := hexToByteArray(row[3], false).LittleEndian()

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

func Test_calcClientProof(t *testing.T) {
	// rows := loadTestData("test_data/calculate_server_proof.csv")
	rows := [][]string{
		strings.Split("7WG6SHZL33JMGPO4 2F409C9AEC0FE203D3673202D57BEA19C931AACBD1FD75C539C34129BD70F83E37BFC0F99CD3A477 0095FE039AFE5E1BADE9AC0CAEC3CB73D2D08BBF4CA8ADDBCDF0CE709ED5103F 00B0C41F58CCE894CFB816FA72CA344C9FE2ED7CE799452ADBA7ABDCD26EAE75 00A4A09E0B5ACA438B8CD837D0816CA26043DBD1EAEF138EEF72DCF3F696D03D 7D07022B4064CCE633D679F61C6B212B6F8BC5C3", " "),
		strings.Split("5710JG4OTXVLEMPT 0FA3BB94B7F4DFF1BF4FC7D170B211746063F82024497D8434986E99B716D4D778CB6A8697AD347C BD46C6E1FD47B110DC323AED761BED3D2E9CCF2FCB70DF1EDFDFC395A8BBAE9F D9EBBC0A87873A0A56E7E4FBA3BE9D4B7A59CD0EC3F5BCB4B86BA32DA0BA26B7 BFDC80CDA8BA261CBA4EFC77FCC0C7BFE6A3CC0E88BEBBD7BFCFEEE0DF3CAEBC BF82301C25720FAFA8360A49DAFBC66A5F58AE68", " "),
		strings.Split("XL3PN03LQF15WNLF C12137EBCF3D83E159CE3B29EFADB97B6BB86199F7120CE5F7A6BB7A901F1812EFE3792E3E2A872D BDF7EAFD517B2F7AC43CBB5EEA100FFBFBE0E4C5F15D4E4ABCB340D8FED1D6AF F06A1AE57ACCD950447568246E3FDC0D927ABDAFC9DD5F9E62CBDDDF681548FF 1D5DD08DAEE73B3FF06FE15EE4B4EBAF738D7DA3EAD994B927198A513532F4AD EB7B43B30325F80AAF46926B3140FAEDE156569B", " "),
		strings.Split("YPVUAHPZSTKMO4LL 69DF397C8C3E89792DF5DEFDCDA809C8919B5AB5352C155555112292807DA8B0FE54EE29305BEEC5 1ABC9AFE6FA708BE5E4F63DC7A9DCEC2BD9DBABC3EA7BFFF76EEB8DA96BCE671 1EC1745EA3A1A6BBDF2ECA884C38CCCFFF1EF65ADDAE3BECB9F8CADBAE90FECA E0DAADEC8FCB57157DAB2EA3DB499CAB03BCCF62AADE503789C4EB227B7D8CBA 6E8DBB106A4111E3F2EAC98A734A31F5AEA624FE", " "),
		strings.Split("VKFHEA0AZA6BHZZO 79E90C76B44389D287B30833CE61DB7BCB29E4239C16A37FAC8FF53AB10C7BFEDC3F504893F2DF84 5AD2BCC6FE7C6BACC3BD898A9C7EFC3DCEFD892BFEA4B93C7D77C7F09C6A22C9 C88CE7AE4EDFDFCEADA9EE9E67BD0E5224C1CFCA3E45AF9D224AD5813F44CF49 B0DBFF06EF84DFD6F34D4DE762D2ECDDDE31963E5E5BA9CA0DE09D3D1E86B3D7 3264022E55627604A4AE9496D174B47269242C19", " "),
		strings.Split("6HDMTYJHXA19ZKIX 359C3B88FCAC11D3F0BE0902BAA07D28CA42191A50BC8BC205B898977FAD09AC51643FBC9CE6A740 BD13DCAA3347E3AF5EAA33A5E9EDA32DE4C454BBB4FF89C9CB2CD3AC85CD240E C3F5B5A8D2F3CDEAFCFEC4B5E9C0A17A2A66CBFBE70D2B7C39458FB8CABB76BA 0B71DE1E27DECE04FD9A9AD23E92C6605DBBDB279AC19A1ED913F6A8BDD2C61B CFCF36CC348916F72E8FCCCD21B3BF4B711A7857", " "),
		strings.Split("0ECS16O8PQZG4GMS FC8B6057667DCB78F388FAE2408FF1BD0CA775199444304C959BF2E652E8FF0550B9AC7C7E6ABFDD 2CDCA2BA30B6AF269C4E21AD0F060BB9928CDAD4EFDE9EED9AFF779D4F55E66E 616EC1BEBCBF0C2D8FA16A0692A4F3DBD6EEFAD3DFDFF16CDD07FC13CCEFCA4C AE0DC7368F2F3EECDFEECB2F62CE4EAB1F5CDC668CFFDEFFBC4B7E0B3FF6D7CA 8FEA516838B89A5B4766B346B2F53479E8464D23", " "),
		strings.Split("PIDHAYGC9AEJ20SJ 34F20E76FEE85680A79214D16DFE82117ABBE57A86DFE0362FA9312167EBD2501DA982711885960A ABBB41CF364ADD5FC0EBD7D4AFDCD912E06E974D22DD6F058FBDBF2B466C8D9B 0EA6ACEEAEAF463CAAFBD9CA4F32AE5FE6CF0BC9AF69CCB6C0F3B0E2DEBDF047 9FD931091A8E3C0062BAC9F64EDEF6EA7E3A5EA7A9AF64DBB20EB9D64E798EBE 6DFFCF8CCA189A2302FF472B160478431E1556AD", " "),
		strings.Split("AXK2OWGZ0DW88EHM C4A7D93AE52574A113246BAAF94D539C9D1B677D07928F4351E3D94BEF708D42DFD3D1354490E479 4D34AFCC03E756F7BE7CA187C1AD8DEB907746E27DEEEF5212CC8F6EDE5A6BDD 02A9B0C50C17E4CDD15D9F6FB67314D48BA7E823BA4612F43BF32FEB6ECCA78B 6FDBC6B31BBB2ACFFB21A210EFCFB7BE2AABED7D24ECBB3BA086E2CE9AF1A4F3 74C492615B1C5999B8CB7A141577C3065B12E334", " "),
	}

	for _, row := range rows {
		username := row[0]
		sessionKey := hexToByteArray(row[1], true)
		clientPublicKey := hexToByteArray(row[2], true).LittleEndian().BigInt()
		serverPublicKey := hexToByteArray(row[3], true).LittleEndian().BigInt()
		salt := hexToByteArray(row[4], true)
		expected := hexToByteArray(row[5], true).LittleEndian()

		assert.Equal(
			t,
			expected,
			calcClientProof(username, sessionKey, clientPublicKey, serverPublicKey, salt),
		)
		// fmt.Printf("%x\n", calcClientProof(username, sessionKey, clientPublicKey, serverPublicKey, salt).Bytes())
	}
}
