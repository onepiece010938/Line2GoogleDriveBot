package v1

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	domainDrive "github.com/onepiece010938/Line2GoogleDriveBot/internal/domain/drive"
)

func Callback(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		events, err := app.LineBotClient.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err)
			} else {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err)
			}
			return
		}

		for _, event := range events {
			// Get user Line ID
			lineID := event.Source.UserID
			// Handle Button Postback
			if event.Type == linebot.EventTypePostback {
				// 如果是 Postback 事件，取得 postback 資料
				postbackData := event.Postback.Data
				log.Printf("Postback data: %s", postbackData)
				// 解析 postback 資料
				values, err := url.ParseQuery(postbackData)
				if err != nil {
					log.Printf("Error parsing postback data: %v", err)
					return
				}
				// 取得特定參數的值，setFolder || openFolder
				action := values.Get("action")
				folderID := values.Get("folderID")

				// 在這裡可以根據 action 和 FolderID 做相應的處理
				log.Printf("Action: %s, FolderID: %s", action, folderID)
				if action == "openFolder" {
					res, err := app.DriveService.ListSelectedFolderCarousel(ctx, lineID, folderID)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err := app.LineBotClient.ReplyMessage(
						event.ReplyToken,
						linebot.NewFlexMessage("打開資料夾", res.CarouselContainer),
					).Do(); err != nil {
						log.Println(err)
						return
					}
				}
				if action == "setFolder" {
					err := app.DriveService.SetUploadPath(ctx, lineID, folderID)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功設定上傳路徑")).Do(); err != nil {
						log.Println(err)
					}
					return

				}
			}
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "[登入]" {
						// fmt.Println("Linebot GET", lineID)
						// profile, _ := app.LineBotClient.GetProfile(lineID).Do()
						// fmt.Println("Hi~ ", profile.DisplayName)

						authURL := app.DriveService.LoginURL(ctx, lineID)
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(authURL)).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list" {
						res, err := app.DriveService.ListFiles(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list folder" {
						res, err := app.DriveService.ListMyDriveFolders(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list shared" {
						res, err := app.DriveService.ListSharedFolders(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "flex carousel" {

						contents := &linebot.CarouselContainer{
							Type: linebot.FlexContainerTypeCarousel,
							Contents: []*linebot.BubbleContainer{
								{
									Type: linebot.FlexContainerTypeBubble,
									Body: &linebot.BoxComponent{
										Type:   linebot.FlexComponentTypeBox,
										Layout: linebot.FlexBoxLayoutTypeVertical,
										Contents:
										// append(title, otherComponents...),
										[]linebot.FlexComponent{

											&linebot.TextComponent{
												Type:   linebot.FlexComponentTypeText,
												Text:   "FOLDER",
												Weight: linebot.FlexTextWeightTypeBold,
												Color:  "#1DB446",
												Size:   linebot.FlexTextSizeTypeSm,
											},
											&linebot.TextComponent{
												Type:   linebot.FlexComponentTypeText,
												Text:   "Folder Name",
												Weight: linebot.FlexTextWeightTypeBold,
												Size:   linebot.FlexTextSizeTypeXxl,
												Margin: linebot.FlexComponentMarginTypeMd,
											},
											&linebot.TextComponent{
												Type:  linebot.FlexComponentTypeText,
												Text:  "/path/to/floder",
												Size:  linebot.FlexTextSizeTypeXs,
												Color: "#aaaaaa",
												Wrap:  true,
											},
											&linebot.SeparatorComponent{
												Type:   linebot.FlexComponentTypeSeparator,
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
											&linebot.BoxComponent{
												Type:    linebot.FlexComponentTypeBox,
												Layout:  linebot.FlexBoxLayoutTypeVertical,
												Margin:  linebot.FlexComponentMarginTypeXxl,
												Spacing: linebot.FlexComponentSpacingTypeSm,
												Contents: []linebot.FlexComponent{
													&linebot.BoxComponent{
														Type:   linebot.FlexComponentTypeBox,
														Layout: linebot.FlexBoxLayoutTypeHorizontal,
														Contents: []linebot.FlexComponent{
															&linebot.TextComponent{
																Type:       linebot.FlexComponentTypeText,
																Text:       "Folder1",
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
																	Data:        "folderid1",
																	DisplayText: "進入Folder1",
																},
																Style:      linebot.FlexButtonStyleTypeLink,
																Height:     linebot.FlexButtonHeightTypeSm,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
																AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
															},
														},
													},
													&linebot.BoxComponent{
														Type:   linebot.FlexComponentTypeBox,
														Layout: linebot.FlexBoxLayoutTypeHorizontal,
														Contents: []linebot.FlexComponent{
															&linebot.TextComponent{
																Type:       linebot.FlexComponentTypeText,
																Text:       "Folder2",
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
																	Data:        "folderid2",
																	DisplayText: "進入Folder2",
																},
																Style:      linebot.FlexButtonStyleTypeLink,
																Height:     linebot.FlexButtonHeightTypeSm,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
																AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
															},
														},
													},
												},
											},
											// Separator
											&linebot.SeparatorComponent{
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
											// Files
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
														Type:  linebot.FlexComponentTypeText,
														Text:  "3",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#111111",
														Align: linebot.FlexComponentAlignTypeEnd,
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE1",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE2",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE3",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											// Separator
											&linebot.SeparatorComponent{
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
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
															Data:        "folderid",
															DisplayText: "設為上傳資料夾",
														},
														Style:      linebot.FlexButtonStyleTypePrimary,
														AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
													},
												},
											},
										},
									},
									Styles: &linebot.BubbleStyle{
										Footer: &linebot.BlockStyle{
											Separator: true,
										},
									},
								},
							},
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("Flex message alt text", contents),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "[我の雲端硬碟]" {
						res, err := app.DriveService.ListFolderCarousel(ctx, lineID, domainDrive.PersonalFolder)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("測試Flex Carousel", res.CarouselContainer),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "[共用硬碟]" {
						res, err := app.DriveService.ListFolderCarousel(ctx, lineID, domainDrive.SharedFolder)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("測試Flex Carousel", res.CarouselContainer),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "[上傳路徑]" {
						path, err := app.DriveService.GetUploadPath(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(path)).Do(); err != nil {
							log.Println(err)
						}
						return
					}

					// samplePK, err := app.SampleService.Sample(ctx, message.Text)
					// if err != nil {
					// 	log.Println(err)
					// 	return
					// }
					// if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(samplePK)).Do(); err != nil {
					// 	log.Println(err)
					// }

				case *linebot.FileMessage:
					content, err := app.LineBotClient.GetMessageContent(message.ID).Do()
					if err != nil {
						log.Println(err)
						return
					}

					log.Printf("Got file: %s", content.ContentType)
					fmt.Printf("File `%s` (%d bytes) received.", message.FileName, message.FileSize)

					err = app.DriveService.UploadFile(ctx, lineID, message.FileName, content.Content)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功上傳: "+message.FileName)).Do(); err != nil {
						log.Println(err)
					}

				}

			}

		}

	}

}
