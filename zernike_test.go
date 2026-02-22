package imghash_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/v2/hashtype"
	"github.com/ajdnik/imghash/v2/similarity"
)

var zernikeCalculateTests = []struct {
	filename   string
	hash       hashtype.Float64
	width      uint
	height     uint
	resizeType imghash.Interpolation
	degree     int
}{
	{"assets/lena.jpg", hashtype.Float64{0.11668305348115725, 0.05156799331790223, 0.08957103711423306, 0.11586025697655308, 0.08189314610341758, 0.042684364826453734, 0.11305945045923636, 0.08143104777533867, 0.0817759699322858, 0.07759822420389888, 0.09294426977458646, 0.021645837784853482, 0.11565074086702563, 0.05688327151867321, 0.061321149567903785, 0.10167170320859784, 0.046987108870450545, 0.015598041503948317, 0.0481076852475451, 0.0593927831958293, 0.022380406697562923, 0.047752603506985146, 0.07148469280766888, 0.07048796318090457}, 64, 64, imghash.Bilinear, 8},
	{"assets/baboon.jpg", hashtype.Float64{0.06159525364585991, 0.09337686546889012, 0.050536289503132985, 0.10434319065344656, 0.09209168953823782, 0.21235273896562354, 0.02182557541748101, 0.06411775556253041, 0.05482733579270972, 0.08836426453428163, 0.029556856072352006, 0.00020280732901036756, 0.0060089372318643435, 0.0773081497248603, 0.08160568841695014, 0.031969891186867876, 0.046085962458639, 0.051415933338780945, 0.061463985953197, 0.10949154883835412, 0.08611625723504393, 0.032026473881818104, 0.009421726938003328, 0.04009796713653944}, 64, 64, imghash.Bilinear, 8},
	{"assets/cat.jpg", hashtype.Float64{0.28306618125837657, 0.22391755269960378, 0.17031766007654178, 0.026314025786357168, 0.06899862690035397, 0.06386868700331225, 0.1697964729998082, 0.07315746226147049, 0.03705128230519558, 0.04404122821425214, 0.06622638147254005, 0.020260231120349435, 0.1183382739266235, 0.08363771806823428, 0.07496781411624184, 0.04634838687124962, 0.03298688655087804, 0.04685165371160982, 0.05853848471039265, 0.09210950918602738, 0.023739139633080893, 0.1054261943165092, 0.0390394371997422, 0.07267920046066043}, 64, 64, imghash.Bilinear, 8},
	{"assets/monarch.jpg", hashtype.Float64{0.04477607848831568, 0.12677185835238777, 0.10266295016843405, 0.03183665304043582, 0.08834299884144185, 0.17508353508205277, 0.014472757292212755, 0.015411795568623942, 0.050496207673955254, 0.0802734508291114, 0.045334250431127744, 0.06657451032236844, 0.025246091085086745, 0.059434266053859174, 0.06584489585728268, 0.07097351739773158, 0.08566780531565811, 0.04273381842173576, 0.03444618969366611, 0.0750550777590373, 0.04867382101513085, 0.003643525976804611, 0.061282891607232895, 0.044101885498916985}, 64, 64, imghash.Bilinear, 8},
	{"assets/peppers.jpg", hashtype.Float64{0.07725142394799082, 0.14216630263413413, 0.1184614100959899, 0.10664718196721855, 0.05785443196507941, 0.11616041791398905, 0.19091851031987048, 0.07581750473623498, 0.057577360146210003, 0.0700116750465979, 0.11515907317567785, 0.1569578602584502, 0.11819462261489701, 0.21898842504449073, 0.04431222952898577, 0.10254145163752351, 0.06781903745914254, 0.16142950694100616, 0.03901099012355343, 0.02295872264855506, 0.09129582769472767, 0.07744783986339221, 0.03043175204765063, 0.021333654450331342}, 64, 64, imghash.Bilinear, 8},
	{"assets/tulips.jpg", hashtype.Float64{0.027261825089668873, 0.22182499405971665, 0.0570153978575255, 0.0772499250201467, 0.11020785880281243, 0.2122878065644418, 0.11465733580577779, 0.005478118071578849, 0.02238971939480202, 0.046968002853662746, 0.0842522900633989, 0.1566057641781387, 0.13392821176657774, 0.04546025752776038, 0.15038142066581675, 0.1557370058717295, 0.05417114947370435, 0.19249324952059077, 0.08158850512255136, 0.20542962269103812, 0.12434589350373026, 0.06207449137611067, 0.11910216764589134, 0.030314128360730892}, 64, 64, imghash.Bilinear, 8},
}

