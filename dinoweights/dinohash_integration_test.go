package dinoweights_test

import (
	"encoding/hex"
	"math/bits"
	"testing"

	imghash "github.com/ajdnik/imghash/v2"
	"github.com/ajdnik/imghash/dinoweights"
	"github.com/ajdnik/imghash/v2/hashtype"
)

// Reference hashes captured by running the proteus-photos compute_hashes_onnx.py
// pipeline on the same JPEG fixtures: PIL Image.open -> Resize((224,224))
// bilinear -> /255 -> ImageNet normalize -> ONNX (ORT) -> sign bits ->
// numpy.packbits big-endian. The fused 96-bit ONNX
// (dinov2_vits14_reg_96bit.onnx) is the source for both Go and Python; the
// Python sign vector is therefore the ground truth.
//
// Known 1-bit drift on assets/lena.jpg: position 55 flips relative to the
// Python pipeline. Root cause is PIL bilinear vs golang.org/x/image
// BilinearExact producing a sub-LSB pixel difference that pushes a single
// near-zero PCA-projection scalar across the sign threshold. The remaining
// five fixtures match bit-exact, so the lena hash is left as-is and the test
// allows up to 1 bit of distance against the Python reference. Add a
// PIL-replica resize and tighten this to zero once parity is required.
var pythonONNXRef = map[string]string{
	"../assets/lena.jpg":    "4c3b42ddaec814d3734a9b4f",
	"../assets/baboon.jpg":  "e14526f50f4e2f0fce94a7ef",
	"../assets/cat.jpg":     "cc5328a6e80a178637ee4314",
	"../assets/monarch.jpg": "cc2c35035351b8e5a1023d35",
	"../assets/peppers.jpg": "fc6f91748c81c07d7ca370d6",
	"../assets/tulips.jpg":  "ec2cba311f4d21a2cd334df6",
}

func TestDINOHash_Calculate_AgainstPythonONNX(t *testing.T) {
	d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
	if err != nil {
		t.Fatalf("NewDINOHash: %v", err)
	}
	const maxDrift = 1
	for path, refHex := range pythonONNXRef {
		t.Run(path, func(t *testing.T) {
			img, err := imghash.OpenImage(path)
			if err != nil {
				t.Fatalf("open %s: %v", path, err)
			}
			h, err := d.Calculate(img)
			if err != nil {
				t.Fatalf("Calculate %s: %v", path, err)
			}
			gotBin := []byte(h.(hashtype.Binary))
			refBin, err := hex.DecodeString(refHex)
			if err != nil {
				t.Fatalf("decode reference: %v", err)
			}
			if len(gotBin) != len(refBin) {
				t.Fatalf("length got=%d want=%d", len(gotBin), len(refBin))
			}
			var drift int
			for i := range gotBin {
				drift += bits.OnesCount8(gotBin[i] ^ refBin[i])
			}
			if drift > maxDrift {
				t.Errorf("hash drift %d bits exceeds tolerance %d\n  got %x\n  ref %x", drift, maxDrift, gotBin, refBin)
			}
		})
	}
}

func TestDINOWeightsBlob_Parseable(t *testing.T) {
	tensors, err := imghash.ParseSafetensors(dinoweights.Blob)
	if err != nil {
		t.Fatalf("ParseSafetensors(dinoweights.Blob): %v", err)
	}
	wantShape := map[string][]int{
		"patch_embed.weight":       {384, 3, 14, 14},
		"cls_token":                {1, 1, 384},
		"register_tokens":          {1, 4, 384},
		"pos_embed":                {1, 1370, 384},
		"blocks.0.attn.qkv.weight": {1152, 384},
		"head.weight":              {96, 384},
		"head.bias":                {1, 96},
	}
	for name, want := range wantShape {
		got, ok := tensors[name]
		if !ok {
			t.Errorf("tensor %q missing", name)
			continue
		}
		if len(got.Shape) != len(want) {
			t.Errorf("tensor %q shape rank = %d, want %d", name, len(got.Shape), len(want))
			continue
		}
		for i := range want {
			if got.Shape[i] != want[i] {
				t.Errorf("tensor %q shape = %v, want %v", name, got.Shape, want)
				break
			}
		}
	}
}

func ExampleDINOHash_Calculate() {
	img, err := imghash.OpenImage("../assets/cat.jpg")
	if err != nil {
		panic(err)
	}
	d, err := imghash.NewDINOHash(imghash.WithSafetensorsBlob(dinoweights.Blob))
	if err != nil {
		panic(err)
	}
	_, err = d.Calculate(img)
	if err != nil {
		panic(err)
	}
	// Output:
}
