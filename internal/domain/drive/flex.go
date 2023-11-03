package drive

import (
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type FolderType int

const (
	PersonalFolder FolderType = iota
	SharedFolder
	OtherFolder
)

type FolderCarousel struct {
	CarouselContainer *linebot.CarouselContainer
}
type NewFolderCarouselParam struct {
	BubbleParams []NewFolderBubbleParam
}

func NewFolderCarousel(carouselParams NewFolderCarouselParam) FolderCarousel {
	var bubbles []*linebot.BubbleContainer
	for _, v := range carouselParams.BubbleParams {
		bubble := newFolderBubble(v)
		bubbles = append(bubbles, bubble)
	}

	return FolderCarousel{
		CarouselContainer: &linebot.CarouselContainer{
			Type: linebot.FlexContainerTypeCarousel,
			// Insert bubbles
			Contents: bubbles,
		},
	}
}

type NewFolderBubbleParam struct {
	Type          string
	Name          string
	Path          string
	ID            string
	InsideFolderM map[string]string
	FileM         map[string]string
}

func newFolderBubble(param NewFolderBubbleParam) *linebot.BubbleContainer {
	var allFlexComponents []linebot.FlexComponent

	title := folderTitleFlexComponents(param.Type, param.Name, param.Path)
	allFlexComponents = append(allFlexComponents, title...)

	folder := folderFlexComponents(param.InsideFolderM)
	allFlexComponents = append(allFlexComponents, folder...)

	file := fileFlexComponents(param.FileM)
	allFlexComponents = append(allFlexComponents, file...)

	button := buttonFlexComponents(param.ID)
	allFlexComponents = append(allFlexComponents, button...)

	// Single Bubble
	bubble := linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			// Insert Components
			Contents: allFlexComponents,
		},
		Styles: &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{
				Separator: true,
			},
		},
	}

	return &bubble
}

func folderTitleFlexComponents(folderType string, name string, path string) []linebot.FlexComponent {
	var title []linebot.FlexComponent
	title = append(title,
		&linebot.TextComponent{
			Type:   linebot.FlexComponentTypeText,
			Text:   folderType,
			Weight: linebot.FlexTextWeightTypeBold,
			Color:  "#1DB446",
			Size:   linebot.FlexTextSizeTypeSm,
		},
		&linebot.TextComponent{
			Type:   linebot.FlexComponentTypeText,
			Text:   name,
			Weight: linebot.FlexTextWeightTypeBold,
			Size:   linebot.FlexTextSizeTypeXxl,
			Margin: linebot.FlexComponentMarginTypeMd,
		},
		&linebot.TextComponent{
			Type:  linebot.FlexComponentTypeText,
			Text:  path,
			Size:  linebot.FlexTextSizeTypeXs,
			Color: "#aaaaaa",
			Wrap:  true,
		},
		&linebot.SeparatorComponent{
			Type:   linebot.FlexComponentTypeSeparator,
			Margin: linebot.FlexComponentMarginTypeXxl,
		})
	return title
}

func folderFlexComponents(folderM map[string]string) []linebot.FlexComponent {
	var folders []linebot.FlexComponent
	var folderFlex []linebot.FlexComponent
	// map為空
	if len(folderM) == 0 {
		folderFlex = append(folderFlex,
			&linebot.BoxComponent{
				Type:    linebot.FlexComponentTypeBox,
				Layout:  linebot.FlexBoxLayoutTypeVertical,
				Margin:  linebot.FlexComponentMarginTypeXxl,
				Spacing: linebot.FlexComponentSpacingTypeSm,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:       linebot.FlexComponentTypeText,
						Text:       "無",
						Size:       linebot.FlexTextSizeTypeSm,
						Color:      "#555555",
						Decoration: linebot.FlexTextDecorationTypeNone,
						MaxLines:   linebot.IntPtr(25),
						Align:      linebot.FlexComponentAlignTypeStart,
						Margin:     linebot.FlexComponentMarginTypeNone,
						Gravity:    linebot.FlexComponentGravityTypeCenter,
						Flex:       linebot.IntPtr(0),
					},
				},
			},
			// Separator
			&linebot.SeparatorComponent{
				Margin: linebot.FlexComponentMarginTypeXxl,
			},
		)
		return folderFlex
	}
	for folderID, folderName := range folderM {
		folders = append(folders, &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:       linebot.FlexComponentTypeText,
					Text:       folderName,
					Size:       linebot.FlexTextSizeTypeSm,
					Color:      "#555555",
					Decoration: linebot.FlexTextDecorationTypeUnderline,
					MaxLines:   linebot.IntPtr(25),
					Align:      linebot.FlexComponentAlignTypeStart,
					Margin:     linebot.FlexComponentMarginTypeNone,
					Gravity:    linebot.FlexComponentGravityTypeCenter,
					Flex:       linebot.IntPtr(0),
				},
				&linebot.FillerComponent{
					Type: linebot.FlexComponentTypeFiller,
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.PostbackAction{
						Label:       "進入資料夾",
						Data:        "action=openFolder&folderID=" + folderID,
						DisplayText: "進入" + folderName,
					},
					Style:      linebot.FlexButtonStyleTypeLink,
					Height:     linebot.FlexButtonHeightTypeSm,
					Gravity:    linebot.FlexComponentGravityTypeCenter,
					Flex:       linebot.IntPtr(0),
					AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
				},
			},
		})
	}
	folderFlex = append(folderFlex,
		&linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Margin:   linebot.FlexComponentMarginTypeXxl,
			Spacing:  linebot.FlexComponentSpacingTypeSm,
			Contents: folders,
		},
		// Separator
		&linebot.SeparatorComponent{
			Margin: linebot.FlexComponentMarginTypeXxl,
		},
	)

	return folderFlex

}

func fileFlexComponents(fileM map[string]string) []linebot.FlexComponent {
	totalFiles := len(fileM)

	fileFlex := []linebot.FlexComponent{
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Margin: linebot.FlexComponentMarginTypeXxl,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  "Total Files",
					Size:  linebot.FlexTextSizeTypeSm,
					Color: "#555555",
					Flex:  linebot.IntPtr(0),
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					// Insert totalFiles
					Text:  strconv.Itoa(totalFiles),
					Size:  linebot.FlexTextSizeTypeSm,
					Color: "#111111",
					Align: linebot.FlexComponentAlignTypeEnd,
				},
			},
		},
	}
	// Generate files
	var files []linebot.FlexComponent
	for _, name := range fileM {
		files = append(files,
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  name,
				Size:  linebot.FlexTextSizeTypeSm,
				Color: "#555555",
			},
		)
	}
	fileFlex = append(fileFlex, files...)
	// Add Separator
	fileFlex = append(fileFlex,
		&linebot.SeparatorComponent{
			Margin: linebot.FlexComponentMarginTypeXxl,
		},
	)
	return fileFlex
}

func buttonFlexComponents(currentFolderID string) []linebot.FlexComponent {
	var buttonFlex []linebot.FlexComponent
	buttonFlex = append(buttonFlex,
		// Button
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Margin: linebot.FlexComponentMarginTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.PostbackAction{
						Label:       "設為上傳資料夾",
						Data:        "action=setFolder&folderID=" + currentFolderID,
						DisplayText: "設為上傳資料夾",
					},
					Style:      linebot.FlexButtonStyleTypePrimary,
					AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
				},
			},
		})
	return buttonFlex

}
