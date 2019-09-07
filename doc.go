// Package xdcodec provides binary codec library for easier use,
// maybe with better performance(TODO: benchmark).
/*
The Encoding Format

1. All size fixed numbers are encoded as big/little endian.

2. varints(unsigned included): https://developers.google.com/protocol-buffers/docs/encoding#varints

3. Sized bytes/string:

    | size    | bytes/string |
    |---------|--------------|
    | uvarint | raw bytes    |

5. Containers:

List and Map are encoded as below:

    | number of elements(k/v pairs) | elements(k/v pairs)      |
    |-------------------------------|--------------------------|
    | uint8 (limit to 255)          | []Typed/map[string]Typed |

K-V pair:

    | key          | value |
    |--------------|-------|
    | sized string | Typed |

The Typed just a basic value with a prefix byte to indicate the type:

    | name        | prefix | value format      |
    |-------------|--------|-------------------|
    | TypedInt    | 'i'    | varint            |
    | TypedUint   | 'u'    | uvarint           |
    | TypedFloat  | 'f'    | uvarint(IEEE 754) |
    | TypedBytes  | 'b'    | sized bytes       |
    | TypedString | 's'    | sized string      |
    | TypedList   | 'l'    | []Typed           |
    | TypedMap    | 'm'    | map[string]Typed  |
*/
package xdcodec
