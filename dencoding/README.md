# dencoding - Dasel Encoding

This package provides encoding implementations for all supported file formats.

The main difference is that it aims to keep maps ordered, where the default encoders/decoders do not.

## Support formats

### Decoding

Custom decoders are required to ensure that map/object values are decoded into the `Map` type rather than a standard `map[string]any`.

### Encoding

The `Map` type must have the appropriate Marshal func on it to ensure marshalling it in the desired format retains the ordering.
