package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/imgproc"
	"github.com/ajdnik/imghash/similarity"
)

var colorMomentCalculateTests = []struct {
	filename   string
	hash       hashtype.Float64
	width      uint
	height     uint
	resizeType imgproc.ResizeType
	kernel     int
	sigma      float64
}{
	{"assets/lena.jpg", hashtype.Float64{0.001682934993677429, 3.9036864281658695e-09, 3.0554937212039115e-12, 1.4691471346003442e-10, -9.080795312922549e-22, 8.897371996993657e-15, 2.977306503666137e-21, 0.0012797434138170877, 4.388911176819886e-10, 3.0096420635170045e-12, 7.46410757760574e-12, 3.524037064463721e-23, 1.3906037391922022e-16, -3.1092091144073665e-24, 0.0009276508011411202, 3.2348966268694354e-09, 6.266786680370576e-13, 1.2595760568413167e-12, -1.108504738482668e-24, -6.834796243064575e-17, -1.534446982575172e-25, 0.0013354710297328763, 6.785675715722817e-09, 1.1937549471672312e-12, 1.1508102249222317e-11, -1.9456273772603476e-24, -8.778039449590288e-16, 4.260993662688662e-23, 0.0009969331042253649, 2.5979805607407044e-10, 2.407192815574133e-13, 1.220985948388969e-13, -1.983273625570263e-26, 6.064435932101452e-19, -6.695730840743158e-27, 0.0013973374569736076, 1.915611857766067e-09, 3.205292709962961e-13, 6.951279189437352e-13, -3.2810501250429035e-25, 3.0408967147378336e-17, -2.9939722908067632e-27}, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/baboon.jpg", hashtype.Float64{0.0024918588682241983, 2.6668418711805042e-08, 4.7072384155490496e-11, 2.830058934952091e-11, 3.462076470427119e-22, -4.50702081334241e-15, -9.731964400970534e-22, 0.0016590587234915884, 5.054257634038301e-09, 3.545330848595035e-11, 1.8371041544595362e-11, 5.55918598879393e-23, 1.3017502062374958e-15, 4.655372607399865e-22, 0.0009649767170685414, 3.4501227635866005e-09, 4.857582036353572e-13, 1.4973890413977926e-12, -1.2509081415494355e-24, 3.721800872696929e-17, 2.571316070355979e-25, 0.0012642851807576274, 3.9509158383639425e-09, 8.160680439530165e-13, 8.980273326874563e-13, -5.132728471031725e-25, -3.0409293931322725e-17, -5.723296844933609e-25, 0.001220347132745941, 5.975055105849723e-10, 1.770272494676388e-14, 4.0416882907207656e-12, -6.582886393254827e-25, 8.709067220205253e-17, 8.575690996612257e-25, 0.0013489776178140863, 7.386006701172241e-10, 4.2582323039863036e-13, 1.1871388847094712e-12, -7.157081255668091e-26, -2.0818832144747544e-17, 8.410077725967483e-25}, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/cat.jpg", hashtype.Float64{0.007840559221738475, 2.2034498359157744e-07, 8.684939637320095e-09, 1.3195340752596305e-08, 1.361994698531565e-16, -3.128016142811829e-12, 3.7465225373114406e-17, 0.0016131824928845837, 4.150246244542602e-08, 7.635890095045098e-11, 2.1849001025983358e-11, 4.1135742177399834e-22, -4.4427924315240715e-15, 7.91976998613191e-22, 0.001045548809225742, 2.772613007864512e-09, 1.337261180170139e-11, 1.6167463698622338e-11, 7.006891104009164e-24, 2.9211456863128017e-16, 2.376195809939422e-22, 0.0011816509074528167, 2.9408079432297224e-09, 3.061401364698586e-11, 3.893655593746319e-11, 2.84267196907831e-22, 9.85794678916472e-16, 1.3139013030470254e-21, 0.0011755591858654255, 1.4306034607639147e-10, 5.391803318981961e-13, 3.8112496595287077e-13, 1.6424203326396093e-25, -2.0998441435316273e-18, 5.360968121680658e-26, 0.0015080347084822899, 1.7821772814069935e-10, 1.1621533470920317e-13, 5.113034793145402e-13, -1.1298961074948206e-25, -2.3943853547045442e-18, 5.261115642777516e-26}, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/monarch.jpg", hashtype.Float64{0.002662394563829838, 8.300904725182467e-08, 1.4588163678577022e-09, 1.7799284933778596e-09, 1.972364977368551e-18, -4.96783175405463e-13, 2.0823418136581693e-18, 0.001250559295148185, 5.920860948159999e-10, 2.5370006565257617e-12, 2.792092761615076e-12, -5.124481646000649e-24, 4.841119239278186e-19, 5.381582530170822e-24, 0.0011622437650202432, 6.743101958369291e-09, 4.991671383226105e-12, 4.0181156557344116e-12, -6.951096541109391e-24, -7.802912877081211e-17, 1.6598476680525018e-23, 0.0015431023988322979, 1.551282558717626e-08, 6.779035641247579e-12, 2.4266605613294086e-12, -9.390163243411213e-24, -2.91505768311683e-16, -2.948916518638777e-24, 0.0010884880068615325, 2.619343498884856e-10, 3.561326317464918e-13, 7.588252697169075e-13, -9.596843455711414e-26, 1.1102086622370933e-17, 3.82622429862677e-25, 0.0015302253016975893, 5.386735944362092e-10, 4.0989260589827466e-14, 1.2472451390685164e-12, -2.777226964511623e-25, 1.103695419076888e-17, -4.898069005701871e-26}, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/peppers.jpg", hashtype.Float64{0.003941680044858082, 4.051369372556929e-07, 2.811861997021288e-11, 5.705112955860052e-10, 5.579971159763003e-20, -2.85837060590033e-13, -4.591066910351847e-20, 0.0010539871619496019, 3.1449666374563825e-09, 4.2502473588362055e-12, 9.629418181698208e-13, 1.2063026456570062e-24, 3.6703906742718216e-17, 1.5296589062296712e-24, 0.0010155149735782548, 2.527290610287061e-10, 1.223285873305063e-12, 5.069661814756587e-12, -1.8058653486884563e-24, 4.7566378252801606e-17, 1.2495209739304686e-23, 0.0013946292126994918, 1.6655999503876514e-08, 2.0374829552924317e-11, 2.7089651453782583e-12, 1.824127748336085e-23, 2.958882060256322e-16, -8.503040144541254e-24, 0.001128283975484718, 4.641665656093276e-09, 7.671235436760118e-12, 2.7676706136901324e-12, -8.148312761650637e-24, 1.2864173499951257e-16, -9.810093860745394e-24, 0.0016928305170451807, 1.8973805484919894e-09, 2.0274681685923785e-12, 1.051845804853099e-12, 1.6855575276247788e-25, 3.631772647036622e-17, -1.5267731530089948e-24}, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/tulips.jpg", hashtype.Float64{0.0025293177387043795, 1.4356598602797603e-08, 1.8704987719836307e-10, 1.6439303696445385e-12, -4.875692849536562e-24, -6.192153867552066e-17, 2.841197567643476e-23, 0.0014781139981409099, 3.637271003609577e-10, 2.917256126972781e-12, 7.791183825470193e-12, 3.6240544737932674e-23, 6.275212749438943e-17, -8.143937096806239e-24, 0.0010970144264479047, 5.188434119386236e-09, 4.112809393169052e-11, 2.4301287246914963e-12, 3.092354101411916e-24, 1.3648878736581163e-16, -2.4097185129778304e-23, 0.0013968503622165257, 7.0323173110166205e-09, 1.015742985374534e-10, 4.539201036842881e-12, 8.445279097508376e-23, 1.5000234971351897e-16, -4.865894404868708e-23, 0.0012375374902843692, 1.6038385709207788e-10, 1.4535056215353631e-12, 3.166203453449624e-14, 4.413719495194982e-27, -3.7221013570936514e-19, 5.162786191531114e-27, 0.0014163020274715553, 2.1211023930238313e-10, 2.631973957869045e-12, 1.81345828005539e-13, 9.398585348571086e-27, -2.3535641028507513e-18, -1.249328608255329e-25}, 512, 512, imgproc.Bicubic, 3, 0},
}

func TestColorMoment_Calculate(t *testing.T) {
	for _, tt := range colorMomentCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewColorMomentWithParams(tt.width, tt.height, tt.resizeType, tt.kernel, tt.sigma)
			img, _ := imgproc.Read(tt.filename)
			if res := hash.Calculate(img); !res.Equal(tt.hash) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

var colorMomentDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imgproc.ResizeType
	kernel      int
	sigma       float64
}{
	{"assets/lena.jpg", "assets/cat.jpg", 0.006173268140224651, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/lena.jpg", "assets/monarch.jpg", 0.0010413351488863956, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/baboon.jpg", "assets/cat.jpg", 0.005352693298027109, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/peppers.jpg", "assets/baboon.jpg", 0.0016168943396224305, 512, 512, imgproc.Bicubic, 3, 0},
	{"assets/tulips.jpg", "assets/monarch.jpg", 0.00036101159922622084, 512, 512, imgproc.Bicubic, 3, 0},
}

func TestColorMoment_Distance(t *testing.T) {
	for _, tt := range colorMomentDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewColorMomentWithParams(tt.width, tt.height, tt.resizeType, tt.kernel, tt.sigma)
			img1, _ := imgproc.Read(tt.firstImage)
			img2, _ := imgproc.Read(tt.secondImage)
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist := similarity.L2Float64(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
