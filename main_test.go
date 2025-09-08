// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/ericlagergren/decimal"
)

// This precision is used to compare results in tests.  Use two digits less
// than the max (34) to simplify rounding issues.
const defaultTestPrecision = 32

func TestRPN(t *testing.T) {
	ctx := decimal.Context128

	casetests := []struct {
		input     string
		want      *decimal.Big
		wantError bool
		precision int
	}{
		// Note: We use a "continuous" stack across operations.
		// Basic operations.
		{input: "1 2 +", want: bigUint(3)},
		{input: "4 5 8 + + +", want: bigUint(20)},
		{input: "5 -", want: bigUint(15)},
		{input: "3 /", want: bigUint(5)},
		{input: "3.5 *", want: bigFloat("17.5")},
		{input: "2 *", want: bigFloat("35.0")},
		{input: "chs", want: bigFloat("-35.0")},
		{input: "inv", want: bigFloat("-0.02857142857142857142857142857142857")},
		{input: "inv", want: bigFloat("-35.00000000000000000000000000000000")},
		{input: "chs", want: bigUint(35)},
		{input: "3 ^", want: bigUint(42875)},
		{input: "10 mod", want: bigUint(5)},
		{input: "fac", want: bigUint(120)},
		{input: "10 %", want: bigUint(12)},
		{input: "2 ^", want: bigUint(144)},
		{input: "sqr", want: bigUint(12)},
		{input: "3 ^", want: bigUint(1728)},
		{input: "cbr", want: bigUint(12)},
		{input: "c 1 2 3 4 sum", want: bigUint(10)},
		{input: "c 1 2 x", want: bigUint(1)},
		{input: "x", want: bigUint(2)},
		{input: "c", want: bigUint(0)},
		{input: "0.25 4 *", want: bigUint(1)},
		{input: ".2 5 *", want: bigUint(1)},
		{input: "$2,500.00 €3,500.00 +", want: bigUint(6000)},
		{input: "2 12345 ^", want: bigFloat(
			"164171010688258216356020741663906501410127235530735881272116" +
				"103087925094171390144280159034536439457734870419127140401667" +
				"195510331085657185332721089236401193044493457116299768844344" +
				"303479235489462436380672117015123283299131391904179287678259" +
				"173308536738761981139958654880852234908448338817289014166774" +
				"169869251339379828599748492918775437864739032217778051333882" +
				"990074116246281269364933724892342134504702491040016637557429" +
				"810893780765197418589477584716543480995722533317862352141459" +
				"217781316266211186486157019262080414077670264642736018426998" +
				"113523445732680856144329876972273300703392584997729207197971" +
				"083945700345494092400147186997307012069454068489589035676979" +
				"448169848060836924945824197706493306108258511936030341393221" +
				"586423523264452449403781993352421885094664052270795527632721" +
				"896121424813173522474674395886155092203404036730748474781710" +
				"715745446135468098139831824083259647919175273503681561172684" +
				"624283384438504776503000432241604550454374116320822227191911" +
				"322123484085063926350606342197146407841178028071147192533942" +
				"517270553513988142925976090769695456221159699052583533011331" +
				"652079347093098173086975483539274464023357456484465482927479" +
				"569437320368592222760278170306076733438801098370797675711274" +
				"671054970711442158930561684343135774118741594506702833147396" +
				"758825015850042983343690345185995956235143825771620543546030" +
				"664562647854656431302644574119873820215595718618624485232422" +
				"006575550007068883734241454686368856734496265385908809403972" +
				"494685137741122866896719678053937285818409751670320140501843" +
				"039224040735870096889596273419106389103662095318937990625980" +
				"136711988237421962315266686856089505981438440850638067589321" +
				"141759499017023839596858455548192000140085142294166987063499" +
				"024792681334843159790936321351919859758669569200541507612099" +
				"780909705198902176026219872201715422096090343686272984351441" +
				"594569506778041062663266799342793856313801540959815845788584" +
				"759033248828248561586450271172777240971795656082001848115815" +
				"260930521663167480173886064019118572778281516735157779555888" +
				"167787064432558595410843987446497881666288423233170060413025" +
				"924629950477303342180149398926073618582715358742250388958231" +
				"281694757980523791263699450732952325727664209947786063982561" +
				"775327638504516918570101319391698412388607603742484414748268" +
				"389669129118026878969735782286841116842656410574647607524418" +
				"900720328045377993386279808768990376289424757351052369393977" +
				"137871998119168898493037938756635621557623138404459266598837" +
				"784229325799838782026060481496865561757031839002257091802876" +
				"949248392744175669112242088439883248336310597001257385980776" +
				"961529351198877747193531054956881808332177946751404038228718" +
				"567911769630971553915410012677600002457982207465176670752102" +
				"117002773980548089696530972476439694599881281812973217265853" +
				"884727906535479745854085338851105144585481994156206497436745" +
				"899944877732531412541279014300324594890623941145509856940982" +
				"863769834430048120562966797907114102689879364945689860493474" +
				"954538422367719507882513166051007352994068319251450666676648" +
				"368200564329382998758875760414259654004977261309988267319806" +
				"354856051784553990936610634733375984159028722378614984450255" +
				"386315585631994503350002142910493190254825610707400589976364" +
				"985748467955131077971641882672895854571236368282811336220769" +
				"174784720113331269084746524204124263475054112841630933586166" +
				"195036115696469686075600480420563557567616835633252622327172" +
				"811002146392754445051182169805284630259703542633955126179520" +
				"113059629914229833688535925729676778028406897316106101038469" +
				"119090984567152591962365415039646394591503830797626339246986" +
				"057077758611413664914168745375266786298141171496573941614387" +
				"744125843685677063619782918759823106021054037757857761587472" +
				"240835040580447360544029064930412569943169729238102162312218" +
				"687930203068055400275795180972382856696655279408212344832"), precision: 34},
		{input: "10 6144 ^", want: bigFloat("1" + strings.Repeat("0", 6144)), precision: 34},
		{input: "10 6145 ^", want: bigFloat("+Infinity")},
		{input: "10 34 ^ 1 -", want: bigFloat("9999999999999999999999999999999999")},
		{input: "c", want: bigUint(0)},
		// Invalid operator should not cause changes to stack.
		{input: "foobar", want: bigUint(0)},

		// Trigonometric functions.
		{input: "deg 90 sin", want: bigFloat("1")},
		{input: "rad 90 PI * 180 / sin", want: bigUint(1)},
		{input: "deg 0 cos", want: bigUint(1)},
		{input: "rad 60 PI * 180 / cos", want: bigFloat("0.5")},
		{input: "deg 30 tan", want: func() *decimal.Big { // Sqrt(3) / 3)
			square := big().Quo(bigUint(1), bigUint(2))
			z := ctx.Pow(big(), bigUint(3), square)
			return z.Quo(z, bigUint(3))
		}()},
		{input: "rad 60 PI * 180 / tan", want: ctx.Pow(big(), bigUint(3), big().Quo(bigUint(1), bigUint(2)))}, // Sqrt(3)
		{input: "deg -1 180 * PI / asin", want: big().Quo(ctx.Pi(big()), bigUint(2)).SetSignbit(true)},        // -Pi / 2
		{input: "rad -1 asin", want: big().Quo(ctx.Pi(big()), bigUint(2)).SetSignbit(true)},                   // -Pi / 2
		{input: "rad 0 acos", want: big().Quo(ctx.Pi(big()), bigUint(2))},                                     // Pi / 2
		{input: "rad 1 acos", want: bigUint(0)},
		{input: "deg 0.5 180 * PI / acos", want: big().Quo(ctx.Pi(big()), bigUint(3))}, // Pi / 3
		{input: "rad 3 sqr atan", want: big().Quo(ctx.Pi(big()), bigUint(3))},          // Pi / 3
		{input: "deg 1 180 * PI / atan", want: big().Quo(ctx.Pi(big()), bigUint(4))},   // Pi / 4

		// Log functions
		{input: "c E ln", want: bigUint(1)},
		{input: "c 1000 log", want: bigUint(3)},

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: bigUint(0xffff)},
		{input: "0x0ff0 and", want: bigUint(0x0ff0)},
		{input: "0x1ee1 xor", want: bigUint(0x1111)},
		{input: "1 lshift", want: bigUint(0x2222)},
		{input: "1 lshift", want: bigUint(0x4444)},
		{input: "2 rshift", want: bigUint(0x1111)},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: bigUint(12419)},
		{input: "c", want: bigUint(0)},

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: bigUint(0xffff)},
		{input: "0x0ff0 and", want: bigUint(0x0ff0)},
		{input: "0x1ee1 xor", want: bigUint(0x1111)},
		{input: "1 lshift", want: bigUint(0x2222)},
		{input: "1 lshift", want: bigUint(0x4444)},
		{input: "2 rshift", want: bigUint(0x1111)},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: bigUint(12419)},
		{input: "d d", want: bigUint(0)},

		// Miscellaneous operations
		{input: "212 f2c", want: bigUint(100)},
		{input: "c2f", want: bigUint(212)},
		{input: "-40 f2c", want: bigFloat("-40")},
		{input: "-10 c2f", want: bigUint(14)},
		{input: "0 c2f", want: bigUint(32)},
		{input: "c", want: bigUint(0)},
		{input: "1 dup dup sum", want: bigUint(3)},
		{input: "c", want: bigUint(0)},
	}

	stack := &stackType{}

	for _, tt := range casetests {
		err := calc(stack, tt.input)
		if !tt.wantError {
			if err != nil {
				t.Fatalf("Got error %q, want no error", err)
			}
			precision := defaultTestPrecision
			if tt.precision != 0 {
				precision = tt.precision
			}
			got := decimal.WithPrecision(precision).Set(stack.top())
			want := decimal.WithPrecision(precision).Set(tt.want)
			if got.CmpTotal(want) != 0 {
				t.Fatalf("diff: input: %s, want: %s, got: %s", tt.input, want, got)
			}
			continue
		}

		// Here, we want to see an error.
		if err == nil {
			t.Errorf("Got no error, want error")
		}
	}
}

