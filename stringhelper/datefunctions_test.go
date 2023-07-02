package stringhelper

import "testing"

func TestExtractDate(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		// the table itself
		{"POS 20230615", "20230615_WS_MVTEST_6003792356_20230615.txt", "20230615"},
		{"POS 20230630 ", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_20230630_08.log", "20230630"},
		{"POS 2023-06-30", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023-06-30_08.log", "20230630"},
		{"POS 2023_06_30", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023_06_30_08.log", "20230630"},
		{"POS 2023.06.30", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023.06.30_08.log", "20230630"},
		{"NEG 2023+06+30", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023+06+30_08.log", ""},
		{"NEG 2023_0630", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023-0630_08.log", ""},
		{"NEG 2023_0630", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_202306-30_08.log", ""},
		{"NEG 2023+06+30", "test-dummy1_XYZ_DF_ABC_JohnDoe_A_2023+06+30_08.log", ""},
		{"POS 2023.06.30", "test-dummy1_XYZ_DF_ABC_JohnDoeA2023063008.log", "20230630"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := ExtractDate(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s : %v", ans, tt.want, err)
			}
		})
	}
}
