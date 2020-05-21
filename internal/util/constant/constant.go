package constant

var (
	//Hostname is using to serve the application
	Hostname string
	//Port is using to serve the application
	Port string
	//CacheControlMaxAge is using to set max-age property of cache-control header
	CacheControlMaxAge int
)

var (
	//ErrorMessage is using to write the HTML response when any exception is occurred
	ErrorMessage = `
		<h1>Oops! Something went wrong...</h1>
		<p>We seem to be having some technical difficulties. Hang tight.</p>
	`
	//LogErrorFormat is using to format the message
	LogErrorFormat = "%s: %s"
	//LogErrorMessage is using to write a log message
	LogErrorMessage = "Unable to write HTTP response message"
)
