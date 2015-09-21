package operations

const (

	//permissions
	DEFAULT_FILE_PERMISSION = 0600
	UPLOAD_FILE_PERMISSION  = 0775

	//middleware tags
	MIDDLEWARE_CURRENT_SESSION      = "CurrentSession"
	MIDDLEWARE_IPADDRESS_HEADER     = "X-Forwarded-For"
	MIDDLEWARE_APPLICATION_HOST     = "ApplicationHost"
	MIDDLEWARE_STATIC_PATH          = "StaticPath"
	MIDDLEWARE_BASE_API             = "baseApi"
	MIDDLEWARE_BASE_CONTROLLER_NAME = "base"
	MIDDLEWARE_BASE_URL             = "BaseUrl"
	MIDDLEWARE_DB                   = "Db"
	MIDDLEWARE_TOKEN_LABEL          = "token"

	//mime types
	CONTENT_TYPE_HEADER = "Content-Type"
	CONTENT_DISPOSITION = "Content-Disposition"
	PDF_MIME_TYPE       = "application/pdf"
	HTML_MIME_TYPE      = "text/html"
	CSS_MIME_TYPE       = "text/css"
	HEADER_STATUS       = "status"

	//issues & projects
	ISSUE_SEQUENCE_NAME   = "seq_issue"
	ISSUE_NO_PROJECT_MASK = "%010d"

	//data formats
	DATE_FORMAT_JSON = "2006-01-02T15:04:05.000Z"
	YEAR_MIN_VALUE   = 2000
)
