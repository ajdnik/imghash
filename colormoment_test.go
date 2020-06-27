package imghash_test

import (
	"fmt"

	"testing"

	. "github.com/ajdnik/imghash"
	"github.com/ajdnik/imghash/hashtype"
	"github.com/ajdnik/imghash/similarity"
)

var colorMomentCalculateTests = []struct {
	filename   string
	hash       hashtype.Float64
	width      uint
	height     uint
	resizeType ResizeType
	kernel     int
	sigma      float64
}{
	{"assets/lena.jpg", hashtype.Float64{0.0016755528384518568463246168676050729118, 0.0000000040498124784862561518617208341890, 0.0000000000024134051482390184451838278958, 0.0000000001457878792439884682850864994715, -0.0000000000000000000010212041931817059114, 0.0000000000000089244393987466171966711188, 0.0000000000000000000025367888103129747102, 0.0012822844318584815555273070941666446743, 0.0000000004155402064615737318827905300901, 0.0000000000030580975657581176886190089176, 0.0000000000074705626500816496955298566298, 0.0000000000000000000000355275454344103652, 0.0000000000000001352683188219727936497612, -0.0000000000000000000000035772492539636305, 0.0009242777066312995223190673854674059839, 0.0000000031991882503733676697522181678595, 0.0000000000006200000492699385018452209439, 0.0000000000012369679570006526337013853373, -0.0000000000000000000000010717475251155950, -0.0000000000000000666397495763894886243027, -0.0000000000000000000000001575244515713180, 0.0013303352874822847035085615630123356823, 0.0000000066766174985757278272022637564693, 0.0000000000011666112896083758944249963500, 0.0000000000112824521136343136844569209571, -0.0000000000000000000000017412784669605137, -0.0000000000000008539092964121314598418431, 0.0000000000000000000000408954521865070584, 0.0009960894533863210234536644804848037893, 0.0000000002619749147911207667144816461170, 0.0000000000002428528724284573006740958232, 0.0000000000001232924231615914259439878156, -0.0000000000000000000000000202323888608194, 0.0000000000000000006147717668894893670633, -0.0000000000000000000000000067674367706608, 0.0013976892912837281766902375323979867972, 0.0000000019348457413695052854176007232564, 0.0000000000003218096705293956832420737824, 0.0000000000006986006425889018010615248983, -0.0000000000000000000000003312129921317522, 0.0000000000000000307130927342116396088624, -0.0000000000000000000000000042441709589950}, 512, 512, Bicubic, 3, 0},
	{"assets/baboon.jpg", hashtype.Float64{0.0025082897664179112581783748225916497177, 0.0000000257634611647961016264032294833738, 0.0000000000510133875101602925445203291232, 0.0000000000341091701054653905242346692143, 0.0000000000000000000005238642764110667969, -0.0000000000000053241716451374538788247812, -0.0000000000000000000013228651669828509237, 0.0016691294052571079183089342024004508858, 0.0000000051232708490931239398285062754991, 0.0000000000386835933072046710881963879810, 0.0000000000189769962900475889543359042514, 0.0000000000000000000000342408134385749918, 0.0000000000000013433288257618152952265489, 0.0000000000000000000005130260953103055128, 0.0009624644881938247817018683427647829376, 0.0000000034351373235256599163087618225947, 0.0000000000004769850044700858595474888753, 0.0000000000014915752757281524826583000864, -0.0000000000000000000000012361229565864938, 0.0000000000000000360348944257926307629283, 0.0000000000000000000000002342046818409293, 0.0012595379743381219291714634422874041775, 0.0000000038918781065010252472275471423337, 0.0000000000008017665970884909338633574042, 0.0000000000008816901281130305446589237332, -0.0000000000000000000000004949334767305258, -0.0000000000000000295863260654084704094374, -0.0000000000000000000000005518845271546131, 0.0012200155852694709833483610594839774421, 0.0000000006022751759023855215462974636365, 0.0000000000000178832218560365947674887750, 0.0000000000040698525295798315852918394538, -0.0000000000000000000000006610754713482645, 0.0000000000000000879564926318580449761802, 0.0000000000000000000000008766525311560759, 0.0013491501705725126268597957590600344702, 0.0000000007418734486369334333224766320583, 0.0000000000004308008335591366955472725521, 0.0000000000011973350315716126938711460383, -0.0000000000000000000000000707314996084050, -0.0000000000000000210888028219638757792654, 0.0000000000000000000000008570135568232099}, 512, 512, Bicubic, 3, 0},
	{"assets/cat.jpg", hashtype.Float64{0.0077465609926574818117073206735767598730, 0.0000005490056530285862217738066491656834, 0.0000000069877748177321331001738656800871, 0.0000000127353966785635419123253212966840, 0.0000000000000001063018133784848749960077, -0.0000000000075681555627740688052955043380, -0.0000000000000000559783873366532040937818, 0.0016528967953893399477072190251192296273, 0.0000000449469556839834998365587517398650, 0.0000000000829838900338398333765696169939, 0.0000000000228359242388381055208319401843, 0.0000000000000000000005074084003436265970, -0.0000000000000048169570835017549534525440, 0.0000000000000000000008548376782063646934, 0.0010304361398509045370103232031055995321, 0.0000000026166711164867077681267444914590, 0.0000000000119210911693469797570657916728, 0.0000000000146601646148954183653597231234, 0.0000000000000000000000022810908071943823, 0.0000000000000002549103310145142953564893, 0.0000000000000000000001937921544489220316, 0.0011698663076284095997670053890260533080, 0.0000000028442823241015531677598310685744, 0.0000000000289759790362963057350070084117, 0.0000000000369540928704824408849853465579, 0.0000000000000000000002539913288813496961, 0.0000000000000009324863670873915845074420, 0.0000000000000000000011822657478604092660, 0.0011684952912512404366029983293628902175, 0.0000000001526426737051212864250004006312, 0.0000000000005482371440142733626626773247, 0.0000000000003831563521767345347866705879, 0.0000000000000000000000001657874748092768, -0.0000000000000000022323136287155900892094, 0.0000000000000000000000000579072324160354, 0.0014980889137293116816773697053122305078, 0.0000000001626361617850252595644909499282, 0.0000000000001004323040026869196973770332, 0.0000000000004791463656889229821920145975, -0.0000000000000000000000000923092124545013, -0.0000000000000000018283176126246665434859, 0.0000000000000000000000000502679391484742}, 512, 512, Bicubic, 3, 0},
	{"assets/monarch.jpg", hashtype.Float64{0.0026203358687636583469748874364313451224, 0.0000000954943915963917657912485531350455, 0.0000000013789342409239323935117611791409, 0.0000000016832123200888739566645880776528, 0.0000000000000000016830263579383413220534, -0.0000000000005042562895607684205983869332, 0.0000000000000000019347865221678934892355, 0.0012613147127821483332865692972291071783, 0.0000000005144842988402207042321010486768, 0.0000000000027726327757788891552718400737, 0.0000000000030940989413286471733212447102, -0.0000000000000000000000062179284637643632, 0.0000000000000000001415332411524857603348, 0.0000000000000000000000065928798597962658, 0.0011487167274111297413097254604963382008, 0.0000000064268073612985080067427374062982, 0.0000000000048066050136691691089794243026, 0.0000000000038598461690966932520346844971, -0.0000000000000000000000061924950937965533, -0.0000000000000000706366460525972588958513, 0.0000000000000000000000154291806230747646, 0.0015290182810008029652693029021293114056, 0.0000000149673359194826251587056973766862, 0.0000000000064418315356170528914541053197, 0.0000000000023085319805204277012528719531, -0.0000000000000000000000084837725244490705, -0.0000000000000002720838066264881866160059, -0.0000000000000000000000026979340181561140, 0.0010849407007777951174321007243861458846, 0.0000000002584843428468601258168398114768, 0.0000000000003643652780152868832092937946, 0.0000000000007669467598632618609706586260, -0.0000000000000000000000000897267079686544, 0.0000000000000000111642380553412019564679, 0.0000000000000000000000003953769506116549, 0.0015253961303048134239002120438044585171, 0.0000000005422194290729974638743401968237, 0.0000000000000363108276016162194445639509, 0.0000000000012825470825009409934957372388, -0.0000000000000000000000002622561258848904, 0.0000000000000000116941217883097770632535, -0.0000000000000000000000000884677616625858}, 512, 512, Bicubic, 3, 0},
	{"assets/peppers.jpg", hashtype.Float64{0.0039853495538433878278561373065258521819, 0.0000004379220941509448899332378624810191, 0.0000000000283336836155211145738681228455, 0.0000000005845186597744546173370021168857, 0.0000000000000000000632207983466212324847, -0.0000000000003236406859557054096328648255, -0.0000000000000000000407626021872640942402, 0.0010551220939480610537530846926301819622, 0.0000000031592146734708583409988078379753, 0.0000000000042600338101543893760258099090, 0.0000000000009673115811414270405230451098, 0.0000000000000000000000012507075125984020, 0.0000000000000000376082798271134883726930, 0.0000000000000000000000015137735560178692, 0.0010118376067164663693886961226553466986, 0.0000000002414967893857140296096756469716, 0.0000000000012118027153046795572486086933, 0.0000000000050036797418591600297702738543, -0.0000000000000000000000019096703756380388, 0.0000000000000000460459199734184184846373, 0.0000000000000000000000121722298180874861, 0.0013889588055681061778967588793420873117, 0.0000000163930935306663927345283723706337, 0.0000000000199409373039913105439517646331, 0.0000000000026662107325706140668740441130, 0.0000000000000000000000176378903603046958, 0.0000000000000002887217252170024421197986, -0.0000000000000000000000081761944043755270, 0.0011276538672013438832641973874615359819, 0.0000000046769519077581235021949259648023, 0.0000000000077180318351752769773669847620, 0.0000000000027843617616638563824962114028, -0.0000000000000000000000082556539468060673, 0.0000000000000001296941495665980786557364, -0.0000000000000000000000099220692517707846, 0.0016947366122165158364154535775014664978, 0.0000000019193769704875286617393436030361, 0.0000000000020535073040896644114857101640, 0.0000000000010677156859905629516664386245, 0.0000000000000000000000001784288091657855, 0.0000000000000000371186028945597997701266, -0.0000000000000000000000015708980230751724}, 512, 512, Bicubic, 3, 0},
	{"assets/tulips.jpg", hashtype.Float64{0.0024602151810521992458813400617145816796, 0.0000000141412524163601800765418389717494, 0.0000000001713855699137267204460115168829, 0.0000000000022311375235246538141891241101, -0.0000000000000000000000243263586006429633, -0.0000000000000001577866149638995696468715, 0.0000000000000000000000362178132747380295, 0.0015007515043890360013911777770090338890, 0.0000000004365702191663231976953156571582, 0.0000000000025529522295693679319803967381, 0.0000000000083984209774083188149863552661, 0.0000000000000000000000376813352261859759, 0.0000000000000001098986686875271112719555, -0.0000000000000000000000096129504490652542, 0.0010875002665389939260354168482081149705, 0.0000000049928295699397907242307931741272, 0.0000000000389993444230972065422070834058, 0.0000000000022847932176397427141156678124, 0.0000000000000000000000023493347447211458, 0.0000000000000001276055432562969052631330, -0.0000000000000000000000214391222868624837, 0.0013855889208886993779062946074986939493, 0.0000000067699910568487908057167860517191, 0.0000000000969522824706934896839667928395, 0.0000000000042985630365151820563807478442, 0.0000000000000000000000754754145281355300, 0.0000000000000001423146099218220081726824, -0.0000000000000000000000447674167061118414, 0.0012329017366625478081126887630603050638, 0.0000000001592555006884552777877563791391, 0.0000000000014485526684628839614324827750, 0.0000000000000319263704802605383323512997, 0.0000000000000000000000000044550641765413, -0.0000000000000000003716587064298564235553, 0.0000000000000000000000000052241441764621, 0.0014113409941877682121647019641841325210, 0.0000000002123880728554150023604062805902, 0.0000000000025893363954155318888348905633, 0.0000000000001785578447932620754885396448, 0.0000000000000000000000000093082649814668, -0.0000000000000000023168935198925849499644, -0.0000000000000000000000001210549449932434}, 512, 512, Bicubic, 3, 0},
}

