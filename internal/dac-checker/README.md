# dac-checker

`dac-checker` is a lightweight Swift command-line utility that verifies if a specified USB DAC (Digital-to-Analog Converter) is connected and operating at a target sample rate. It is designed for use on macOS and intended to be integrated with automation tools or system watchers.

## Features

- Detects USB audio devices by name
- Verifies current sample rate
- Returns clear status codes for scripting integration

## Usage

```sh
./dac-checker "<DAC Name>" <Sample Rate>
```

### Example

```sh
./dac-checker "FiiO K5 Pro" 48000
```

### Output

- `READY`: DAC found and sample rate matches
- `WRONG_RATE`: DAC found but sample rate does not match
- `NOT_FOUND`: DAC not found

### Exit Codes

- `0`: DAC is ready (found and correct sample rate)
- `1`: DAC not found or incorrect sample rate

## Building

From the root of the repository:

```sh
cd internal/dac-checker
swift build -c release
```

The compiled binary will be located at:

```
.build/release/dac-checker
```

You can then copy it to `/usr/local/bin` or another location in your `PATH`.

## License

MIT License
