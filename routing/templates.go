package routing



const indexHtml = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>Webservice</title>
        </head>
        <body>
            <h1 style="color: {{ .Color }}">Hello World, again!</h1>
            <p>Version: {{ .Version }}</p>
        </body>
    </html>
`

type indexHtmlData struct {
    Version string
    Color   string
}


const metricsText = `
    # HELP state_entries_quantity The current number of state entries being stored
    # TYPE state_entries_quantity gauge
    state_entries_quantity {{ .Count }}
`

type metricsTextData struct {
    Count int
}