func TestColorMoment_Calculate(t *testing.T) {
	for _, tt := range colorMomentCalculateTests {
		t.Run(tt.filename, func(t *testing.T) {
			hash := NewColorMomentWithParams(tt.width, tt.height, tt.resizeType, tt.kernel, tt.sigma)
			img, _ := ReadImageCV(tt.filename)
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
	resizeType  ResizeType
	kernel      int
	sigma       float64
}{
	{"assets/lena.jpg", "assets/cat.jpg", 0.0060886219783440835, 512, 512, Bicubic, 3, 0},
	{"assets/lena.jpg", "assets/monarch.jpg", 0.0010035467596423348, 512, 512, Bicubic, 3, 0},
	{"assets/baboon.jpg", "assets/cat.jpg", 0.005241874421998388, 512, 512, Bicubic, 3, 0},
	{"assets/peppers.jpg", "assets/baboon.jpg", 0.0016449495116979428, 512, 512, Bicubic, 3, 0},
	{"assets/tulips.jpg", "assets/monarch.jpg", 0.00037707969600449026, 512, 512, Bicubic, 3, 0},
}

func TestColorMoment_Distance(t *testing.T) {
	for _, tt := range colorMomentDistanceTests {
		t.Run(fmt.Sprintf("%v %v", tt.firstImage, tt.secondImage), func(t *testing.T) {
			hash := NewColorMomentWithParams(tt.width, tt.height, tt.resizeType, tt.kernel, tt.sigma)
			img1, _ := ReadImageCV(tt.firstImage)
			img2, _ := ReadImageCV(tt.secondImage)
			h1 := hash.Calculate(img1)
			h2 := hash.Calculate(img2)
			dist := similarity.L2Float64(h1, h2)
			if !dist.Equal(tt.distance) {
				t.Errorf("got %v, want %v", dist, tt.distance)
			}
		})
	}
}
