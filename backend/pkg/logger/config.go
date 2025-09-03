package logger

type Config struct {
}

const mapping = `
{
    "mappings": {
        "properties": {
            "name": {
                "type": "text"
            },
            "log_level": {
                "type": "text"
            },
            "debug": {
                "type": "text"
            },
            "content": {
                "type": "text",
				"analyzer": "standard"
            },
            "time": {
                "type": "date"
            }
        }
    }
}
`
