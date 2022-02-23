package adapters

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/adlio/trello"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type TrelloAdapter struct {
	KEY      string
	TOKEN    string
	BOARD_ID string
	LIST_ID  string
	LABELS   []string
}

func initTrelloStruct() TrelloAdapter {
	return TrelloAdapter{
		KEY:      viper.GetString("adapters.trello.key"),
		TOKEN:    viper.GetString("adapters.trello.token"),
		BOARD_ID: viper.GetString("adapters.trello.board_id"),
		LIST_ID:  viper.GetString("adapters.trello.list_id"),
		LABELS:   viper.GetStringSlice("adapters.trello.labels"),
	}
}

func (t *TrelloAdapter) Preflight() error {
	pterm.Info.Println("Preflight check for Trello Adapter")

	conf := initTrelloStruct()

	if conf.KEY == "" || conf.TOKEN == "" || conf.BOARD_ID == "" || conf.LIST_ID == "" {
		return errors.New("trello configuration is not set properly")
	}

	api := trello.NewClient(conf.KEY, conf.TOKEN)

	// Validate that user can request a specific board
	_, err := api.GetBoard(conf.BOARD_ID, trello.Defaults())
	if err != nil {
		return err
	}

	return nil
}

func (t *TrelloAdapter) Execute(imagePath string) error {
	conf := initTrelloStruct()

	/*
		Steps:
		1. Create a card without the labels attached
		2. Check if labels exist
			a. If they exist get all board labels
			b. Try to match text labels with label IDs
			c. If no match, ignore this label
			d. Make request to add each label matched
		3. Attach images to the card
	*/

	api := trello.NewClient(conf.KEY, conf.TOKEN)

	currentTime := time.Now().Format(time.RFC1123)

	// Create card
	card := trello.Card{
		Name:   "Cyclops Testing 🌁",
		Desc:   fmt.Sprintf("Date: %s", currentTime),
		Pos:    1,
		IDList: conf.LIST_ID,
		// For some reason, labels cannot be added directly to the card
		// and only if the card is created.
		// Issue: https://github.com/adlio/trello/issues/45
		// Labels: -----,
	}

	list, err := api.GetList(conf.LIST_ID, trello.Defaults())

	if err != nil {
		return err
	}

	err = list.AddCard(&card)

	if len(conf.LABELS) > 0 {
		// Check and find labels
		board, err := api.GetBoard(conf.BOARD_ID, trello.Defaults())
		if err != nil {
			return nil
		}
		labels, err := board.GetLabels()

		// Match Label string value to Label ID based on Trello's response
		labelIds := make(map[string]string)
		for _, userLabel := range conf.LABELS {
			found := false
			for _, label := range labels {
				if label.Name == userLabel {
					labelIds[userLabel] = label.ID
					found = true
					break
				}
			}
			if !found {
				pterm.Warning.Printfln("Label '%v' not found on Trello board", userLabel)
			}
		}

		if err != nil {
			return err
		}

		for labelName, labelID := range labelIds {
			err = card.AddIDLabel(labelID)
			if err != nil {
				pterm.Warning.Printfln("Failed to add label '%v' to card", labelName)
			}
		}

		if err != nil {
			return err
		}
	}

	// Add images as attachments
	filesSortedCleaned, err := SortedImageFileNames(imagePath)

	if err != nil {
		return err
	}

	for _, file := range filesSortedCleaned {
		data, err := os.Open(fmt.Sprintf("%s/%s", imagePath, file))
		if err != nil {
			return err
		}

		attachment := trello.Attachment{
			Name: file,
		}

		card.AddFileAttachment(&attachment, file, data)
	}

	pterm.Success.Sprintfln("Trello card updated. Visit here for info: %s\n",
		"https://trello.com/c/https://trello.com/c/KAPUAWc3/419-loceye-invoice-info")

	return nil
}
