package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

const Version = "v2021-05-11"

func setupLogging(arguments docopt.Opts) {
	verbose, _ := arguments.Bool("--verbose")
	logLevel := log.InfoLevel
	if verbose {
		logLevel = log.DebugLevel
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stderr)
	log.SetLevel(logLevel)
}

func doDownload(arguments docopt.Opts) int {
	// Download supports 2 different modes.
	// When the <filepath> is absent, the filepath is inferred from the s3path.
	//
	// download <s3path>
	// download <s3path> <filepath>
	s3path, _ := arguments.String("<s3path>")
	s3pathparts := strings.Split(s3path, "/")
	bucket := s3pathparts[0]
	key := path.Join(s3pathparts[1:]...)

	var filepath string
	if val, ok := arguments["<filepath>"].(string); ok {
		filepath = val
	} else {
		filepath = s3pathparts[len(s3pathparts)-1]
	}

	// For debugging, ensure the command you are running is what you intend to.
	dryrun, _ := arguments.Bool("--dry-run")
	if dryrun {
		fmt.Printf("download %s/%s %s\n", bucket, key, filepath)
		return 0
	}

	// Connect to S3 and perform the download operation.
	downloadObject(bucket, key, filepath)
	return 0
}

func doUpload(arguments docopt.Opts) int {
	// Upload supports 2 different modes.
	// When the <s3path> is absent, the filepath is assumed to be data-bysha256
	// and the file will be keyed by its sha256 hexdigest.
	//
	// upload <filepath>
	// upload <filepath> <s3path>
	filepath, _ := arguments.String("<filepath>")
	hexdigest, _ := computeSha256(filepath)

	var s3path string
	if val, ok := arguments["<s3path>"].(string); ok {
		s3path = val
	} else {
		s3path = fmt.Sprintf("data-bysha256/%s", hexdigest)
	}
	s3pathparts := strings.Split(s3path, "/")
	bucket := s3pathparts[0]
	key := path.Join(s3pathparts[1:]...)
	if strings.HasSuffix(s3path, "/") {
		// If the s3path ends with a slash, append the basename of the target
		// uploaded file.
		key = path.Join(key, path.Base(filepath))
	}

	// For debugging, ensure the command you are running is what you intend to.
	dryrun, _ := arguments.Bool("--dry-run")
	if dryrun {
		fmt.Printf("upload %s %s/%s\n", filepath, bucket, key)
		return 0
	}

	// Connect to S3 and perform the upload operation.
	uploadObject(bucket, key, filepath)
	return 0
}

func main() {
	usage := `S3 Client.

Populate AWS credentials:

AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

Other important environment variables:

AWS_ENDPOINT_URL
AWS_REGION
INSECURE_SKIP_VERIFY

When uploading without an explicit s3path, the client will perform a sha256
checksum of the input file and upload that to the data-bysha256 bucket.

Usage:
  s3cli download [options] <s3path>
  s3cli download [options] <s3path> <filepath>
  s3cli upload [options] <filepath>
  s3cli upload [options] <filepath> <s3path>

Options:
  -h --help     Show this screen.
  -n --dry-run  Display what action will be performed.
  -v --verbose  Enable verbose logging.
  --version     Show version.
`

	arguments, _ := docopt.ParseArgs(usage, nil, Version)

	// Parse the incoming arguments.
	download, _ := arguments.Bool("download")
	upload, _ := arguments.Bool("upload")
	setupLogging(arguments)

	// Perform S3 Action.
	if download {
		os.Exit(doDownload(arguments))
	} else if upload {
		os.Exit(doUpload(arguments))
	}
}
