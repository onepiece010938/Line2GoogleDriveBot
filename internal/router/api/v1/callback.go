package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
)

type LineHandler struct {
	bot *linebot.Client
}

func Callback(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("In Callback2")
		log.Println(c.Request)
		ctx := c.Request.Context()
		events, err := app.LineBotClient.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {

				c.JSON(http.StatusBadRequest, nil)
			} else {
				c.JSON(http.StatusInternalServerError, nil)
			}
			return
		}
		lineHandler := LineHandler{bot: app.LineBotClient}
		for _, event := range events {
			log.Printf("Got event %v", event)
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if err := lineHandler.handleText(ctx, app, message, event.ReplyToken, event.Source); err != nil {
						log.Print(err)
					}

				case *linebot.FileMessage:
					if err := lineHandler.handleFile(message, event.ReplyToken); err != nil {
						log.Print(err)
					}
				default:
					log.Printf("Unknown message: %v", message)
				}
			case linebot.EventTypeFollow:
				if err := lineHandler.replyText(event.ReplyToken, "Got followed event"); err != nil {
					log.Print(err)
				}
			case linebot.EventTypeUnfollow:
				log.Printf("Unfollowed this bot: %v", event)
			case linebot.EventTypeJoin:
				if err := lineHandler.replyText(event.ReplyToken, "Joined "+string(event.Source.Type)); err != nil {
					log.Print(err)
				}
			case linebot.EventTypeLeave:
				log.Printf("Left: %v", event)
			case linebot.EventTypePostback:
				data := event.Postback.Data
				if data == "DATE" || data == "TIME" || data == "DATETIME" {
					data += fmt.Sprintf("(%v)", *event.Postback.Params)
				}
				if err := lineHandler.replyText(event.ReplyToken, "Got postback: "+data); err != nil {
					log.Print(err)
				}
			case linebot.EventTypeBeacon:
				if err := lineHandler.replyText(event.ReplyToken, "Got beacon: "+event.Beacon.Hwid); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown event: %v", event)
			}
		}
		//events, err := myBot.ParseRequest(c.Request)
		// ctx := c.Request.Context()
		// err := AnalyzeService.CreateAnalyze(ctx, analyze.CreateAnalyzeParm{})
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// c.JSON(http.StatusOK, "")

	}

}

