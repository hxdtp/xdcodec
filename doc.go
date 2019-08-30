// Package bincodec provides binary codec library for easier use,
// maybe with better performance(TODO: benchmark).
/*
The Encoding Format

1. All size fixed numbers are encoded as big endian.

2. varints(unsigned included): https://developers.google.com/protocol-buffers/docs/encoding#varints

3. Sized bytes/string:
    | size    | bytes/string |
    | uvarint | raw bytes    |

5. Containers:

List and Map are encoded as below:

    | number of elements(k/v pair) | elements(k/v paris)      |
    | uint8 (limit to 255)         | []Typed/map[string]Typed |

The Typed is only for internal use, it just a basic type with a byte prefix:

    | name        | type              | prefix |
    | TypedInt    | varint            | 'i'    |
    | TypedUint   | uvarint           | 'u'    |
    | TypedFloat  | uvarint(IEEE 754) | 'f'    |
    | TypedBytes  | sized bytes       | 'b'    |
    | TypedString | sized string      | 's'    |
    | TypedList   | []Typed           | 'l'    |
    | TypedMap    | map[strng]Typed   | 'm'    |
*/
package xdcodec
