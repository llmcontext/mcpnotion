package tools

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
)

func getPageContent(ctx context.Context, client *notionapi.Client, pageId string) ([]string, error) {
	content := []string{}
	// let's iterate over the page content to find the databases
	startCursor := notionapi.Cursor("")
	hasMoreBlocks := true

	for hasMoreBlocks {
		children, err := client.Block.GetChildren(
			ctx,
			notionapi.BlockID(notionapi.PageID(pageId)),
			&notionapi.Pagination{
				StartCursor: startCursor,
				PageSize:    100, // Maximum allowed page size
			},
		)
		if err != nil {
			return content, err
		}
		for _, child := range children.Results {
			var blockType = child.GetType().String()

			switch blockType {
			case notionapi.BlockTypeParagraph.String():
				paragraphBlock := child.(*notionapi.ParagraphBlock)
				content = append(content, getTextFromRichText(paragraphBlock.Paragraph.RichText))

			case notionapi.BlockTypeHeading1.String():
				heading1Block := child.(*notionapi.Heading1Block)
				content = append(content, fmt.Sprintf("# %s", getTextFromRichText(heading1Block.Heading1.RichText)))
				content = append(content, "")

			case notionapi.BlockTypeHeading2.String():
				heading2Block := child.(*notionapi.Heading2Block)
				content = append(content, fmt.Sprintf("## %s", getTextFromRichText(heading2Block.Heading2.RichText)))
				content = append(content, "")

			case notionapi.BlockTypeHeading3.String():
				heading3Block := child.(*notionapi.Heading3Block)
				content = append(content, fmt.Sprintf("### %s", getTextFromRichText(heading3Block.Heading3.RichText)))
				content = append(content, "")

			case notionapi.BlockTypeBulletedListItem.String():
				bulletedListItemBlock := child.(*notionapi.BulletedListItemBlock)
				content = append(content, fmt.Sprintf("- %s", getTextFromRichText(bulletedListItemBlock.BulletedListItem.RichText)))

			case notionapi.BlockTypeNumberedListItem.String():
				numberedListItemBlock := child.(*notionapi.NumberedListItemBlock)
				content = append(content, fmt.Sprintf("1. %s", getTextFromRichText(numberedListItemBlock.NumberedListItem.RichText)))

			}

			// if there are more blocks, we need to get the next cursor
			if children.HasMore {
				startCursor = notionapi.Cursor(children.NextCursor)
			} else {
				hasMoreBlocks = false
			}
		}
	}
	return content, nil
}

func getTextFromRichText(richText []notionapi.RichText) string {
	text := ""
	for _, richText := range richText {
		text += richText.PlainText
	}
	return text
}
