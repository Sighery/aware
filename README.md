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
of the pair that you spent.

```
{
	"symbol": "ETHEUR",
	"first_amount": 0.451,
	"second_amount": 1000.28192
}
```

This trade would mean you bought 0.451 ETH, with 1000.28192 EUR.

[example.aware.toml]: ./example.aware.toml
[example.aware.json]: ./example.aware.json
