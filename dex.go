package main

import (
	"fmt"
	"strings"
	"strconv"
	"math"
	"encoding/hex"

	d "github.com/deroholic/derogo"
	"github.com/yourbasic/graph"
)

type Token struct {
        n int
        contract string
        decimals int
        bridgeFee uint64
        bridgeable bool
        swapable bool
        native_symbol string
        name string
        eth_contract string
}

type Pair struct {
        contract string
        fee uint64
        val1 uint64
        val2 uint64
        sharesOutstanding uint64
        adds uint64
        rems uint64
        swaps uint64
}

var tokens map[string]Token
var pairs map[string]Pair
var tokenList []string
var tokenGraph *graph.Mutable

var bridgeRegistry string
var swapRegistry string

func getTokens() {
        tokens = make(map[string]Token)
        n := 0

        // bridgeable tokens
        bridgeVars, bridgeValid := d.DeroGetVars(bridgeRegistry)
        if bridgeValid {
                for key, value := range bridgeVars {
                        s := strings.Split(key, ":")
                        if s[0] == "s" {
                                var tok Token
                                tok.n = n
                                n++
                                tok.contract = value.(string)
                                tok.bridgeable = true

                                fee_str, _ := d.DeroGetVar(value.(string), "bridgeFee")
                                fee, _ := strconv.Atoi(fee_str)
                                tok.bridgeFee = uint64(fee)

                                dec_str, _ := d.DeroGetVar(value.(string), "decimals")
                                tok.decimals, _ = strconv.Atoi(dec_str)

                                hex_str, _ := d.DeroGetVar(value.(string), "name")
                                bytes,_ := hex.DecodeString(hex_str)
                                tok.name = string(bytes)

                                hex_str, _ = d.DeroGetVar(value.(string), "native_symbol")
                                bytes,_ = hex.DecodeString(hex_str)
                                tok.native_symbol = string(bytes)

                                hex_str,_ = d.DeroGetVar(bridgeRegistry, "d:" + value.(string))
                                bytes,_ = hex.DecodeString(hex_str)
                                tok.eth_contract = string(bytes)

                                tokens[s[1]] = tok
                        }
                }
        }

        // swappable tokens
        swapVars, swapValid := d.DeroGetVars(swapRegistry)
        if (swapValid) {
                for key, value := range swapVars {
                        s := strings.Split(key, ":")
                        if s[0] == "t" && s[2] == "c" {
                                var tok Token = tokens[s[1]]

                                if tok == (Token{}) {
                                        tok.n = n
                                        n++
                                        tok.contract = value.(string)

                                        dec_str, _ := d.DeroGetVar(swapRegistry, "t:" + s[1] + ":d")
                                        tok.decimals, _ = strconv.Atoi(dec_str)
                                }

                                tok.swapable = true
                                tokens[s[1]] = tok
                        }
                }
        }

        // build list
        tokenList = make([]string, len(tokens))
        for k, v := range tokens {
                tokenList[v.n] = k
        }
}

func getPairs() {
        pairs = make(map[string]Pair)
        tokenGraph = graph.New(len(tokens))
        swapVars, swapValid := d.DeroGetVars(swapRegistry)

        if (swapValid) {
                for key, value := range swapVars {
                        s := strings.Split(key, ":")
                        if s[0] == "p" {
                                var pair Pair

                                pair.contract = value.(string)

                                fee_str, _ := d.DeroGetVar(pair.contract, "fee")
                                fee, _ := strconv.Atoi(fee_str)
                                pair.fee = uint64(fee)

                                val1_str, _ := d.DeroGetVar(pair.contract, "val1")
                                val1, _ := strconv.Atoi(val1_str)
                                pair.val1 = uint64(val1)

                                val2_str, _ := d.DeroGetVar(pair.contract, "val2")
                                val2, _ := strconv.Atoi(val2_str)
                                pair.val2 = uint64(val2)

                                adds_str, _ := d.DeroGetVar(pair.contract, "adds")
                                adds, _ := strconv.Atoi(adds_str)
                                pair.adds = uint64(adds)

                                rems_str, _ := d.DeroGetVar(pair.contract, "rems")
                                rems, _ := strconv.Atoi(rems_str)
                                pair.rems = uint64(rems)

                                swaps_str, _ := d.DeroGetVar(pair.contract, "swaps")
                                swaps, _ := strconv.Atoi(swaps_str)
                                pair.swaps = uint64(swaps)

                                shares_str, _ := d.DeroGetVar(pair.contract, "sharesOutstanding")
                                shares, _ := strconv.Atoi(shares_str)
                                pair.sharesOutstanding = uint64(shares)

                                pairs[s[1] + ":" + s[2]] = pair

                                if pair.val1 > 0 {
                                        tok1 := tokens[s[1]]
                                        tok2 := tokens[s[2]]

                                        val1_float := float64(pair.val1) / math.Pow(10, float64(tok1.decimals))
                                        val2_float := float64(pair.val2) / math.Pow(10, float64(tok2.decimals))

                                        tokenGraph.AddCost(tok1.n, tok2.n, int64(val2_float / val1_float * math.Pow(10, 7)))
                                        tokenGraph.AddCost(tok2.n, tok1.n, int64(val1_float / val2_float * math.Pow(10, 7)))
                                }
                        }
                }
        }
}

