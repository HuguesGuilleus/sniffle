# Fetch

Each request is saved in a separaated file.

## File path

The path for `scheme://host/...` in base is `base/{scheme}/{host}/{ID}.http`

The ID is the SHA256 in hexadecimal of:

```txt
METHOD\n
URL\n
HEADER1: VALUE1.1\n
HEADER1: VALUE1.2\n
HEADER2: VALUE2.1\n
\n
BODY...\n
```

> [!NOTE]
> A newline at end. The newline character is `\n` or `0x0A`. The header name is
> `Aaaa-Bbbb-Cccc`.

Example

```txt
GET
https://example.com/file?a=1&b=2
Accept-Encoding: application/json

Hello World!
```

Result: `150ee78242120c0e38fc747a175c56068c2f07f8b0c57345a7ee6cdd5a172d05`

## Format

```txt
"HTTP" json_len:big_endian_uint32 json:[json_len]byte responseBody:[...]byte
```

The json contain meta info of the request in json:

NOTE: At end of JSON, add two `\n\n` only for display file. This two byte are
included in the json length.

```json
{
	"time": "yyyy-mm-ddThh:mm:ss.nnnnnnnnnZ",
	"requestMethod": "GET",
	"requestURL": "https://example.net/dir/file.txt?a=1",
	"requestHeader": {
		"header1": ["value1", "value2"]
	},
	"requestBody": "<base64 of the body>",
	"status": 200,
	"responseHeader": {
		"header1": ["value1", "value2"]
	}
}

{
	"time": "2024-11-12T22:44:59.862240355Z",
	"requestMethod": "GET",
	"requestURL": "https://example.net/dir/file.txt?a=1",
	"requestHeader": {
		"K1": ["v1"],
		"K2": ["v2"]
	},
	"requestBody": "Ym9keQ==",
	"status": 200,
	"responseHeader": {
		"H1": ["v1","v2"]
	}
}
```
