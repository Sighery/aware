# aware

Because I'm too lazy to be constantly opening up Binance and looking at
graphs, I decided to build this tool to look at my trades, and then compare
against the current price and tell me the percentage difference.

## Usage

```sh
aware myconfig.aware.toml
```

You should provide a TOML configuration file. You can take a look at the
[example.aware.toml][] to figure out the current schema.

In the `trades_file` configuration option, you should specify a JSON file with
all your relevant investments.

You can take a look at the [example.aware.json] to figure out the schema
there. You should provide the `symbol` (they can be repeated), as well as the
`first_amount` and `second_amount`. The `first_amount` is the amount you
bought of the first part of the pair. While `second_amount` is the second part
of the pair that you spent. `remaining` just keeps track of how much of this
trade is remaining (in case you sold 0.2 out of 0.5 BTC, you'd still want the
remaining 0.3 BTC to be calculated against the original price you bought at).

This trade would mean you bought 0.451 ETH, with 1000.28192 EUR, and that you
didn't sell any of it:

```
{
	"symbol": "ETHEUR",
	"first_amount": 0.451,
	"remaining": 0.451,
	"second_amount": 1000.28192
}
```

This trade would mean you bought 0.451 ETH, with 1000.28192 EUR, and that you
already sold 0.251 ETH, meaning you still have 0.2 ETH. If you
`second_amount / first_amount`, you can get the price you bought at (the price
where 1 ETH equals X EUR). For this trade, this was 2217.92 EUR per 1 ETH.
Meaning you still have 0.2 ETH at the 2217.92 EUR price, so any comparisons
should be calculated against that original price:

```
{
	"symbol": "ETHEUR",
	"first_amount": 0.451,
	"remaining": 0.2,
	"second_amount": 1000.28192
}
```


[example.aware.toml]: ./example.aware.toml
[example.aware.json]: ./example.aware.json