func float64ApproxEqual(a, b hashtype.Float64, maxDiff float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > maxDiff {
			return false
		}
	}
	return true
}

func TestZernike_Calculate(t *testing.T) {
	for _, tt := range zernikeCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash, err := imghash.NewZernike(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType), imghash.WithDegree(tt.degree))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img, err := imghash.OpenImage(tt.filename)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.filename, err)
			}
			result, err := hash.Calculate(img)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			res := result.(hashtype.Float64)
			if !float64ApproxEqual(res, tt.hash, 1e-9) {
				t.Errorf("got %v, want %v", res, tt.hash)
			}
		})
	}
}

func ExampleZernike_Calculate() {
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	z, err := imghash.NewZernike()
	if err != nil {
		panic(err)
	}
	hash, err := z.Calculate(img)
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}

var zernikeDistanceTests = []struct {
	firstImage  string
	secondImage string
	distance    similarity.Distance
	width       uint
	height      uint
	resizeType  imghash.Interpolation
	degree      int
}{
	{"assets/lena.jpg", "assets/cat.jpg", 0.300630123737467, 64, 64, imghash.Bilinear, 8},
	{"assets/lena.jpg", "assets/monarch.jpg", 0.264993575543639, 64, 64, imghash.Bilinear, 8},
	{"assets/baboon.jpg", "assets/cat.jpg", 0.399675840586009, 64, 64, imghash.Bilinear, 8},
	{"assets/peppers.jpg", "assets/baboon.jpg", 0.375665121373897, 64, 64, imghash.Bilinear, 8},
	{"assets/tulips.jpg", "assets/monarch.jpg", 0.344466318423521, 64, 64, imghash.Bilinear, 8},
}

func TestZernike_Distance(t *testing.T) {
	for _, tt := range zernikeDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash, err := imghash.NewZernike(imghash.WithSize(tt.width, tt.height), imghash.WithInterpolation(tt.resizeType), imghash.WithDegree(tt.degree))
			if err != nil {
				t.Fatalf("failed to create hasher: %v", err)
			}
			img1, err := imghash.OpenImage(tt.firstImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.firstImage, err)
			}
			img2, err := imghash.OpenImage(tt.secondImage)
			if err != nil {
				t.Fatalf("failed to open %s: %v", tt.secondImage, err)
			}
			h1, err := hash.Calculate(img1)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.firstImage, err)
			}
			h2, err := hash.Calculate(img2)
			if err != nil {
				t.Fatalf("failed to calculate hash for %s: %v", tt.secondImage, err)
			}
			dist, err := hash.Compare(h1, h2)
			if err != nil {
				t.Fatalf("failed to compute distance: %v", err)
			}
			if math.Abs(float64(dist)-float64(tt.distance)) > 1e-9 {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}

func TestNewZernike_Defaults(t *testing.T) {
	z, err := imghash.NewZernike()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := z.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	if h.Len() != 24 {
		t.Fatalf("expected hash length 24, got %d", h.Len())
	}
}

func TestNewZernike_CustomDegree(t *testing.T) {
	z, err := imghash.NewZernike(imghash.WithDegree(4))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	img, err := imghash.OpenImage("assets/cat.jpg")
	if err != nil {
		t.Fatalf("failed to open image: %v", err)
	}
	h, err := z.Calculate(img)
	if err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}
	if h.Len() != 8 {
		t.Fatalf("expected hash length 8, got %d", h.Len())
	}
}

func TestNewZernike_Errors(t *testing.T) {
	tests := []struct {
		name string
		opts []imghash.ZernikeOption
		err  error
	}{
		{"zero width", []imghash.ZernikeOption{imghash.WithSize(0, 64)}, imghash.ErrInvalidSize},
		{"zero height", []imghash.ZernikeOption{imghash.WithSize(64, 0)}, imghash.ErrInvalidSize},
		{"invalid interpolation", []imghash.ZernikeOption{imghash.WithInterpolation(imghash.Interpolation(999))}, imghash.ErrInvalidInterpolation},
		{"zero degree", []imghash.ZernikeOption{imghash.WithDegree(0)}, imghash.ErrInvalidDegree},
		{"negative degree", []imghash.ZernikeOption{imghash.WithDegree(-1)}, imghash.ErrInvalidDegree},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := imghash.NewZernike(tt.opts...)
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}
