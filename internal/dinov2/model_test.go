package dinov2

import (
	"fmt"
	"maps"
	"strings"
	"testing"

	"github.com/ajdnik/imghash/v2/internal/vit"
)

// Synthetic weights keep loader tests independent of the optional 85 MB blob.
func modelTensors() map[string]Tensor {
	sizes := map[string]int{
		"patch_embed.weight": vit.Dim * vit.Channels * vit.PatchSize * vit.PatchSize,
		"patch_embed.bias":   vit.Dim, "cls_token": vit.Dim,
		"register_tokens": vit.NumRegisters * vit.Dim,
		"pos_embed":       PosEmbedNativeLen * vit.Dim,
		"norm.weight":     vit.Dim, "norm.bias": vit.Dim,
		"head.weight": vit.HashBits * vit.Dim, "head.bias": vit.HashBits,
	}
	for i := 0; i < vit.NumBlocks; i++ {
		for name, size := range map[string]int{
			"norm1.weight": vit.Dim, "norm1.bias": vit.Dim,
			"attn.qkv.weight": 3 * vit.Dim * vit.Dim, "attn.qkv.bias": 3 * vit.Dim,
			"attn.proj.weight": vit.Dim * vit.Dim, "attn.proj.bias": vit.Dim,
			"ls1.gamma": vit.Dim, "norm2.weight": vit.Dim, "norm2.bias": vit.Dim,
			"mlp.fc1.weight": vit.MLPHidden * vit.Dim, "mlp.fc1.bias": vit.MLPHidden,
			"mlp.fc2.weight": vit.Dim * vit.MLPHidden, "mlp.fc2.bias": vit.Dim,
			"ls2.gamma": vit.Dim,
		} {
			sizes[fmt.Sprintf("blocks.%d.%s", i, name)] = size
		}
	}
	tensors := make(map[string]Tensor, len(sizes))
	for name, size := range sizes {
		data := make([]float32, size)
		data[0] = float32(len(tensors) + 1)
		tensors[name] = Tensor{Data: data}
	}
	return tensors
}

func TestLoadModel(t *testing.T) {
	tensors := modelTensors()
	for _, withLayerScale := range []bool{true, false} {
		t.Run(fmt.Sprintf("layer_scale=%t", withLayerScale), func(t *testing.T) {
			tensors := maps.Clone(tensors)
			if !withLayerScale {
				for name := range tensors {
					if strings.HasSuffix(name, ".gamma") {
						delete(tensors, name)
					}
				}
			}
			m, err := LoadModel(tensors)
			if err != nil {
				t.Fatal(err)
			}
			check := func(name string, got []float32) {
				t.Helper()
				want := tensors[name].Data
				if len(got) != len(want) {
					t.Fatalf("%s length = %d, want %d", name, len(got), len(want))
				}
				if len(want) > 0 && &got[0] != &want[0] {
					t.Errorf("%s does not reference the supplied tensor", name)
				}
			}
			check("patch_embed.weight", m.PatchEmbedW)
			check("patch_embed.bias", m.PatchEmbedB)
			check("cls_token", m.CLSToken)
			check("register_tokens", m.RegTokens)
			check("norm.weight", m.NormGamma)
			check("norm.bias", m.NormBeta)
			check("head.weight", m.HeadW)
			check("head.bias", m.HeadB)
			if len(m.PosEmbed) != vit.PosEmbedLen*vit.Dim || m.PosEmbed[0] != tensors["pos_embed"].Data[0] {
				t.Fatal("position embedding must be resized and preserve the CLS position")
			}
			if len(m.Blocks) != vit.NumBlocks {
				t.Fatalf("blocks = %d, want %d", len(m.Blocks), vit.NumBlocks)
			}
			for i, b := range m.Blocks {
				prefix := fmt.Sprintf("blocks.%d.", i)
				if b.Dim != vit.Dim || b.Eps != vit.Eps || b.Attn.Dim != vit.Dim || b.Attn.NumHeads != vit.NumHeads || b.Attn.HeadDim != vit.HeadDim || b.MLP.Dim != vit.Dim || b.MLP.Hidden != vit.MLPHidden {
					t.Fatalf("block %d has incorrect architecture", i)
				}
				for name, data := range map[string][]float32{
					"norm1.weight": b.Norm1Gamma, "norm1.bias": b.Norm1Beta,
					"attn.qkv.weight": b.Attn.QKVWeight, "attn.qkv.bias": b.Attn.QKVBias,
					"attn.proj.weight": b.Attn.ProjWeight, "attn.proj.bias": b.Attn.ProjBias,
					"ls1.gamma": b.LS1Gamma, "norm2.weight": b.Norm2Gamma, "norm2.bias": b.Norm2Beta,
					"mlp.fc1.weight": b.MLP.FC1Weight, "mlp.fc1.bias": b.MLP.FC1Bias,
					"mlp.fc2.weight": b.MLP.FC2Weight, "mlp.fc2.bias": b.MLP.FC2Bias,
					"ls2.gamma": b.LS2Gamma,
				} {
					check(prefix+name, data)
				}
			}
		})
	}
}

func TestLoadModel_InvalidTensors(t *testing.T) {
	tensors := modelTensors()
	for name, tensor := range maps.Clone(tensors) {
		for _, missing := range []bool{true, false} {
			if missing && strings.HasSuffix(name, ".gamma") {
				continue // LayerScale tensors are optional.
			}
			t.Run(fmt.Sprintf("%s/missing=%t", name, missing), func(t *testing.T) {
				if missing {
					delete(tensors, name)
				} else {
					tensors[name] = Tensor{Data: tensor.Data[:len(tensor.Data)-1]}
				}
				defer func() { tensors[name] = tensor }()
				m, err := LoadModel(tensors)
				if m != nil || err == nil || !strings.Contains(err.Error(), fmt.Sprintf("%q", name)) {
					t.Fatalf("LoadModel returned model=%v, error=%v; want error naming %q", m != nil, err, name)
				}
				want := "missing tensor"
				if !missing {
					want = fmt.Sprintf("has %d elements, want %d", len(tensor.Data)-1, len(tensor.Data))
				}
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want %q", err, want)
				}
			})
		}
	}
}
