package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MsgUpdateParams", func() {
	BeforeEach(func() {
		ks.Reset()
		Expect(ks.ms).ToNot(BeNil())
	})
	Context("", func() {
		It("update StartBaseBlock for TokenFeeder", func() {
			p := ks.k.GetParams(ks.ctx)
			p.TokenFeeders[1].StartBaseBlock = 10
			ks.k.SetParams(ks.ctx, p)
			p.TokenFeeders[1].StartBaseBlock = 5
			_, err := ks.ms.UpdateParams(ks.ctx, &types.MsgUpdateParams{
				Params: types.Params{
					TokenFeeders: []*types.TokenFeeder{
						{
							TokenID:        1,
							StartBaseBlock: 5,
						},
					},
				},
			})
			Expect(err).Should(BeNil())
			p = ks.k.GetParams(ks.ctx)
			Expect(p.TokenFeeders[1].StartBaseBlock).Should(BeEquivalentTo(5))
		})
		It("Add AssetID for Token", func() {
			_, err := ks.ms.UpdateParams(ks.ctx, &types.MsgUpdateParams{
				Params: types.Params{
					Tokens: []*types.Token{
						{
							Name:    "ETH",
							ChainID: 1,
							AssetID: "0x83e6850591425e3c1e263c054f4466838b9bd9e4_0x9ce1",
						},
					},
				},
			})
			Expect(err).Should(BeNil())
			p := ks.k.GetParams(ks.ctx)
			Expect(p.Tokens[1].AssetID).Should(BeEquivalentTo("0x83e6850591425e3c1e263c054f4466838b9bd9e4_0x9ce1"))
		})
	})
})
