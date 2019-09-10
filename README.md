## XD Codec

It is just a [binary](https://golang.org/pkg/encoding/binary/)-like package, inspired by [TLV](https://en.wikipedia.org/wiki/Type-length-value), [amqp.table](https://www.amqp.org/), [Protocol Buffers](https://developers.google.com/protocol-buffers/), etc.


The Encoding Format:

1. All size fixed numbers are encoded as big/little endian.

2. varints(unsigned included): https://developers.google.com/protocol-buffers/docs/encoding#varints

3. Sized bytes/string:

    | size    | bytes/string |
    |---------|--------------|
    | uvarint | raw bytes    |

5. Containers:

    `List` and `Map` are encoded as below:

    | number of elements(k/v pairs) | elements(k/v pairs)      |
    |-------------------------------|--------------------------|
    | uint8 (limit to 255)          | []Typed/map[string]Typed |

    K-V pair:

    | key          | value |
    |--------------|-------|
    | sized string | Typed |

    The `Typed` just a basic type with a prefix byte to indicate the type:

    | name        | prefix | value format      |
    |-------------|--------|-------------------|
    | TypedInt    | 'i'    | varint            |
    | TypedUint   | 'u'    | uvarint           |
    | TypedFloat  | 'f'    | uvarint(IEEE 754) |
    | TypedBytes  | 'b'    | sized bytes       |
    | TypedString | 's'    | sized string      |
    | TypedList   | 'l'    | []Typed           |
    | TypedMap    | 'm'    | map[string]Typed  |

    **NOTE**: The **limitation** of varints is that, the larger numbers may require more than 8 bytes, but I think it is big enough..