func conversion(sym1 string, sym2 string) (ratio float64, path string) {
        if tokens[sym1] == (Token{}) || tokens[sym2] == (Token{}) {
                return
        }

        n1 := tokens[sym1].n
        n2 := tokens[sym2].n

        p, d := graph.ShortestPath(tokenGraph, n1, n2)
        if d == -1 {
                return
        }

        ratio = float64(1.0)

        n := n1
        path = sym1

        for i := 1; i < len(p); i++ {
                ratio *= (float64(tokenGraph.Cost(n, p[i])) / math.Pow(10, 7))
                path += " => " + tokenList[p[i]]
                n = p[i]
        }

        return
}

func Tokens() (reply string) {
	getTokens()

	reply += "```"
	reply += fmt.Sprintf("%-10s %-64s %-7s %-7s\n\n", "TOKEN", "CONTRACT", "SWAP", "BRIDGE")
	for key, tok := range tokens {
		swap_check := "\u2716"
		bridge_check := "\u2716"

		if (tok.swapable) {
			swap_check = "\u2714"
		}
		if (tok.bridgeable) {
			bridge_check = "\u2714"
		}

		reply += fmt.Sprintf("%-10s %64s    %s       %s\n", key, tok.contract, swap_check, bridge_check)
	}
	reply += "```"

	return
}

func Quote(words []string) (reply string) {
        if len(words) != 2 {
                reply = "quote requires 2 arguments"
                return
        }

        getPairs()

        ratio, path := conversion(words[0], words[1])
        if len(path) == 0 {
                reply = "Cannot find path between '" + words[0] + "' and '" + words[1] + "'\n"
                return
        }

        reply += path + "\n"
        reply += fmt.Sprintf("1 %s == %0.7f %s\n", words[0], ratio, words[1])

	return
}

func QuoteDero() (quote float64) {
	getPairs()
        quote, _ = conversion("DERO", "DUSDT")

	return
}

func Pairs() (reply string) {
        tlv := float64(0)
        getPairs()

	reply += "```"
        reply += fmt.Sprintf("%-15s %30s\n\n", "PAIR", "TOTAL LIQUIDITY")
        for key, pair := range pairs {
                if pair.sharesOutstanding > 0 {
                        s := strings.Split(key, ":")
                        tokenA := tokens[s[0]]
                        tokenB := tokens[s[1]]

                        val1 := d.DeroFormatMoneyPrecision(pair.val1, tokenA.decimals)
                        val2 := d.DeroFormatMoneyPrecision(pair.val2, tokenB.decimals)

                        ratio1, _ := conversion(s[0], "DUSDT")
                        ratio2, _ := conversion(s[1], "DUSDT")

                        val1_float, _ := val1.Float64()
                        val2_float, _ := val2.Float64()

                        tlv += val1_float * ratio1
                        tlv += val2_float * ratio2

                        reply += fmt.Sprintf("%-15s %18.7f/%18.7f, %10s = %0.7f %s\n", key, val1, val2, s[0], val2_float / val1_float, s[1])
                } else {
                        reply += fmt.Sprintf("%-15s %18.7f/%18.7f\n", key, 0.0, 0.0)
                }
        }

        reply += fmt.Sprintf("\n")
        reply += fmt.Sprintf("TLV: %.2f USDT\n", tlv)
	reply += "```"

	return
}

func DexInit() {
	d.DeroInit(config.Daemon)

	bridgeRegistry, _ = d.DeroGetKeyHex("dex.bridge.registry")
	swapRegistry, _ = d.DeroGetKeyHex("dex.swap.registry")

	getTokens()
	getPairs()
}
