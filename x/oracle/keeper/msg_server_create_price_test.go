package keeper_test

import (
	reflect "reflect"

	math "cosmossdk.io/math"
	dogfoodkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/testdata"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination mock_validator_test.go -package keeper_test github.com/cosmos/cosmos-sdk/x/staking/types ValidatorI

var _ = Describe("MsgCreatePrice", func() {
	//	var operator1, operator2, operator3 sdk.ValAddress
	var consKey1, consKey2, consKey3 sdk.AccAddress
	var c *cache.Cache
	var p *Patches
	BeforeEach(func() {
		ks.Reset()
		Expect(ks.ms).ToNot(BeNil())

		privVal1 := mock.NewPV()
		pubKey1, _ := privVal1.GetPubKey()
		//		operator1 = sdk.ValAddress(pubKey1.Address())
		consKey1 = sdk.AccAddress(pubKey1.Address())

		privVal2 := mock.NewPV()
		pubKey2, _ := privVal2.GetPubKey()
		consKey2 = sdk.AccAddress(pubKey2.Address())

		privVal3 := mock.NewPV()
		pubKey3, _ := privVal3.GetPubKey()
		consKey3 = sdk.AccAddress(pubKey3.Address())

		// TODO: remove monkey patch for test
		p = ApplyMethod(reflect.TypeOf(dogfoodkeeper.Keeper{}), "GetLastTotalPower", func(k dogfoodkeeper.Keeper, ctx sdk.Context) math.Int { return math.NewInt(3) })
		p.ApplyMethod(reflect.TypeOf(dogfoodkeeper.Keeper{}), "GetAllExocoreValidators", func(k dogfoodkeeper.Keeper, ctx sdk.Context) []dogfoodtypes.ExocoreValidator {
			return []dogfoodtypes.ExocoreValidator{
				{
					Address: pubKey1.Address(),
					Power:   1,
				},
				{
					Address: pubKey2.Address(),
					Power:   1,
				},
				{
					Address: pubKey3.Address(),
					Power:   1,
				},
			}
		})

		Expect(ks.ctx.BlockHeight()).To(Equal(int64(2)))
	})

	AfterEach(func() {
		ks.ctrl.Finish()
		if p != nil {
			p.Reset()
		}
	})

	Context("3 validators with 1 voting power each", func() {
		BeforeEach(func() {
			ks.ms.CreatePrice(ks.ctx, &types.MsgCreatePrice{
				// Creator:    operator1.String(),
				Creator:    consKey1.String(),
				FeederID:   1,
				Prices:     testdata.PS1,
				BasedBlock: 1,
				Nonce:      1,
			})

			c = keeper.GetCaches()
			pRes := &common.Params{}
			c.GetCache(cache.ItemP(pRes))
			p4Test := types.DefaultParams()
			p4Test.TokenFeeders[1].StartBaseBlock = 1
			Expect(*pRes).Should(BeEquivalentTo(p4Test))
		})

		It("success on 3rd message", func() {
			iRes := make([]*cache.ItemM, 0)
			c.GetCache(&iRes)
			Expect(iRes[0].Validator).Should(Equal(consKey1.String()))

			ks.ms.CreatePrice(ks.ctx, &types.MsgCreatePrice{
				Creator:    consKey2.String(),
				FeederID:   1,
				Prices:     testdata.PS2,
				BasedBlock: 1,
				Nonce:      1,
			},
			)
			ks.ms.CreatePrice(ks.ctx, &types.MsgCreatePrice{})
			c.GetCache(&iRes)
			Expect(len(iRes)).Should(Equal(2))

			ks.ms.CreatePrice(ks.ctx, &types.MsgCreatePrice{
				Creator:    consKey3.String(),
				FeederID:   1,
				Prices:     testdata.PS4,
				BasedBlock: 1,
				Nonce:      1,
			},
			)
			c.GetCache(&iRes)
			Expect(len(iRes)).Should(Equal(0))
			prices := ks.k.GetAllPrices(sdk.UnwrapSDKContext(ks.ctx))
			Expect(prices[0]).Should(BeEquivalentTo(types.Prices{
				TokenID:     1,
				NextRoundID: 2,
				PriceList: []*types.PriceTimeRound{
					{
						Price:     testdata.PTD2.Price,
						Decimal:   testdata.PTD2.Decimal,
						Timestamp: prices[0].PriceList[0].Timestamp,
						RoundID:   1,
					},
				},
			}))
		})
	})
})
