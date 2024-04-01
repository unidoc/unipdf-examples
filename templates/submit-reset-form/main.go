package main

import (
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/common"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io.
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	runTemplate()
}

func runTemplate() {
	c := creator.New()
	var (
		tpl = `
		<division>
			<form>
				<text-field label="Full Name" name="full name" rect="123.97, 619.02, 343.99, 633.6">
					<paragraph x="43" y="162.98">
						<text-chunk>Full Name</text-chunk>
					</paragraph>
					<line thickness="1" x1="123.97" y1="172.98" x2="343" y2="162.98"></line>
				</text-field>
				<submit-button x="400" y="400" height="20" width="50" label="Submit" url="https://unidoc.io" fill-color="#00FF00" label-color="#FF0000"></submit-button>
				<reset-button x="100" y="400" height="20" width="50" label="Reset" fill-color="#808080" label-color="#FFFFFF"></reset-button>
			</form>
		</division>
		`
	)

	// Draw template.
	err := c.DrawTemplate(strings.NewReader(tpl), nil, nil)
	if err != nil {
		panic(err)
	}

	if err := c.WriteToFile("form-template_example.pdf"); err != nil {
		panic(err)
	}
}
