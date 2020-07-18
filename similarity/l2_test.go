package similarity_test

import (
	"fmt"
	"testing"

	"github.com/ajdnik/imghash/hashtype"
	. "github.com/ajdnik/imghash/similarity"
)

var l2Float64Tests = []struct {
	name  string
	hash1 hashtype.Float64
	hash2 hashtype.Float64
	out   Distance
}{
	{"same integer hashes", hashtype.Float64{1, 2, 3, 4, 5, 6, 7}, hashtype.Float64{1, 2, 3, 4, 5, 6, 7}, Distance(0)},
	{"different integer hashes", hashtype.Float64{2, 2, 3, 4, 5, 6, 7}, hashtype.Float64{1, 2, 3, 4, 5, 6, 7}, Distance(1)},
	{"same decimal hashes", hashtype.Float64{0.001682934993677429, 3.9036864281658695e-09, 3.0554937212039115e-12, 1.4691471346003442e-10, -9.080795312922549e-22, 8.897371996993657e-15, 2.977306503666137e-21, 0.0012797434138170877, 4.388911176819886e-10, 3.0096420635170045e-12, 7.46410757760574e-12, 3.524037064463721e-23, 1.3906037391922022e-16, -3.1092091144073665e-24, 0.0009276508011411202, 3.2348966268694354e-09, 6.266786680370576e-13, 1.2595760568413167e-12, -1.108504738482668e-24, -6.834796243064575e-17, -1.534446982575172e-25, 0.0013354710297328763, 6.785675715722817e-09, 1.1937549471672312e-12, 1.1508102249222317e-11, -1.9456273772603476e-24, -8.778039449590288e-16, 4.260993662688662e-23, 0.0009969331042253649, 2.5979805607407044e-10, 2.407192815574133e-13, 1.220985948388969e-13, -1.983273625570263e-26, 6.064435932101452e-19, -6.695730840743158e-27, 0.0013973374569736076, 1.915611857766067e-09, 3.205292709962961e-13, 6.951279189437352e-13, -3.2810501250429035e-25, 3.0408967147378336e-17, -2.9939722908067632e-27}, hashtype.Float64{0.001682934993677429, 3.9036864281658695e-09, 3.0554937212039115e-12, 1.4691471346003442e-10, -9.080795312922549e-22, 8.897371996993657e-15, 2.977306503666137e-21, 0.0012797434138170877, 4.388911176819886e-10, 3.0096420635170045e-12, 7.46410757760574e-12, 3.524037064463721e-23, 1.3906037391922022e-16, -3.1092091144073665e-24, 0.0009276508011411202, 3.2348966268694354e-09, 6.266786680370576e-13, 1.2595760568413167e-12, -1.108504738482668e-24, -6.834796243064575e-17, -1.534446982575172e-25, 0.0013354710297328763, 6.785675715722817e-09, 1.1937549471672312e-12, 1.1508102249222317e-11, -1.9456273772603476e-24, -8.778039449590288e-16, 4.260993662688662e-23, 0.0009969331042253649, 2.5979805607407044e-10, 2.407192815574133e-13, 1.220985948388969e-13, -1.983273625570263e-26, 6.064435932101452e-19, -6.695730840743158e-27, 0.0013973374569736076, 1.915611857766067e-09, 3.205292709962961e-13, 6.951279189437352e-13, -3.2810501250429035e-25, 3.0408967147378336e-17, -2.9939722908067632e-27}, Distance(0)},
	{"different decimal hashes", hashtype.Float64{0.0024918588682241983, 2.6668418711805042e-08, 4.7072384155490496e-11, 2.830058934952091e-11, 3.462076470427119e-22, -4.50702081334241e-15, -9.731964400970534e-22, 0.0016590587234915884, 5.054257634038301e-09, 3.545330848595035e-11, 1.8371041544595362e-11, 5.55918598879393e-23, 1.3017502062374958e-15, 4.655372607399865e-22, 0.0009649767170685414, 3.4501227635866005e-09, 4.857582036353572e-13, 1.4973890413977926e-12, -1.2509081415494355e-24, 3.721800872696929e-17, 2.571316070355979e-25, 0.0012642851807576274, 3.9509158383639425e-09, 8.160680439530165e-13, 8.980273326874563e-13, -5.132728471031725e-25, -3.0409293931322725e-17, -5.723296844933609e-25, 0.001220347132745941, 5.975055105849723e-10, 1.770272494676388e-14, 4.0416882907207656e-12, -6.582886393254827e-25, 8.709067220205253e-17, 8.575690996612257e-25, 0.0013489776178140863, 7.386006701172241e-10, 4.2582323039863036e-13, 1.1871388847094712e-12, -7.157081255668091e-26, -2.0818832144747544e-17, 8.410077725967483e-25}, hashtype.Float64{0.007840559221738475, 2.2034498359157744e-07, 8.684939637320095e-09, 1.3195340752596305e-08, 1.361994698531565e-16, -3.128016142811829e-12, 3.7465225373114406e-17, 0.0016131824928845837, 4.150246244542602e-08, 7.635890095045098e-11, 2.1849001025983358e-11, 4.1135742177399834e-22, -4.4427924315240715e-15, 7.91976998613191e-22, 0.001045548809225742, 2.772613007864512e-09, 1.337261180170139e-11, 1.6167463698622338e-11, 7.006891104009164e-24, 2.9211456863128017e-16, 2.376195809939422e-22, 0.0011816509074528167, 2.9408079432297224e-09, 3.061401364698586e-11, 3.893655593746319e-11, 2.84267196907831e-22, 9.85794678916472e-16, 1.3139013030470254e-21, 0.0011755591858654255, 1.4306034607639147e-10, 5.391803318981961e-13, 3.8112496595287077e-13, 1.6424203326396093e-25, -2.0998441435316273e-18, 5.360968121680658e-26, 0.0015080347084822899, 1.7821772814069935e-10, 1.1621533470920317e-13, 5.113034793145402e-13, -1.1298961074948206e-25, -2.3943853547045442e-18, 5.261115642777516e-26}, Distance(0.005352693298027109)},
}

