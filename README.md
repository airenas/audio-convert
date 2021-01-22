# audio-convert

![Go](https://github.com/airenas/audio-convert/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/airenas/audio-convert/badge.svg?branch=main)](https://coveralls.io/github/airenas/audio-convert?branch=main) ![CodeQL](https://github.com/airenas/audio-convert/workflows/CodeQL/badge.svg)

Service to convert wav audio to mp3, m4a formats. The service wraps *ffmpeg* to do the conversion.

To test the service look into [examples/docker-compose](examples/docker-compose).  

```bash
    cd examples/docker-compose
    docker-compose up -d
    curl -X POST http://localhost:8003/convert -H 'content-type: application/json' -d @data.json
```

Input is a base64 encoded wav file put into json. See sample [data.json](examples/docker-compose/data.json).

The result is a json with base64 encoded data (mp3 or m4a file):

```json
{"audio":"AAAAH..."}
```

---

## License

Copyright © 2021, [Airenas Vaičiūnas](https://github.com/airenas).

Released under the [The 3-Clause BSD License](LICENSE).

---