func (l *LineHandler) handleText(ctx context.Context, app *app.Application, message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "sample":
		test, err := app.AnalyzeService.AnalyzeTest(ctx)
		if err != nil {
			return l.replyText(replyToken, "ERROR")
		}
		return l.replyText(replyToken, test)
	case "profile":
		if source.UserID != "" {
			profile, err := l.bot.GetProfile(source.UserID).Do()
			if err != nil {
				return l.replyText(replyToken, err.Error())
			}
			if _, err := l.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("Display name: "+profile.DisplayName),
				linebot.NewTextMessage("Status message: "+profile.StatusMessage),
			).Do(); err != nil {
				return err
			}
		} else {
			return l.replyText(replyToken, "Bot can't use profile API without user ID")
		}

	case "confirm":
		template := linebot.NewConfirmTemplate(
			"Do it?",
			linebot.NewMessageAction("Yes", "Yes!"),
			linebot.NewMessageAction("No", "No!"),
		)
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewTemplateMessage("Confirm alt text", template),
		).Do(); err != nil {
			return err
		}

	case "datetime":
		template := linebot.NewButtonsTemplate(
			"", "", "Select date / time !",
			linebot.NewDatetimePickerAction("date", "DATE", "date", "", "", ""),
			linebot.NewDatetimePickerAction("time", "TIME", "time", "", "", ""),
			linebot.NewDatetimePickerAction("datetime", "DATETIME", "datetime", "", "", ""),
		)
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewTemplateMessage("Datetime pickers alt text", template),
		).Do(); err != nil {
			return err
		}
	case "flex":
		// {
		//   "type": "bubble",
		//   "body": {
		//     "type": "box",
		//     "layout": "horizontal",
		//     "contents": [
		//       {
		//         "type": "text",
		//         "text": "Hello,"
		//       },
		//       {
		//         "type": "text",
		//         "text": "World!"
		//       }
		//     ]
		//   }
		// }
		contents := &linebot.BubbleContainer{
			Type: linebot.FlexContainerTypeBubble,
			Body: &linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeHorizontal,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type: linebot.FlexComponentTypeText,
						Text: "Hello,",
					},
					&linebot.TextComponent{
						Type: linebot.FlexComponentTypeText,
						Text: "World!",
					},
				},
			},
		}
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("Flex message alt text", contents),
		).Do(); err != nil {
			return err
		}
	case "flex carousel":
		// {
		//   "type": "carousel",
		//   "contents": [
		//     {
		//       "type": "bubble",
		//       "body": {
		//         "type": "box",
		//         "layout": "vertical",
		//         "contents": [
		//           {
		//             "type": "text",
		//             "text": "First bubble"
		//           }
		//         ]
		//       }
		//     },
		//     {
		//       "type": "bubble",
		//       "body": {
		//         "type": "box",
		//         "layout": "vertical",
		//         "contents": [
		//           {
		//             "type": "text",
		//             "text": "Second bubble"
		//           }
		//         ]
		//       }
		//     }
		//   ]
		// }
		contents := &linebot.CarouselContainer{
			Type: linebot.FlexContainerTypeCarousel,
			Contents: []*linebot.BubbleContainer{
				{
					Type: linebot.FlexContainerTypeBubble,
					Body: &linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.TextComponent{
								Type: linebot.FlexComponentTypeText,
								Text: "First bubble",
							},
						},
					},
				},
				{
					Type: linebot.FlexContainerTypeBubble,
					Body: &linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.TextComponent{
								Type: linebot.FlexComponentTypeText,
								Text: "Second bubble",
							},
						},
					},
				},
			},
		}
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("Flex message alt text", contents),
		).Do(); err != nil {
			return err
		}
	case "flex json":
		jsonString := `{
  "type": "bubble",
  "hero": {
    "type": "image",
    "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_1_cafe.png",
    "size": "full",
    "aspectRatio": "20:13",
    "aspectMode": "cover",
    "action": {
      "type": "uri",
      "uri": "http://linecorp.com/"
    }
  },
  "body": {
    "type": "box",
    "layout": "vertical",
    "contents": [
      {
        "type": "text",
        "text": "Brown Cafe",
        "weight": "bold",
        "size": "xl"
      },
      {
        "type": "box",
        "layout": "baseline",
        "margin": "md",
        "contents": [
          {
            "type": "icon",
            "size": "sm",
            "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
          },
          {
            "type": "icon",
            "size": "sm",
            "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
          },
          {
            "type": "icon",
            "size": "sm",
            "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
          },
          {
            "type": "icon",
            "size": "sm",
            "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
          },
          {
            "type": "icon",
            "size": "sm",
            "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png"
          },
          {
            "type": "text",
            "text": "4.0",
            "size": "sm",
            "color": "#999999",
            "margin": "md",
            "flex": 0
          }
        ]
      },
      {
        "type": "box",
        "layout": "vertical",
        "margin": "lg",
        "spacing": "sm",
        "contents": [
          {
            "type": "box",
            "layout": "baseline",
            "spacing": "sm",
            "contents": [
              {
                "type": "text",
                "text": "Place",
                "color": "#aaaaaa",
                "size": "sm",
                "flex": 1
              },
              {
                "type": "text",
                "text": "Miraina Tower, 4-1-6 Shinjuku, Tokyo",
                "wrap": true,
                "color": "#666666",
                "size": "sm",
                "flex": 5
              }
            ]
          },
          {
            "type": "box",
            "layout": "baseline",
            "spacing": "sm",
            "contents": [
              {
                "type": "text",
                "text": "Time",
                "color": "#aaaaaa",
                "size": "sm",
                "flex": 1
              },
              {
                "type": "text",
                "text": "10:00 - 23:00",
                "wrap": true,
                "color": "#666666",
                "size": "sm",
                "flex": 5
              }
            ]
          }
        ]
      }
    ]
  },
  "footer": {
    "type": "box",
    "layout": "vertical",
    "spacing": "sm",
    "contents": [
      {
        "type": "button",
        "style": "link",
        "height": "sm",
        "action": {
          "type": "uri",
          "label": "CALL",
          "uri": "https://linecorp.com"
        }
      },
      {
        "type": "button",
        "style": "link",
        "height": "sm",
        "action": {
          "type": "uri",
          "label": "WEBSITE",
          "uri": "https://linecorp.com",
          "altUri": {
            "desktop": "https://line.me/ja/download"
          }
        }
      },
      {
        "type": "spacer",
        "size": "sm"
      }
    ],
    "flex": 0
  }
}`
		contents, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonString))
		if err != nil {
			return err
		}
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("Flex message alt text", contents),
		).Do(); err != nil {
			return err
		}
	case "imagemap":
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewImagemapMessage(
				"appBaseURL"+"/static/rich",
				"Imagemap alt text",
				linebot.ImagemapBaseSize{Width: 1040, Height: 1040},
				linebot.NewURIImagemapAction("LINE Store Manga", "https://store.line.me/family/manga/en", linebot.ImagemapArea{X: 0, Y: 0, Width: 520, Height: 520}),
				linebot.NewURIImagemapAction("LINE Store Music", "https://store.line.me/family/music/en", linebot.ImagemapArea{X: 520, Y: 0, Width: 520, Height: 520}),
				linebot.NewURIImagemapAction("LINE Store Play", "https://store.line.me/family/play/en", linebot.ImagemapArea{X: 0, Y: 520, Width: 520, Height: 520}),
				linebot.NewMessageImagemapAction("URANAI!", "URANAI!", linebot.ImagemapArea{X: 520, Y: 520, Width: 520, Height: 520}),
			),
		).Do(); err != nil {
			return err
		}

	default:
		log.Printf("Echo message to %s: %s", replyToken, message.Text)
		if _, err := l.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(message.Text),
		).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (l *LineHandler) handleFile(message *linebot.FileMessage, replyToken string) error {
	return l.replyText(replyToken, fmt.Sprintf("File `%s` (%d bytes) received.", message.FileName, message.FileSize))
}

func (l *LineHandler) replyText(replyToken, text string) error {
	if _, err := l.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

/*
func (l *LineHandler) handleHeavyContent(messageID string, callback func(*os.File) error) error {
	content, err := l.bot.GetMessageContent(messageID).Do()
	if err != nil {
		return err
	}
	defer content.Content.Close()
	log.Printf("Got file: %s", content.ContentType)
	originalContent, err := saveContent(content.Content)
	if err != nil {
		return err
	}
	return callback(originalContent)
}

func saveContent(content io.ReadCloser) (*os.File, error) {
	file, err := ioutil.TempFile("downloadDir", "")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return nil, err
	}
	log.Printf("Saved %s", file.Name())
	return file, nil
}
*/
