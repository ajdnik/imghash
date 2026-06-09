package dinov2

import (
	"fmt"

	"github.com/ajdnik/imghash/v2/internal/vit"
)

// LoadModel populates a *vit.Model from a parsed safetensors tensor map. It
// validates each tensor's shape against the architectural constants in the
// vit package and reports the first mismatch with the tensor's name. The
// pos_embed tensor is bicubically resampled from its native 37x37 grid to
// the 16x16 grid used at 224x224 inference resolution.
func LoadModel(tensors map[string]Tensor) (*vit.Model, error) {
	get := func(name string, want int) ([]float32, error) {
		t, ok := tensors[name]
		if !ok {
			return nil, fmt.Errorf("imghash: missing tensor %q", name)
		}
		if len(t.Data) != want {
			return nil, fmt.Errorf("imghash: tensor %q has %d elements, want %d", name, len(t.Data), want)
		}
		return t.Data, nil
	}
	optional := func(name string, want int) ([]float32, error) {
		t, ok := tensors[name]
		if !ok {
			return nil, nil
		}
		if len(t.Data) != want {
			return nil, fmt.Errorf("imghash: tensor %q has %d elements, want %d", name, len(t.Data), want)
		}
		return t.Data, nil
	}

	m := &vit.Model{}
	var err error
	if m.PatchEmbedW, err = get("patch_embed.weight", vit.Dim*vit.Channels*vit.PatchSize*vit.PatchSize); err != nil {
		return nil, err
	}
	if m.PatchEmbedB, err = get("patch_embed.bias", vit.Dim); err != nil {
		return nil, err
	}
	if m.CLSToken, err = get("cls_token", vit.Dim); err != nil {
		return nil, err
	}
	if m.RegTokens, err = get("register_tokens", vit.NumRegisters*vit.Dim); err != nil {
		return nil, err
	}
	rawPos, err := get("pos_embed", PosEmbedNativeLen*vit.Dim)
	if err != nil {
		return nil, err
	}
	m.PosEmbed = InterpolatePosEmbed(rawPos)

	m.Blocks = make([]*vit.Block, vit.NumBlocks)
	for i := 0; i < vit.NumBlocks; i++ {
		prefix := fmt.Sprintf("blocks.%d.", i)
		b := &vit.Block{Dim: vit.Dim, Eps: vit.Eps}
		if b.Norm1Gamma, err = get(prefix+"norm1.weight", vit.Dim); err != nil {
			return nil, err
		}
		if b.Norm1Beta, err = get(prefix+"norm1.bias", vit.Dim); err != nil {
			return nil, err
		}
		attn := &vit.Attention{Dim: vit.Dim, NumHeads: vit.NumHeads, HeadDim: vit.HeadDim}
		if attn.QKVWeight, err = get(prefix+"attn.qkv.weight", 3*vit.Dim*vit.Dim); err != nil {
			return nil, err
		}
		if attn.QKVBias, err = get(prefix+"attn.qkv.bias", 3*vit.Dim); err != nil {
			return nil, err
		}
		if attn.ProjWeight, err = get(prefix+"attn.proj.weight", vit.Dim*vit.Dim); err != nil {
			return nil, err
		}
		if attn.ProjBias, err = get(prefix+"attn.proj.bias", vit.Dim); err != nil {
			return nil, err
		}
		b.Attn = attn
		if b.LS1Gamma, err = optional(prefix+"ls1.gamma", vit.Dim); err != nil {
			return nil, err
		}
		if b.Norm2Gamma, err = get(prefix+"norm2.weight", vit.Dim); err != nil {
			return nil, err
		}
		if b.Norm2Beta, err = get(prefix+"norm2.bias", vit.Dim); err != nil {
			return nil, err
		}
		mlp := &vit.MLP{Dim: vit.Dim, Hidden: vit.MLPHidden}
		if mlp.FC1Weight, err = get(prefix+"mlp.fc1.weight", vit.MLPHidden*vit.Dim); err != nil {
			return nil, err
		}
		if mlp.FC1Bias, err = get(prefix+"mlp.fc1.bias", vit.MLPHidden); err != nil {
			return nil, err
		}
		if mlp.FC2Weight, err = get(prefix+"mlp.fc2.weight", vit.Dim*vit.MLPHidden); err != nil {
			return nil, err
		}
		if mlp.FC2Bias, err = get(prefix+"mlp.fc2.bias", vit.Dim); err != nil {
			return nil, err
		}
		b.MLP = mlp
		if b.LS2Gamma, err = optional(prefix+"ls2.gamma", vit.Dim); err != nil {
			return nil, err
		}
		m.Blocks[i] = b
	}

	if m.NormGamma, err = get("norm.weight", vit.Dim); err != nil {
		return nil, err
	}
	if m.NormBeta, err = get("norm.bias", vit.Dim); err != nil {
		return nil, err
	}
	if m.HeadW, err = get("head.weight", vit.HashBits*vit.Dim); err != nil {
		return nil, err
	}
	if m.HeadB, err = get("head.bias", vit.HashBits); err != nil {
		return nil, err
	}
	return m, nil
}
