package main

import (
	"github.com/goodcake/facebook-messenger"
	"github.com/joho/godotenv"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	msng := &messenger.Messenger{
		AccessToken:     os.Getenv("AccessToken"),
		VerifyToken:     os.Getenv("VerifyToken"),
		PageID:          os.Getenv("PageID"),
		MessageReceived: messageReceived, // your function for handling received messages, defined below
	}
	// in init or afterwards, you can also specify events when receiving postbacks and message delivery reports from Facebook
	// if you don't want to manage this events, just comment/don't use this receivers
	msng.PostbackReceived = postbackReceived // comment/delete if not used
	// msng.DeliveryReceived = deliveryReceived // comment/delete if not used

	// set URL for your webhook and directly use msng as http Handler
	http.Handle("/webhook", msng)
	appengine.Main()
	//http.ListenAndServe(":8080", nil)
}

// messageReceived is called when you receive message on you webhook i.e. when someone sends message to your chat bot
// params: messenger that received the message, then the user id that sent us message and message data itself
func messageReceived(msng *messenger.Messenger, userID int64, m messenger.FacebookMessage) {

	//log.Println(userID)
	// message received, now lets check what user has sent to us
	switch m.Text {
	case "hello", "hi":
		// someone sent hello or hi, reply with simple text message
		//msng.SendTextMessage(userID, "Hello there")
		qr := msng.NewQuickReplyMessage(userID, " Please Choose:")
		qr.AddNewQuickReply(messenger.TextQuickReply, "OK", "OK", "")
		qr.AddNewQuickReply(messenger.TextQuickReply, "NO", "NO", "")
		qr.AddNewQuickReply(messenger.TextQuickReply, "100", "100", "")
		msng.SendMessage(qr)

	case "send me website":
		// now lets send him some structured message with image and link
		gm := msng.NewGenericMessage(userID)
		gm.AddNewElement("Title", "Subtitle", "http://mysite.com", "http://mysite.com/some-photo.jpeg", nil)

		// GenericMessage can contain up to 10 elements, they are represented as cards and can be scoreled horicontally in messenger
		// So lets add one more element, this time with buttons
		btn1 := msng.NewWebURLButton("Contact US", "http://mysite.com/contact")
		btn2 := msng.NewPostbackButton("Ok", "THIS_DATA_YOU_WILL_RECEIVE_AS_POSTBACK_WHEN_USER_CLICK_THE_BUTTON")
		gm.AddNewElement("Site title", "Subtitle", "http://mysite.com", "http://mysite.com/some-photo.jpeg", []messenger.Button{btn1, btn2})

		// ok, message is ready, lets send
		msng.SendMessage(gm)

	default:
		// upthere we haven't check for errors and responses for cleaner example code
		// but keep in mind that SendMessage returns FacebookResponse struct and error
		// errors are received from Facebook if sometnihg went wrong with message sending
		resp, err := msng.SendTextMessage(userID, m.Text) // echo, send back to user the same text he sent to us
		if err != nil {
			log.Println(err)
			return // if there is an error, resp is empty struct, useless
		}
		log.Println("Message ID", resp.MessageID, "sent to user", resp.RecipientID)
		// store resp.MessageID if you want to track delivery reports that will be sent later from Facebook
	}
}

// postbackReceived is called when you reiceive postback event from Facebook server
func postbackReceived(msng *messenger.Messenger, userID int64, p messenger.FacebookPostback) {
	if p.Payload == "THIS_DATA_YOU_WILL_RECEIVE_AS_POSTBACK_WHEN_USER_CLICK_THE_BUTTON" {
		// user just clicked Ok button from previouse example, lets just send him a message
		msng.SendTextMessage(userID, "Ok, I'm always online, chat with me anytime :)")
	}
}

// deliveryReceived is used if you want to track delivery reports for sent messages
func deliveryReceived(msng *messenger.Messenger, userID int64, d messenger.FacebookDelivery) {
	for _, mid := range d.Mids {
		log.Println("Message delivered, msgID:", mid)
	}
}