func TestL2Float64(t *testing.T) {
	for _, tt := range l2Float64Tests {
		t.Run(tt.name, func(t *testing.T) {
			res := L2Float64(tt.hash1, tt.hash2)
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleL2Float64() {
	hash1 := hashtype.Float64{-6.582886393254827e-25, 8.709067220205253e-17, 8.575690996612257e-25}
	hash2 := hashtype.Float64{7.006891104009164e-24, 2.9211456863128017e-16, 2.376195809939422e-22, 0.0011816509074528167, 2.9408079432297224e-09}
	hash3 := hashtype.Float64{-1.983273625570263e-26, 6.064435932101452e-19, -6.695730840743158e-27, 0.0013973374569736076, 1.915611857766067e-09, 3.205292709962961e-13}

	fmt.Println(L2Float64(hash1, hash2))
	fmt.Println(L2Float64(hash1, hash3))
	// Output:
	// 2.050238964293645e-16
	// 8.648422860884238e-17
}

var l2Uint8Tests = []struct {
	name  string
	hash1 hashtype.UInt8
	hash2 hashtype.UInt8
	out   Distance
}{
	{"same integer hashes", hashtype.UInt8{132, 0, 255, 247, 54, 127, 136, 143, 64, 77, 159, 158, 113, 146, 101, 142, 138, 156, 140, 89, 128, 124, 124, 179, 108, 137, 134, 122, 145, 126, 134, 129, 118, 146, 133, 130, 124, 133, 132, 133}, hashtype.UInt8{132, 0, 255, 247, 54, 127, 136, 143, 64, 77, 159, 158, 113, 146, 101, 142, 138, 156, 140, 89, 128, 124, 124, 179, 108, 137, 134, 122, 145, 126, 134, 129, 118, 146, 133, 130, 124, 133, 132, 133}, Distance(0)},
	{"different hashes", hashtype.UInt8{71, 68, 254, 38, 179, 87, 159, 70, 0, 106, 39, 66, 72, 101, 61, 65, 82, 82, 85, 53, 55, 71, 55, 60, 67, 86, 64, 58, 72, 68, 75, 70, 65, 76, 83, 69, 55, 68, 64, 73}, hashtype.UInt8{166, 246, 10, 0, 124, 194, 254, 203, 219, 156, 116, 176, 226, 154, 138, 184, 195, 174, 155, 143, 213, 154, 170, 209, 125, 152, 173, 167, 181, 170, 165, 183, 157, 179, 174, 162, 161, 171, 157, 194}, Distance(717.8217048822082)},
}

func TestL2Uint8(t *testing.T) {
	for _, tt := range l2Uint8Tests {
		t.Run(tt.name, func(t *testing.T) {
			res := L2UInt8(tt.hash1, tt.hash2)
			if !res.Equal(tt.out) {
				t.Errorf("got %v, want %v", res, tt.out)
			}
		})
	}
}

func ExampleL2UInt8() {
	hash1 := hashtype.UInt8{60, 67, 86, 64, 58, 72, 68, 75}
	hash2 := hashtype.UInt8{143, 213, 154, 170, 209, 125, 152, 173, 167, 181}
	hash3 := hashtype.UInt8{0, 255, 247, 54, 127}

	fmt.Println(L2UInt8(hash1, hash2))
	fmt.Println(L2UInt8(hash1, hash3))
	// Output:
	// 293.82818108547724
	// 264.05681206891825
}