func TestFormatNumber(t *testing.T) {
	ctx := decimal.Context128

	casetests := []struct {
		base  int
		input *decimal.Big
		want  string
	}{
		// Decimal
		{10, bigUint(0), "0"},
		{10, bigUint(1), "1"},
		{10, bigUint(999), "999"},
		{10, bigUint(1000), "1000 (1,000)"},
		{10, bigUint(1000000), "1000000 (1,000,000)"},
		{10, bigUint(1000000000000000), "1000000000000000 (1,000,000,000,000,000)"},
		{10, bigFloat("10000.333333"), "10000.333333 (10,000.333333)"},
		{10, bigFloat("-10000.333333"), "-10000.333333 (-10,000.333333)"},
		{10, ctx.Quo(big(), bigUint(567), bigUint(999)), "0.567568"},
		{10, ctx.Pow(big(), bigUint(2), bigUint(64)), "18446744073709551616 (18,446,744,073,709,551,616)"},
		{10, ctx.Pow(big(), bigUint(2), bigUint(1234567890)), "Infinity"},
		{10, ctx.Quo(big(), bigUint(0), bigUint(0)), "NaN"},
		{10, ctx.Quo(big(), bigUint(1), bigUint(0)), "Infinity"},
		{10, ctx.Quo(big(), bigFloat("-1"), bigUint(0)), "-Infinity"},

		// Binary
		{2, bigUint(0b11111111), "0b11111111"},
		{2, big().Add(bigUint(0b11111111), bigFloat("0.5")), "0b11111111 (truncated from 255.5)"},
		{2, big().Add(bigUint(0b11111111), bigFloat("0.5")).SetSignbit(true), "-0b11111111 (truncated from -255.5)"},

		// Octal
		{8, bigUint(0377), "0377"},
		{8, big().Add(bigUint(0377), bigFloat("0.5")), "0377 (truncated from 255.5)"},
		{8, big().Add(bigUint(0377), bigFloat("0.5")).SetSignbit(true), "-0377 (truncated from -255.5)"},

		// Hex
		{16, bigUint(0xff), "0xff"},
		{16, big().Add(bigUint(0xff), bigFloat("0.5")).SetSignbit(true), "-0xff (truncated from -255.5)"},
	}
	for _, tt := range casetests {
		got := formatNumber(ctx, tt.input, tt.base, 6, false)
		if got != tt.want {
			t.Fatalf("diff: base: %d, input: %v, want: %q, got: %q", tt.base, tt.input, tt.want, got)
		}
		continue
	}
}

func Example_main() {
	os.Args = []string{"rpn", "1", "2", "3", "+", "+", "6", "-"}
	main()
	// Output: 0
}
