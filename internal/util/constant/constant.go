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

var (
	//Patterns is using to store valid route patterns
	Patterns = [4]string{
		`/i/(?P<slug>gdrive|gridfs|abs|s3|url)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gtc]{1,3})/(?P<path>.*)`,
		`/i/(?P<slug>gdrive|gridfs|abs|s3|url)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		`/i/(?P<slug>gdrive|gridfs|abs|s3|url)/(?P<webp>webp)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gtc]{1,3})/(?P<path>.*)`,
		`/i/(?P<slug>gdrive|gridfs|abs|s3|url)/(?P<webp>webp)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
	}
)
