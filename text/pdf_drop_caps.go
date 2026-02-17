/*
 * Drop caps example for StyledParagraph.
 * Demonstrates both traditional drop caps (DropCapsDrop) and inline drop caps (DropCapsInline).
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v4/common"
	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/creator"
	"github.com/unidoc/unipdf/v4/model"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		panic(err)
	}

	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

func main() {
	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Load fonts
	regularFont := model.NewStandard14FontMustCompile(model.HelveticaName)
	boldFont := model.NewStandard14FontMustCompile(model.HelveticaBoldName)

	// Title
	title := c.NewStyledParagraph()
	tc := title.Append("Drop Caps Examples")
	tc.Style.Font = boldFont
	tc.Style.FontSize = 24
	title.SetMargins(0, 0, 20, 0)
	if err := c.Draw(title); err != nil {
		log.Fatal(err)
	}

	textStyle := creator.TextStyle{
		Font:     regularFont,
		FontSize: 12,
	}

	// Example 1: Traditional drop cap with DropCapsDrop (first character)
	section1 := c.NewStyledParagraph()
	ch := section1.Append("1. Traditional Drop Caps (DropCapsDrop - First Character)")
	ch.Style.Font = boldFont
	ch.Style.FontSize = 14
	section1.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section1); err != nil {
		log.Fatal(err)
	}

	// Regular paragraph for comparison
	regularP := c.NewStyledParagraph()
	regularP.SetTextAlignment(creator.TextAlignmentLeft)
	regularP.SetLineHeight(1.2)
	regularP.SetMargins(0, 0, 10, 0)

	regularChunk := regularP.Append("This is a regular paragraph without any drop caps, shown for comparison. Notice how the text flows normally from left to right without any special formatting for the first character or word.")
	regularChunk.Style = textStyle

	if err := c.Draw(regularP); err != nil {
		log.Fatal(err)
	}

	// Another regular paragraph for comparison
	regularP2 := c.NewStyledParagraph()
	regularP2.SetTextAlignment(creator.TextAlignmentLeft)
	regularP2.SetLineHeight(1.2)
	regularP2.SetMargins(0, 0, 10, 0)

	regularChunk2 := regularP2.Append("Here is another regular paragraph to show the contrast before we apply the drop cap effect. The text begins at the left margin and flows naturally across the full width of the page.")
	regularChunk2.Style = textStyle

	if err := c.Draw(regularP2); err != nil {
		log.Fatal(err)
	}

	// Single chunk with FirstCharacter scope - "T" will be drop cap
	p1 := c.NewStyledParagraph()
	p1.SetTextAlignment(creator.TextAlignmentLeft)
	p1.SetLineHeight(1.2)
	p1.SetMargins(0, 0, 10, 0)

	chunk1 := p1.Append("This is a paragraph with a beautiful drop cap. The first three lines will flow beside the large T, with reduced width. After that, the text continues below the drop cap at full width, just like in traditional magazine layouts. This demonstrates the classic drop caps effect where only the first N lines are affected.")
	chunk1.Style = textStyle
	p1.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsDrop,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 3,
		Gap:      5.0,
	})

	if err := c.Draw(p1); err != nil {
		log.Fatal(err)
	}

	// Regular paragraph after first drop cap example
	regularP3 := c.NewStyledParagraph()
	regularP3.SetTextAlignment(creator.TextAlignmentLeft)
	regularP3.SetLineHeight(1.2)
	regularP3.SetMargins(0, 0, 10, 0)

	regularChunk3 := regularP3.Append("Notice how the drop cap paragraph above has the first three lines flowing beside the large letter T, and then continues at full width below it. This is the classic drop caps effect used in books and magazines.")
	regularChunk3.Style = textStyle

	if err := c.Draw(regularP3); err != nil {
		log.Fatal(err)
	}

	// Example 2: Traditional drop cap with first word
	section2 := c.NewStyledParagraph()
	ch2 := section2.Append("2. Traditional Drop Caps (DropCapsDrop - First Word)")
	ch2.Style.Font = boldFont
	ch2.Style.FontSize = 14
	section2.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section2); err != nil {
		log.Fatal(err)
	}

	p2 := c.NewStyledParagraph()
	p2.SetTextAlignment(creator.TextAlignmentLeft)
	p2.SetLineHeight(1.2)
	p2.SetMargins(0, 0, 10, 0)

	// Single chunk with FirstWord scope - "CHAPTER" will be drop cap
	chunk2 := p2.Append("CHAPTER ONE: It was the best of times, it was the worst of times. This example shows how an entire word can be used as a drop cap, which is useful for chapter markers or special formatting. The word CHAPTER appears enlarged and the first three lines flow beside it.")
	chunk2.Style = textStyle
	p2.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsDrop,
		Scope:    creator.DropCapsFirstWord,
		NumLines: 3,
		Gap:      8.0,
	})

	if err := c.Draw(p2); err != nil {
		log.Fatal(err)
	}

	// Example 3: Inline drop cap (all lines flow beside)
	section3 := c.NewStyledParagraph()
	ch3 := section3.Append("3. Inline Drop Caps (DropCapsInline - All Lines Beside)")
	ch3.Style.Font = boldFont
	ch3.Style.FontSize = 14
	section3.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section3); err != nil {
		log.Fatal(err)
	}

	p3 := c.NewStyledParagraph()
	p3.SetTextAlignment(creator.TextAlignmentLeft)
	p3.SetLineHeight(1.2)
	p3.SetMargins(0, 0, 10, 0)

	chunk3 := p3.Append("In this paragraph style, ALL lines maintain reduced width and flow beside the drop cap. No text appears below the enlarged letter. This is the key difference from traditional drop caps - the entire paragraph is affected, not just the first N lines. This style is less common but can be useful for shorter paragraphs or special design effects.")
	chunk3.Style = textStyle
	p3.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsInline,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 2,
		Gap:      10.0,
	})

	if err := c.Draw(p3); err != nil {
		log.Fatal(err)
	}

	// Example 4: Drop caps with mixed styling (multiple chunks)
	section4 := c.NewStyledParagraph()
	ch4 := section4.Append("4. Drop Caps with Mixed Text Styling")
	ch4.Style.Font = boldFont
	ch4.Style.FontSize = 14
	section4.SetMargins(0, 0, 10, 0)
	if err := c.Draw(section4); err != nil {
		log.Fatal(err)
	}

	p4 := c.NewStyledParagraph()
	p4.SetTextAlignment(creator.TextAlignmentLeft)
	p4.SetLineHeight(1.2)
	p4.SetMargins(0, 0, 10, 0)

	// First chunk with drop cap
	chunk4a := p4.Append("This paragraph demonstrates ")
	chunk4a.Style = textStyle

	// Second chunk with bold font
	chunk4b := p4.Append("mixed styling ")
	chunk4b.Style.Font = boldFont
	chunk4b.Style.FontSize = 12

	// Third chunk with larger font size
	chunk4c := p4.Append("with different fonts ")
	chunk4c.Style.Font = regularFont
	chunk4c.Style.FontSize = 24

	// Fourth chunk back to regular
	chunk4d := p4.Append("and font sizes within the same paragraph. The drop cap affects the first three lines, and the text flows naturally despite having multiple chunks with different styling properties. This shows that drop caps work seamlessly with complex text formatting.")
	chunk4d.Style = textStyle

	p4.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsDrop,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 3,
		Gap:      5.0,
	})

	if err := c.Draw(p4); err != nil {
		log.Fatal(err)
	}

	// Example 5: Single word as drop cap
	section5 := c.NewStyledParagraph()
	ch5 := section5.Append("5. Single Word as Drop Cap")
	ch5.Style.Font = boldFont
	ch5.Style.FontSize = 14
	section5.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section5); err != nil {
		log.Fatal(err)
	}

	p5 := c.NewStyledParagraph()
	p5.SetTextAlignment(creator.TextAlignmentLeft)
	p5.SetLineHeight(1.2)
	p5.SetMargins(0, 0, 10, 0)

	chunk5 := p5.Append("WORD")
	chunk5.Style = textStyle
	p5.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsDrop,
		Scope:    creator.DropCapsFirstWord,
		NumLines: 2,
		Gap:      5.0,
	})

	if err := c.Draw(p5); err != nil {
		log.Fatal(err)
	}

	pb := c.NewPageBreak()
	if err := c.Draw(pb); err != nil {
		log.Fatal(err)
	}

	// Example 6: Drop caps with custom color
	section6 := c.NewStyledParagraph()
	ch6 := section6.Append("6. Drop Caps with Custom Color (Red)")
	ch6.Style.Font = boldFont
	ch6.Style.FontSize = 14
	section6.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section6); err != nil {
		log.Fatal(err)
	}

	p6 := c.NewStyledParagraph()
	p6.SetTextAlignment(creator.TextAlignmentLeft)
	p6.SetLineHeight(1.2)
	p6.SetMargins(0, 0, 10, 0)

	chunk6 := p6.Append("This is a paragraph with a red drop cap. The first character is styled with a custom color while the rest of the text remains in the default style. This demonstrates how you can make the drop cap stand out visually using the new Style field in DropCapsOptions.")
	chunk6.Style = textStyle

	redStyle := creator.TextStyle{
		Font:     boldFont,
		FontSize: 12,                                  // Will be overridden by NumLines calculation
		Color:    creator.ColorRGBFrom8bit(200, 0, 0), // Red
	}
	p6.SetDropCapsWithStyle(creator.DropCapsDrop, creator.DropCapsFirstCharacter, 3, 5.0, &redStyle)

	if err := c.Draw(p6); err != nil {
		log.Fatal(err)
	}

	// Example 7: Drop caps with different font and color (Blue)
	section7 := c.NewStyledParagraph()
	ch7 := section7.Append("7. Drop Caps with Bold Font and Blue Color")
	ch7.Style.Font = boldFont
	ch7.Style.FontSize = 14
	section7.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section7); err != nil {
		log.Fatal(err)
	}

	p7 := c.NewStyledParagraph()
	p7.SetTextAlignment(creator.TextAlignmentLeft)
	p7.SetLineHeight(1.2)
	p7.SetMargins(0, 0, 10, 0)

	chunk7 := p7.Append("Another example with a blue drop cap using a bold font. The drop cap will be rendered in bold blue while the rest of the paragraph uses the regular font. This gives you complete control over the visual appearance of the drop cap, including font family, weight, color, and other text properties.")
	chunk7.Style = textStyle

	blueStyle := creator.TextStyle{
		Font:     boldFont,
		FontSize: 12,                                  // Will be overridden
		Color:    creator.ColorRGBFrom8bit(0, 0, 200), // Blue
	}
	p7.SetDropCapsWithStyle(creator.DropCapsDrop, creator.DropCapsFirstCharacter, 3, 8.0, &blueStyle)

	if err := c.Draw(p7); err != nil {
		log.Fatal(err)
	}

	// Example 8: Drop caps first word with custom style (Green)
	section8 := c.NewStyledParagraph()
	ch8 := section8.Append("8. Drop Caps First Word with Custom Style (Green)")
	ch8.Style.Font = boldFont
	ch8.Style.FontSize = 14
	section8.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section8); err != nil {
		log.Fatal(err)
	}

	p8 := c.NewStyledParagraph()
	p8.SetTextAlignment(creator.TextAlignmentLeft)
	p8.SetLineHeight(1.2)
	p8.SetMargins(0, 0, 10, 0)

	chunk8 := p8.Append("CHAPTER ONE begins here with the entire first word styled as a drop cap in green color. This is useful for chapter openings and other decorative text elements where you want to emphasize the beginning with custom styling. The word CHAPTER appears enlarged, bold, and colored.")
	chunk8.Style = textStyle

	greenStyle := creator.TextStyle{
		Font:     boldFont,
		FontSize: 12,                                  // Will be overridden
		Color:    creator.ColorRGBFrom8bit(0, 150, 0), // Green
	}
	p8.SetDropCapsWithStyle(creator.DropCapsDrop, creator.DropCapsFirstWord, 4, 10.0, &greenStyle)

	if err := c.Draw(p8); err != nil {
		log.Fatal(err)
	}

	// Example 9: Inline drop caps with custom style (Purple)
	section9 := c.NewStyledParagraph()
	ch9 := section9.Append("9. Inline Drop Caps with Custom Style (Purple)")
	ch9.Style.Font = boldFont
	ch9.Style.FontSize = 14
	section9.SetMargins(0, 0, 20, 0)
	if err := c.Draw(section9); err != nil {
		log.Fatal(err)
	}

	p9 := c.NewStyledParagraph()
	p9.SetTextAlignment(creator.TextAlignmentLeft)
	p9.SetLineHeight(1.2)
	p9.SetMargins(0, 0, 10, 0)

	chunk9 := p9.Append("Inline drop caps work differently - all lines flow beside the enlarged character. This example uses a purple drop cap with a bold font to create a distinctive look that maintains the paragraph flow. The Style field allows you to customize every aspect of the drop cap appearance.")
	chunk9.Style = textStyle

	purpleStyle := creator.TextStyle{
		Font:     boldFont,
		FontSize: 12,                                    // Will be overridden
		Color:    creator.ColorRGBFrom8bit(128, 0, 128), // Purple
	}
	p9.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsInline,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 2,
		Gap:      5.0,
		Style:    &purpleStyle,
	})

	if err := c.Draw(p9); err != nil {
		log.Fatal(err)
	}

	pb = c.NewPageBreak()
	if err := c.Draw(pb); err != nil {
		log.Fatal(err)
	}

	// Example 10: Comparison - traditional vs inline
	section10 := c.NewStyledParagraph()
	section10.SetMargins(0, 0, 20, 0)

	ch10 := section10.Append("10. Side-by-Side Comparison")
	ch10.Style.Font = boldFont
	ch10.Style.FontSize = 14
	if err := c.Draw(section10); err != nil {
		log.Fatal(err)
	}

	compTable := c.NewTable(2)
	compTable.SetMargins(0, 0, 10, 0)

	// Left cell for traditional drop cap
	cell1 := compTable.NewCell()

	// Right cell for inline drop cap
	cell2 := compTable.NewCell()

	// Traditional
	compareText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."

	p10a := c.NewStyledParagraph()
	p10a.SetTextAlignment(creator.TextAlignmentLeft)
	p10a.SetLineHeight(1.2)
	p10a.SetMargins(0, 0, 10, 0)

	chunk10a := p10a.Append(compareText)
	chunk10a.Style = textStyle
	p10a.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsDrop,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 2,
		Gap:      5.0,
	})

	if err := cell1.SetContent(p10a); err != nil {
		log.Fatal(err)
	}

	// Inline
	p10b := c.NewStyledParagraph()
	p10b.SetTextAlignment(creator.TextAlignmentLeft)
	p10b.SetLineHeight(1.2)
	p10b.SetMargins(0, 0, 10, 0)

	chunk10b := p10b.Append(compareText)
	chunk10b.Style = textStyle
	p10b.SetDropCaps(creator.DropCapsOptions{
		Type:     creator.DropCapsInline,
		Scope:    creator.DropCapsFirstCharacter,
		NumLines: 2,
		Gap:      5.0,
	})

	if err := cell2.SetContent(p10b); err != nil {
		log.Fatal(err)
	}

	if err := c.Draw(compTable); err != nil {
		log.Fatal(err)
	}

	// Save to file
	err := c.WriteToFile("pdf_drop_caps.pdf")
	if err != nil {
		log.Fatalf("Error writing PDF: %v", err)
	}

	fmt.Println("PDF created successfully: pdf_drop_caps.pdf")
}
