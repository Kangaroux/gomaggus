package main

import (
	"encoding/csv"
	"encoding/hex"
	"os"
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
	rows := loadTestData("test_data/calculate_x.csv")

	for _, row := range rows {
		username := row[0]
		password := row[1]
		salt := hexToByteArray(row[2], false)
		expected := hexToByteArray(row[3], false).BigInt()

		assert.Equal(t, expected, calcX(username, password, salt))
	}
}

func Test_passVerify(t *testing.T) {
	rows := loadTestData("test_data/calculate_verifier.csv")

	for _, row := range rows {
		username := row[0]
		password := row[1]
		salt := hexToByteArray(row[2], false)
		expected := hexToByteArray(row[3], false).BigInt()

		assert.Equal(t, expected, passVerify(username, password, salt))
	}
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

		// fmt.Printf("%x\n", calcInterleave(S).Bytes())
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

// func Test_calcClientSessionKey(t *testing.T) {
// 	type testCase struct {
// 		username         string // Plaintext
// 		password         string // Plaintext
// 		clientPublicKey  string // Little endian hex
// 		clientPrivateKey string // Little endian hex
// 		serverPublicKey  string // Little endian hex
// 		salt             string // Big endian hex
// 		expected         string // Little endian hex
// 	}

// 	// First 10 testCases from:
// 	// https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/verification_values/calculate_server_session_key.txt
// 	testCases := []testCase{
// 		{"JZsyczxJwDVDXswZ", "ihPe779qnoDix6i5", "42ABC7A7CA3AF505EFF560EBAFA6ECE1C8AE65395A17EB39B98CDFEE1F5FAAEB", "28817DEEDCFA4DF43A0BD9CAB95E72CCA4B2FEECFABEEEE5F40861D4903B407E", "AAD74ABDDEA3BCCE8E4C2E1C8BA05BBAEF1F3DC1B52AAB2A2AC6F390B1276ABD", "9BBEE762EBCCE595E76EF891CD1A7F6C54CC8B65E30CA7CD3375C49433E0E3F7", "80D3FB86282E6A900517CF226BC00D7A4BD1402102648A6C1404C634E21E0BF78D322CE5D1B9DBE5"},
// 		{"lbOhuFSWnkUvb85X", "3Oxrd5k3tVLp46Zy", "CEC2B6F8D7C8B0792DAA6252138FE3CFC9838D7CB6BD4DBA1E3ABE4C7AAD6B7B", "025CBCF2CE5C1CED1BE41D89AC2B7A8DCC323FA59D9E14CFA40B28A13598FF7C", "83A0BFC06AF6E0782AEECAA2E9C2C26FBB7C8FD9C15BB1CDE672D05A6D1A319D", "234AD29E312E3629EA8E85F9CFEEBD394FCDBB164E23E4F145F3E64BFC971EFF", "A5B07A94B65534D8BAF559405C8AA1384427BA4767F8F624F808B5E7890D39D40DEDBD6ECA2DEBFA"},
// 		{"sFHQ6aoh3wfJTXaN", "2pBFKAuEByNkWqXF", "41CE6F062B4771FE3BDA667F43D0D53E5720E7E7CE16EFA91D6ACC169E279CE3", "4FC2843DC720DD76BDD3167AABE8F1E4BA190D883BE52B92B3C85CB83CCA10AC", "7CBFC0F46A7D67CFAD7E781C20E68DDAE0A5D3F05B7D4B3C2B8140A25F3EAD7E", "AEC9AD986CDF54D3E5DB4EEAFF9B2E38AD7EE2948FCAB6C91CD2DD1B6F219E3F", "2D68257416E3D2317A9111CB9219C88E16EDBB803CDFC7F25552AB833A04E851DABAF1A5611EFBC9"},
// 		{"UgvbPPR1Pvd2SDkL", "DS1LSWujYSwmUhEO", "C8DA4FF25BFA0C5EFAEBC9B14B021AAD83A0A4A66AFA5C8FDDCA1A521A2D668C", "FCEBCF90D7FE0CEEA8FABB4625DF31BC63A6096A9ED88C883DFAA7A650EE3D24", "5811ECFCC8AE8DFBF88499F796ACCAE1CDBCCF1E6B2735687AB40DEBBC08FB7C", "9AE0DBBC0D5BBBF1F2BADFA9018E5EBE511BBE59D52AAEB4670BCEDC9297FE2F", "4453847EB907A70A5BBE222EE1C381D4A5C548E3E32447BA98A4AB7E4B53E1A49FDCB72B931B80F5"},
// 		{"ncFm33aGlxKcn0QY", "bCuSAmJFX4e9KYFS", "8C4CA8A45FA68AC36D15C76DCB3C7BDB83F5EDBB0E5AF71C2A3F100F2CE3FCCC", "D16FDD5C7BCDDF9FD0A4A87AC484E1717FC5A823433316FBEF7E7436EA228743", "AF7AAA32FF6264BABDA84B8AE2FA1D7AFCEEE8DFBE5F59FE555ED0DC081D0BAB", "A73EEE0106A7EA6D6D5E5D8FBD7BBC5DDDCEBB0B882BB7FA84F1AD92F9845A53", "273FF5D625CC096E734ED237802990A47322BC7911CC1512EC301036B835E9E094DD246E1CEE124D"},
// 		{"zrfts7MNEdmFJZwM", "Vp1dOwk2VnDLa62O", "39BFCCEAADCD3F2E2AACBC13FF684D3D9F44DAEDD6EDE5FBF43AAC1BDD4D1C1B", "C427678ACEF7ABB4A8B54E9D9F93644A42EECEFEDE1BFB27FC2F24768BADCA1F", "8FB989DDBC4F15A4304E1CEFA6FACC0CB6F43CFAFF785FB45194ED7BD1CDBB55", "AFE2C9F7F977EB1BFAACD56CFEE25C4D7D087B340EA2A0F3BA1C0F541D5808DB", "B9870CFE1BE23EFFDB4FA8FD747D18413EAF4223D2F91FE8CC011B49D50618EF332BF98BE10D5943"},
// 		{"ULOcMyaT7MXceser", "pFNs1qclpBKs6g1I", "3BB3A4FFDCD585E9CCA95BFFFB1CDB1ABEC31B4C1D199BDE585D0BACD9CEF60B", "ABFAEF249E6F411DCCEAEE7C8EFBBD2BA87AB6DE4A759B8DB5C5ED1097A2E8EB", "5EEFF6AF6FDFF390F462349D221AE880FE3FDDC2E7FCCAD4FCCD0FB916EE6580", "E7DBF2EAEF5FABC6D2DD426E63869AEBC57C7D14CBCCC7DDE3389F1EDE48CE59", "3E73FF08D7B7FF403B22C1610F3242B618F1BA1C616DBC6597C28E06415BD9F4C1371E003D7622C4"},
// 		{"RyNiasaTUVw23N6n", "2Bf8oLRL4JjAmK5u", "C82FDE9C1D7CC05A72BD9587CD9BC6D79CE59F842E2BD4E266EBC78FA2B3CAAA", "AE5EABA1EEE08FF8167B433BDFB0E0B9FDBEE5A81A54A9C6D84EBD60FB25ADC6", "8D3C0C195ABAFDBA5E3DDBDA6345F15BB66DB51442AB55326EBA643ECAB609C0", "6B901F8D6E7D35BDBFBD4D5E71DE6DBEF7E7C6B0C037ADBF2104EAD79D3FAFEA", "83A793F2C7641EE68A73B849669F315189BCB11E24B9959248F69C0F70AC4905B34F4A3E01BFAB66"},
// 		{"oPSHo1D79V4487EZ", "boKq0rP7JCPDlCSX", "B7EBA995F9FAEF59C70FB1A55DF0EACFDE696CD8CC06E7BCF3DB9ADCEF9FACDB", "37291BFBBB0BE359B3C382C3140DE2AEE6EADDAEC500C3EA27DA95D7C31AABCE", "9DDB55AF8FAFBACE02D3CEF11CFBCCE5E8AFAC4AC39BE7B12C2F8E8FC8308F7F", "F07EB4C97729F4FBABBB55BBEC20AF2FE5C8F8CADEEBE80EF03A12E53EC7B096", "87E9708A96C09013D6D63898BE7DFC0478E38F1CDCB04693BE5204610765EDC29C88F11A99541E79"},
// 		{"L8SniqtpIDdyLv4p", "nhzoSdKmN3KlRKy6", "0307CC3C23E5DC2AD282FEFDE44D66CDFAAE0C7D7E1ABCA57CE5088EC0D6BDBB", "86EDFF57605F29FAD7A4D0ADA12D11414016D6DE7CA03CA3FFE71D5BE480C5CF", "85B8295EE4EB11FF21AFF5C104C20E908BBBED31DD53F8ACFD8FE83EE2A12FEE", "ADCDBB0CE9CCDCDFCC43D24E330CCF758F98B2E37DD6BB5DC14FFFEB6F8E6BBB", "CA0AE27695468B71CB92E462E6DAE418405F19A8C82FB26EB500742E80951330D73AEA688CC0F00F"},
// 	}

// 	for _, tc := range testCases {
// 		clientPublicKey := hexToByteArray(tc.clientPublicKey, false).BigInt()
// 		clientPrivateKey := hexToByteArray(tc.clientPrivateKey, false).BigInt()
// 		serverPublicKey := hexToByteArray(tc.serverPublicKey, false).BigInt()
// 		salt := hexToByteArray(tc.salt, true)
// 		expected := hexToByteArray(tc.expected, false)
// 		assert.Equal(
// 			t,
// 			expected,
// 			calcClientSessionKey(tc.username,
// 				tc.password,
// 				serverPublicKey,
// 				clientPrivateKey,
// 				clientPublicKey,
// 				salt,
// 			),
// 		)
// 	}
// }
